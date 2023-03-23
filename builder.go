package dbkit

import (
	"context"
	"fmt"
	"github.com/lab210-dev/dbkit/connector/drivers"
	"github.com/lab210-dev/dbkit/connector/drivers/operators"
	"github.com/lab210-dev/dbkit/schema"
	"github.com/lab210-dev/dbkit/specs"
)

type Builder[T specs.Model] interface {
	Get(primaryKey any) (T, error)
	Delete(primaryKey any) error

	Create() (err error)
	Update() error

	Find() error
	FindAll() error

	Fields(field ...string) Builder[T]
	Where(condition specs.WhereCondition) Builder[T]
	Limit(limit int) Builder[T]
	Offset(offset int) Builder[T]
	OrderBy(fields ...string) Builder[T]

	Count() (total int64, err error)

	Payload() specs.PayloadAugmented[T]
}

type builder[T specs.Model] struct {
	context.Context
	specs.Connector

	model  *T
	schema specs.Schema
	fields []string

	focusedSchemaFields []specs.SchemaField

	driverFields []specs.DriverField
	driverJoins  []specs.DriverJoin
	driverWheres []specs.DriverWhere

	payload specs.PayloadAugmented[T]

	// focusedSchemaFieldsCopy []any
}

func (o *builder[T]) countFocusedSchemaFields() int {
	return len(o.focusedSchemaFields)
}

func (o *builder[T]) buildFieldsFromSchema() (err error) {
	for _, fieldName := range o.fields {
		field := o.schema.GetFieldByName(fieldName)

		if field == nil {
			// TODO (Lab210-dev) Use a real name of schema.
			err = &FieldNotFoundError{message: fmt.Sprintf("field `%s` not found in schema `%s`", fieldName, o.schema.TableName()), field: fieldName}
			return
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

	o.driverWheres = append(o.driverWheres, drivers.NewWhere().SetFrom(o.schema.GetPrimaryKeyField().Field()).SetOperator(operators.Equal).SetTo(primaryKeyValue))

	primaryKeyField := o.schema.GetPrimaryKeyField()
	if primaryKeyField == nil {
		err = &FieldNotFoundError{message: fmt.Sprintf("primary key not found in schema `%s`", o.schema.TableName())}
		return
	}

	err = o.buildFieldsFromSchema()
	if err != nil {
		return
	}

	if o.countFocusedSchemaFields() == 0 {
		// TODO (Lab210-dev) : Factory for error. TIP SchemaError, SchemaFieldError, etc.
		err = &FieldNotFoundError{message: "no fields selected"}
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

func (o *builder[T]) Delete(primaryKeyValue any) error {
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

func (o *builder[T]) Fields(field ...string) Builder[T] {
	o.fields = field
	return o
}

func (o *builder[T]) Where(condition specs.WhereCondition) Builder[T] {
	//TODO implement me
	panic("implement me")
}

func (o *builder[T]) Limit(limit int) Builder[T] {
	//TODO implement me
	panic("implement me")
}

func (o *builder[T]) Offset(offset int) Builder[T] {
	//TODO implement me
	panic("implement me")
}

func (o *builder[T]) OrderBy(fields ...string) Builder[T] {
	//TODO implement me
	panic("implement me")
}

func (o *builder[T]) Count() (total int64, err error) {
	//TODO implement me
	panic("implement me")
}

func Use[T specs.Model](ctx context.Context, connector specs.Connector) Builder[T] {
	var model T

	var builder Builder[T] = &builder[T]{
		Context:   ctx,
		Connector: connector,

		model:  &model,
		schema: schema.New(model).Parse(),
	}

	return builder
}
