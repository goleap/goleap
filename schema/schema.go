package schema

import (
	"reflect"
)

type Schema interface {
	Model

	Fields() []Field
	FieldByName() map[string]Field

	GetFieldByName(name string) Field

	SetFromField(fromField Field) Schema
	FromField() Field

	SetIndex(index int) Schema
	Index() int
	Counter() int

	Parse() Schema

	ModelValue() reflect.Value
	ModelOrigin() reflect.Value

	GetPrimaryKeyField() Field
}

type schema struct {
	Model

	index       int
	counter     int
	modelOrigin reflect.Value
	modelType   reflect.Type
	modelValue  reflect.Value

	fields      []Field
	fieldByName map[string]Field

	fromField Field
}

func (schema *schema) ModelOrigin() reflect.Value {
	return schema.modelOrigin
}

func (schema *schema) ModelValue() reflect.Value {
	return schema.modelValue
}

func New(model Model) Schema {
	schema := new(schema)
	schema.Model = model

	schema.modelType = reflect.TypeOf(model)
	schema.modelValue = reflect.ValueOf(model)
	schema.modelOrigin = schema.modelValue

	if !schema.modelValue.CanAddr() {
		schema.modelValue = reflect.New(schema.modelType)
		schema.modelValue = schema.modelValue.Elem()
		schema.modelType = schema.modelValue.Type()
	}

	if schema.modelType.Kind() == reflect.Ptr {
		schema.modelType = schema.modelType.Elem()
		schema.modelValue = schema.modelValue.Elem()
	}

	schema.fieldByName = make(map[string]Field)

	return schema
}

func (schema *schema) SetIndex(index int) Schema {
	schema.index = index
	schema.counter = index
	return schema
}

func (schema *schema) GetPrimaryKeyField() Field {
	for _, field := range schema.fields {

		if field.Schema() != schema {
			continue
		}

		if field.IsPrimaryKey() {
			return field
		}

	}
	return nil
}

func (schema *schema) Index() int {
	return schema.index
}

func (schema *schema) Counter() int {
	if schema.FromField() != nil {
		return schema.FromField().Schema().Counter()
	}

	schema.counter++
	return schema.counter
}

func (schema *schema) Fields() []Field {
	return schema.fields
}

func (schema *schema) FieldByName() map[string]Field {
	return schema.fieldByName
}

func (schema *schema) AddField(field Field) {
	if field.HasEmbeddedSchema() {
		schema.fields = append(schema.fields, field.EmbeddedSchema().Fields()...)

		for key, value := range field.EmbeddedSchema().FieldByName() {
			schema.fieldByName[key] = value
		}
		return
	}

	schema.fields = append(schema.fields, field)
	schema.fieldByName[field.RecursiveFullName()] = field
}

func (schema *schema) FromField() Field {
	return schema.fromField
}

func (schema *schema) SetFromField(fromField Field) Schema {
	schema.fromField = fromField
	return schema
}

func (schema *schema) GetFieldByName(name string) Field {
	return schema.fieldByName[name]
}

func (schema *schema) Parse() Schema {
	for i := 0; i < schema.modelType.NumField(); i++ {
		fieldStruct := schema.modelType.Field(i)
		if !fieldStruct.IsExported() {
			continue
		}

		field := schema.parseField(i)
		if field == nil {
			continue
		}

		schema.AddField(field)
	}

	return schema
}
