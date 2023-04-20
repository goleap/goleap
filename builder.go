package dbkit

import (
	"context"
	"fmt"
	"github.com/kitstack/dbkit/connector/drivers"
	"github.com/kitstack/dbkit/connector/drivers/operators"
	"github.com/kitstack/dbkit/definitions"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/depkit"
	"github.com/kitstack/structkit"
	structKitSpecs "github.com/kitstack/structkit/specs"
	"reflect"
	"strings"
	"sync"
)

func init() {
	depkit.Register[structKitSpecs.Get](structkit.Get)
	depkit.Register[structKitSpecs.Set](structkit.Set)
}

var (
	QueryTypeGet     = "Get"
	QueryTypeFindAll = "FindAll"
	QueryTypeFind    = "Find"
)

type builder[T specs.Model] struct {
	context.Context
	specs.Connector
	sync.Mutex

	queryType string

	model           T
	modelDefinition specs.ModelDefinition

	fields []string
	wheres []specs.Condition

	selectedFieldsDefinition []specs.FieldDefinition

	driverFields []specs.DriverField
	driverJoins  []specs.DriverJoin
	driverWheres []specs.DriverWhere

	payload specs.PayloadAugmented[T]

	dependencies map[string]specs.ModelDefinition
}

func (o *builder[T]) setQueryType(queryType string) specs.Builder[T] {
	if o.queryType != "" {
		return o
	}
	o.queryType = queryType
	return o
}

func (o *builder[T]) QueryType() string {
	return o.queryType
}

func (o *builder[T]) execute(flow ...func() error) (err error) {
	for _, f := range flow {
		if err := f(); err != nil {
			return err
		}
	}

	return
}

func (o *builder[T]) valideRequiredField() error {
	if len(o.selectedFieldsDefinition) > 0 {
		return nil
	}

	return NewFieldRequiredError(o.QueryType())
}

func (o *builder[T]) buildFields() (err error) {
	for _, fieldName := range o.fields {
		field, err := o.modelDefinition.GetFieldByName(fieldName)

		if err != nil {
			return err
		}

		if field.FromSlice() {
			o.addDependency(field)
			continue
		}

		o.selectedFieldsDefinition = append(o.selectedFieldsDefinition, field)
	}

	return
}

func (o *builder[T]) buildWheres() (err error) {
	for _, where := range o.wheres {

		fieldDefinition, err := o.modelDefinition.GetFieldByName(where.From())
		if err != nil {
			return err
		}

		o.driverWheres = append(o.driverWheres, drivers.NewWhere().SetFrom(fieldDefinition.Field()).SetOperator(where.Operator()).SetTo(where.To()))
	}

	return
}

func (o *builder[T]) addDependency(field specs.FieldDefinition) {
	o.dependencies[field.Model().FromField().RecursiveFullName()] = field.Model()
}

func (o *builder[T]) getDriverFields() []specs.DriverField {

	for _, field := range o.selectedFieldsDefinition {
		o.driverFields = append(o.driverFields, field.Field())
	}

	return o.driverFields
}

func (o *builder[T]) getDriverJoins() ([]specs.DriverJoin, error) {

	uniqueJoins := map[string]specs.DriverJoin{}
	for _, field := range o.selectedFieldsDefinition {

		for _, join := range field.Join() {
			formatted, err := join.Formatted()
			if err != nil {
				return nil, err
			}
			uniqueJoins[formatted] = join
		}

	}

	for _, join := range uniqueJoins {
		o.driverJoins = append(o.driverJoins, join)
	}

	return o.driverJoins, nil
}

func (o *builder[T]) getDriverWheres() []specs.DriverWhere {
	return o.driverWheres
}

func (o *builder[T]) buildPayload() error {
	o.payload = NewPayload[T](o.model)
	o.payload.SetFields(o.getDriverFields())
	o.payload.SetWheres(o.getDriverWheres())

	joins, err := o.getDriverJoins()
	if err != nil {
		return err
	}

	o.payload.SetJoins(joins)

	return nil
}

func (o *builder[T]) Payload() specs.PayloadAugmented[T] {
	return o.payload
}

func (o *builder[T]) Get(primaryKeyValue any) (result T, err error) {
	o.setQueryType(QueryTypeGet)

	primaryField, err := o.modelDefinition.GetPrimaryField()
	if err != nil {
		return
	}

	o.Where(NewCondition().SetFrom(primaryField.RecursiveFullName()).
		SetOperator(operators.Equal).
		SetTo(primaryKeyValue))

	return o.Find()
}

