package dbkit

import (
	"context"
	"github.com/lab210-dev/dbkit/connector/drivers"
	"github.com/lab210-dev/dbkit/connector/drivers/operators"
	"github.com/lab210-dev/dbkit/modeldefinition"
	"github.com/lab210-dev/dbkit/specs"
)

type builder[T specs.Model] struct {
	context.Context
	specs.Connector

	queryType string

	model           *T
	modelDefinition specs.ModelDefinition

	fields []string
	wheres []specs.Condition

	selectedFieldsDefinition []specs.FieldDefinition

	driverFields []specs.DriverField
	driverJoins  []specs.DriverJoin
	driverWheres []specs.DriverWhere

	payload specs.PayloadAugmented[T]
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

	return NewFieldRequiredError("")
}

func (o *builder[T]) buildFields() (err error) {
	for _, fieldName := range o.fields {
		field, err := o.modelDefinition.GetFieldByName(fieldName)

		if err != nil {
			return err
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

func (o *builder[T]) getDriverFields() []specs.DriverField {

	for _, field := range o.selectedFieldsDefinition {
		o.driverFields = append(o.driverFields, field.Field())
	}

	return o.driverFields
}

func (o *builder[T]) getDriverJoins() []specs.DriverJoin {

	for _, field := range o.selectedFieldsDefinition {
		o.driverJoins = append(o.driverJoins, field.Join()...)
	}

	return o.driverJoins
}

func (o *builder[T]) getDriverWheres() []specs.DriverWhere {
	return o.driverWheres
}

func (o *builder[T]) buildPayload() error {
	o.payload = NewPayload[T]()
	o.payload.SetFields(o.getDriverFields())
	o.payload.SetJoins(o.getDriverJoins())
	o.payload.SetWheres(o.getDriverWheres())

	return nil
}

func (o *builder[T]) Payload() specs.PayloadAugmented[T] {
	return o.payload
}

func (o *builder[T]) Get(primaryKeyValue any) (result T, err error) {
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

	err = o.execute(
		o.buildFields,
		o.valideRequiredField,
		o.buildWheres,
		o.buildPayload,
	)

	if err != nil {
		return
	}

	err = o.Connector.Select(o.Context, o.payload)
	if err != nil {
		return
	}

	if len(o.Payload().Result()) == 0 {
		err = NewNotFoundError(o.modelDefinition.ModelValue().Type().Name())
		return
	}

	// NOTE : Result is always a slice, so we need to get the first element.
	result = o.Payload().Result()[0]

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

func (o *builder[T]) FindAll() error {
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

func Use[T specs.Model](ctx context.Context, connector specs.Connector) specs.Builder[T] {
	var model T

	var builder specs.Builder[T] = &builder[T]{
		Context:   ctx,
		Connector: connector,

		model:           &model,
		modelDefinition: modeldefinition.Use(model).Parse(),
	}

	return builder
}
