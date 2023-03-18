package dbkit

import (
	"context"
	"fmt"
	"github.com/lab210-dev/dbkit/connector/drivers"
	"github.com/lab210-dev/dbkit/connector/drivers/operators"
	"github.com/lab210-dev/dbkit/schema"
	"github.com/lab210-dev/dbkit/specs"
)

type DbKit[T specs.Model] interface {
	Get(primaryKey any) (T, error)
	Delete(primaryKey any) error

	Create() (err error)
	Update() error

	Find() error
	FindAll() error

	Fields(field ...string) DbKit[T]
	Where(condition specs.Condition) DbKit[T]
	Limit(limit int) DbKit[T]
	Offset(offset int) DbKit[T]
	OrderBy(fields ...string) DbKit[T]

	Count() (total int64, err error)

	Payload() specs.PayloadAugmented[T]
}

type dbKit[T specs.Model] struct {
	context.Context
	specs.Connector

	model  *T
	schema specs.Schema
	fields []string

	focusedSchemaFields []specs.SchemaField

	driverFields []specs.DriverField
	driverJoins  []specs.DriverJoin
	driverWheres []specs.DriverWhere

	focusedSchemaFieldsCopy []any
}

func (o *dbKit[T]) countFocusedSchemaFields() int {
	return len(o.focusedSchemaFields)
}

func (o *dbKit[T]) buildFieldsFromSchema() (err error) {
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

func (o *dbKit[T]) getDriverFields() []specs.DriverField {

	if len(o.driverFields) > 0 {
		return o.driverFields
	}

	for _, field := range o.focusedSchemaFields {
		o.driverFields = append(o.driverFields, field.Field())
	}

	return o.driverFields
}

func (o *dbKit[T]) getDriverJoins() []specs.DriverJoin {

	if len(o.driverJoins) > 0 {
		return o.driverJoins
	}

	for _, field := range o.focusedSchemaFields {
		o.driverJoins = append(o.driverJoins, field.Join()...)
	}

	return o.driverJoins
}

func (o *dbKit[T]) getDriverWheres() []specs.DriverWhere {
	return o.driverWheres
}

func (o *dbKit[T]) Payload() specs.PayloadAugmented[T] {
	return NewPayload[T](o.schema, o.getDriverFields(), o.getDriverJoins(), o.getDriverWheres())
}

func (o *dbKit[T]) Get(primaryKeyValue any) (result T, err error) {

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
		err = &FieldNotFoundError{message: fmt.Sprintf("no fields selected")}
		return
	}

	getPayload := o.Payload()
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

func (o *dbKit[T]) Delete(primaryKeyValue any) error {
	//TODO implement me
	panic("implement me")
}

func (o *dbKit[T]) Create() (err error) {
	//TODO implement me
	panic("implement me")
}

func (o *dbKit[T]) Update() error {
	//TODO implement me
	panic("implement me")
}

func (o *dbKit[T]) Find() error {
	//TODO implement me
	panic("implement me")
}

func (o *dbKit[T]) FindAll() error {
	//TODO implement me
	panic("implement me")
}

func (o *dbKit[T]) Fields(field ...string) DbKit[T] {
	o.fields = field
	return o
}

func (o *dbKit[T]) Where(condition specs.Condition) DbKit[T] {
	//TODO implement me
	panic("implement me")
}

func (o *dbKit[T]) Limit(limit int) DbKit[T] {
	//TODO implement me
	panic("implement me")
}

func (o *dbKit[T]) Offset(offset int) DbKit[T] {
	//TODO implement me
	panic("implement me")
}

func (o *dbKit[T]) OrderBy(fields ...string) DbKit[T] {
	//TODO implement me
	panic("implement me")
}

func (o *dbKit[T]) Count() (total int64, err error) {
	//TODO implement me
	panic("implement me")
}

func Use[T specs.Model](ctx context.Context, connector specs.Connector) DbKit[T] {
	var model T

	return &dbKit[T]{
		Context:   ctx,
		Connector: connector,

		model:  &model,
		schema: schema.New(model).Parse(),
	}
}