func (o *builder[T]) Find() (result T, err error) {
	data, err := o.setQueryType(QueryTypeFind).FindAll()
	if err != nil {
		return
	}

	if len(data) == 0 {
		err = NewNotFoundError(o.modelDefinition.ModelValue().Type().Name())
		return
	}

	return data[0], nil
}

func (o *builder[T]) FindAll() ([]T, error) {
	o.setQueryType(QueryTypeFindAll)

	err := o.execute(
		o.buildFields,
		o.valideRequiredField,
		o.buildWheres,
		o.buildPayload,
	)

	if err != nil {
		return nil, err
	}

	err = o.Connector.Select(o.Context, o.payload)
	if err != nil {
		return nil, err
	}

	for fundamentalName, m := range o.dependencies {
		from, err := m.FromField().GetByColumn()
		if err != nil {
			return []T{}, err
		}

		to, err := m.FromField().GetToColumn()
		if err != nil {
			return []T{}, err
		}

		var in []any
		var mapping = map[any][]int{}
		for index, result := range o.Payload().Result() {
			v := depkit.Get[structKitSpecs.Get]()(result, from.RecursiveFullName())
			if v == nil {
				continue
			}

			in = append(in, v)
			mapping[v] = append(mapping[v], index)
		}

		toFieldName := strings.Replace(to.RecursiveFullName(), fmt.Sprintf("%s.", fundamentalName), "", 1)

		sub := Use[specs.Model](o.Context, o.Connector).
			SetModel(to.Model().Copy()).
			Fields(append(o.extractFieldsByFundamentalName(fundamentalName), toFieldName)...)

		sub.Where(NewCondition().SetFrom(toFieldName).SetOperator(operators.In).SetTo(in))
		for _, where := range o.extractWheresByFundamentalName(fundamentalName) {
			sub.Where(where)
		}

		data, err := sub.FindAll()

		if err != nil {
			return []T{}, err
		}

		for _, result := range data {
			index := depkit.Get[structKitSpecs.Get]()(result, toFieldName)
			for _, index := range mapping[index] {
				err := depkit.Get[structKitSpecs.Set]()(o.Payload().Result()[index], fmt.Sprintf("%s.%s", fundamentalName, "[*]"), result)
				if err != nil {
					return []T{}, err
				}
			}
		}
	}

	return o.Payload().Result(), nil
}

func (o *builder[T]) extractFieldsByFundamentalName(fundamentalName string) (fields []string) {
	for _, field := range o.fields {
		if strings.HasPrefix(field, fmt.Sprintf("%s.", fundamentalName)) {
			fields = append(fields, strings.Replace(field, fmt.Sprintf("%s.", fundamentalName), "", 1))
		}
	}
	return
}

func (o *builder[T]) extractWheresByFundamentalName(fundamentalName string) (fields []specs.Condition) {
	for _, field := range o.wheres {
		if strings.HasPrefix(field.From(), fmt.Sprintf("%s.", fundamentalName)) {
			field.SetFrom(strings.Replace(field.From(), fmt.Sprintf("%s.", fundamentalName), "", 1))
			fields = append(fields, field)
		}
	}
	return
}

func (o *builder[T]) Delete(_ any) error {
	//TODO implement me
	panic("implement me")
}

func (o *builder[T]) Create() (err error) {
	//TODO implement me
	panic("implement me")
}

func (o *builder[T]) Update() error {
	//TODO implement me
	panic("implement me")
}

func (o *builder[T]) Fields(field ...string) specs.Builder[T] {
	o.fields = field
	return o
}

func (o *builder[T]) Where(where specs.Condition) specs.Builder[T] {
	o.wheres = append(o.wheres, where)
	return o
}

func (o *builder[T]) Limit(_ int) specs.Builder[T] {
	//TODO implement me
	panic("implement me")
}

func (o *builder[T]) Offset(_ int) specs.Builder[T] {
	//TODO implement me
	panic("implement me")
}

func (o *builder[T]) OrderBy(_ ...string) specs.Builder[T] {
	//TODO implement me
	panic("implement me")
}

func (o *builder[T]) Count() (total int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (o *builder[T]) SetModel(model T) specs.Builder[T] {
	o.model = model
	o.modelDefinition = definitions.Use(model).Parse()

	return o
}

func Use[T specs.Model](ctx context.Context, connector specs.Connector) specs.Builder[T] {
	var model T

	var builder specs.Builder[T] = &builder[T]{
		Context:      ctx,
		Connector:    connector,
		dependencies: make(map[string]specs.ModelDefinition),
	}

	t := reflect.ValueOf(model)
	if t.IsValid() {
		builder.SetModel(model)
	}

	return builder
}
