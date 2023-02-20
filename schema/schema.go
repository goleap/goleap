package schema

import (
	"reflect"
	"sync"
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
	Get() Model
	Copy() Schema
}

type schema struct {
	Model
	sync.Mutex

	index       int
	counter     int
	parsed      bool
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

func (schema *schema) Copy() Schema {
	tmp := reflect.New(reflect.TypeOf(schema)).Elem()
	tmp.Set(reflect.ValueOf(schema))

	return tmp.Interface().(Schema)
}

func New(model Model) Schema {
	schema := new(schema)
	schema.Model = model

	schema.modelType = reflect.TypeOf(model)
	schema.modelValue = reflect.ValueOf(model)
	schema.modelOrigin = schema.modelValue

	if schema.modelType.Kind() == reflect.Struct {
		schema.modelValue = reflect.New(schema.modelType).Elem().Addr()
		schema.modelType = schema.modelValue.Type()
	}

	if schema.modelType.Kind() == reflect.Ptr {
		schema.modelType = schema.modelType.Elem()
		schema.modelValue = schema.modelValue.Elem()
	}

	schema.fieldByName = make(map[string]Field)

	return schema
}

func (schema *schema) Get() Model {
	return schema.modelValue.Interface().(Model)
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
	defer schema.Unlock()
	schema.Lock()

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
	// no need to parse again normally...
	if schema.parsed {
		return schema
	}

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

	// TODO use a setter to set this value
	schema.parsed = true

	return schema
}
