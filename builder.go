package dbkit

import (
	"context"
	"github.com/kitstack/dbkit/connector/drivers"
	"github.com/kitstack/dbkit/connector/drivers/operators"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/depkit"
	"reflect"
	"sort"
	"sync"
)

var (
	QueryTypeGet     = "Get"
	QueryTypeFindAll = "FindAll"
	QueryTypeFind    = "Find"
)

type builder[T specs.Model] struct {
	sync.Mutex

	context   context.Context
	connector specs.Connector

	subBuilder specs.SubBuilder[T]

	queryType string

	model           T
	modelDefinition specs.ModelDefinition

	fields []string
	wheres []specs.Condition

	selectedFieldsDefinition []specs.FieldDefinition
	filterFieldsDefinition   []specs.FieldDefinition

	driverFields []specs.DriverField
	driverJoins  []specs.DriverJoin
	driverWheres []specs.DriverWhere

	payload specs.PayloadAugmented[T]
}

func (o *builder[T]) Context() context.Context {
	return o.context
}

func (o *builder[T]) SubBuilder() specs.SubBuilder[T] {
	return o.subBuilder
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

	return NewErrFieldRequired(o.QueryType())
}

func (o *builder[T]) buildFields() (err error) {
	for _, fieldName := range o.fields {
		field, err := o.modelDefinition.GetFieldByName(fieldName)

		if err != nil {
			return err
		}

		if field.FromSlice() {
			o.subBuilder.AddJob(o, field.Model().FromField().FundamentalName(), field.Model())
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

		o.filterFieldsDefinition = append(o.filterFieldsDefinition, fieldDefinition)
		o.driverWheres = append(o.driverWheres, drivers.NewWhere().SetFrom(fieldDefinition.Field()).SetOperator(where.Operator()).SetTo(where.To()))
	}

	return
}

func (o *builder[T]) getDriverFields() []specs.DriverField {

	for _, field := range o.selectedFieldsDefinition {
		o.driverFields = append(o.driverFields, field.Field())
	}

	return o.driverFields
}

func (o *builder[T]) getDriverJoins() ([]specs.DriverJoin, error) {
	uniqueJoins := map[string]specs.DriverJoin{}
	uniqueFields := map[string]bool{}
	mergedField := append(o.selectedFieldsDefinition, o.filterFieldsDefinition...)

	for _, field := range mergedField {
		recursiveFullName := field.RecursiveFullName()
		if _, ok := uniqueFields[recursiveFullName]; ok {
			continue
		}

		for _, join := range field.Join() {
			formatted, err := join.Formatted()
			if err != nil {
				return nil, err
			}
			uniqueJoins[formatted] = join
			uniqueFields[recursiveFullName] = true
		}
	}

	for _, join := range uniqueJoins {
		o.driverJoins = append(o.driverJoins, join)
	}

	sort.SliceStable(o.driverJoins, func(i, j int) bool {
		return o.driverJoins[i].To().Index() < o.driverJoins[j].To().Index()
	})

	return o.driverJoins, nil
}

func (o *builder[T]) getDriverWheres() []specs.DriverWhere {
	return o.driverWheres
}

func (o *builder[T]) Connector() (specs.Connector, error) {
	if o.connector == nil {
		connector, err := depkit.Get[specs.ConnectorsInstance]()().Get(o.model.ConnectorName())
		if err != nil {
			return nil, err
		}
		o.connector = connector
	}

	return o.connector, nil
}

func (o *builder[T]) buildPayload() error {
	o.payload = depkit.Get[specs.NewPayload[T]]()(o.model)
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

	o.wheres = append([]specs.Condition{NewCondition().SetFrom(primaryField.RecursiveFullName()).SetOperator(operators.Equal).SetTo(primaryKeyValue)}, o.wheres...)

	return o.Find()
}

func (o *builder[T]) Find() (result T, err error) {
	data, err := o.setQueryType(QueryTypeFind).FindAll()
	if err != nil {
		return
	}

	if len(data) == 0 {
		err = NewErrNotFound(o.modelDefinition.TypeName())
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

	connector, err := o.Connector()
	if err != nil {
		return nil, err
	}

	err = connector.Select(o.Context(), o.Payload())
	if err != nil {
		return nil, err
	}

	err = o.SubBuilder().Execute()
	if err != nil {
		return nil, err
	}

	return o.Payload().Result(), nil
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

func (o *builder[T]) SetFields(field ...string) specs.Builder[T] {
	o.fields = field
	return o
}

func (o *builder[T]) SetWhere(where specs.Condition) specs.Builder[T] {
	o.wheres = append(o.wheres, where)
	return o
}

func (o *builder[T]) Wheres() []specs.Condition {
	return o.wheres
}

func (o *builder[T]) Fields() []string {
	return o.fields
}

func (o *builder[T]) SetLimit(_ int) specs.Builder[T] {
	//TODO implement me
	panic("implement me")
}

func (o *builder[T]) SetOffset(_ int) specs.Builder[T] {
	//TODO implement me
	panic("implement me")
}

func (o *builder[T]) SetOrderBy(_ ...string) specs.Builder[T] {
	//TODO implement me
	panic("implement me")
}

func (o *builder[T]) Count() (total int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (o *builder[T]) SetModel(model T) specs.Builder[T] {
	o.model = model
	o.modelDefinition = depkit.Get[specs.UseModelDefinition]()(model).Parse()

	depkit.Register[specs.NewPayload[T]](NewPayload[T])

	return o
}

func Use[T specs.Model](ctx context.Context) specs.Builder[T] {
	var model T

	injectGenericDependencies[T]()

	var builder specs.Builder[T] = &builder[T]{
		context:    ctx,
		subBuilder: depkit.Get[specs.NewSubBuilder[T]]()(),
	}

	t := reflect.ValueOf(model)
	if t.IsValid() {
		builder.SetModel(model)
	}

	return builder
}
