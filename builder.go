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

	model           *T
	modelDefinition specs.ModelDefinition
	fields          []string

	focusedSchemaFields []specs.FieldDefinition

	driverFields []specs.DriverField
	driverJoins  []specs.DriverJoin
	driverWheres []specs.DriverWhere

	payload specs.PayloadAugmented[T]
}

func (o *builder[T]) countFocusedSchemaFields() int {
	return len(o.focusedSchemaFields)
}

func (o *builder[T]) buildFieldsFromModelDefinition() (err error) {
	for _, fieldName := range o.fields {
		field, err := o.modelDefinition.GetFieldByName(fieldName)

		if err != nil {
			return err
		}

		o.focusedSchemaFields = append(o.focusedSchemaFields, field)
	}

	return
}

func (o *builder[T]) getDriverFields() []specs.DriverField {

	if len(o.driverFields) > 0 {
		return o.driverFields
	}

	for _, field := range o.focusedSchemaFields {
		o.driverFields = append(o.driverFields, field.Field())
	}

	return o.driverFields
}

func (o *builder[T]) getDriverJoins() []specs.DriverJoin {

	if len(o.driverJoins) > 0 {
		return o.driverJoins
	}

	for _, field := range o.focusedSchemaFields {
		o.driverJoins = append(o.driverJoins, field.Join()...)
	}

	return o.driverJoins
}

func (o *builder[T]) getDriverWheres() []specs.DriverWhere {
	return o.driverWheres
}

func (o *builder[T]) buildPayload() specs.PayloadAugmented[T] {
	o.payload = NewPayload[T]()
	o.payload.SetFields(o.getDriverFields())
	o.payload.SetJoins(o.getDriverJoins())
	o.payload.SetWheres(o.getDriverWheres())

	return o.payload
}

func (o *builder[T]) Payload() specs.PayloadAugmented[T] {
	return o.payload
}

func (o *builder[T]) Get(primaryKeyValue any) (result T, err error) {

	primaryKeyField, err := o.modelDefinition.GetPrimaryField()
	if err != nil {
		return
	}

	err = o.buildFieldsFromModelDefinition()
	if err != nil {
		return
	}

	o.driverWheres = append(o.driverWheres, drivers.NewWhere().SetFrom(primaryKeyField.Field()).SetOperator(operators.Equal).SetTo(primaryKeyValue))

	if o.countFocusedSchemaFields() == 0 {
		// TODO (Lab210-dev) : Factory for error. TIP SchemaError, SchemaFieldError, etc.
		// err = &errors.FieldNotFoundError{message: "no fields selected"}
		return
	}

	getPayload := o.buildPayload()

	err = o.Connector.Select(o.Context, getPayload)
	if err != nil {
		return
	}

	// TODO (Lab210-dev) : Do we throw a "not found" error if the result is empty?
	if len(getPayload.Result()) == 0 {
		return
	}

	// NOTE : Result is always a slice, so we need to get the first element.
	result = getPayload.Result()[0]

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

func (o *builder[T]) Find() error {
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

func (o *builder[T]) Where(_ specs.WhereCondition) specs.Builder[T] {
	//TODO implement me
	panic("implement me")
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
