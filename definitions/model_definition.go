package definitions

import (
	"github.com/kitstack/dbkit/specs"
	"reflect"
	"sync"
)

type modelDefinition struct {
	specs.Model
	sync.Mutex

	index       int
	counter     int
	parsed      bool
	modelOrigin reflect.Value
	modelType   reflect.Type
	modelValue  reflect.Value

	fields      []specs.FieldDefinition
	fieldByName map[string]specs.FieldDefinition

	fromField specs.FieldDefinition
}

func (md *modelDefinition) ModelOrigin() reflect.Value {
	return md.modelOrigin
}

func (md *modelDefinition) ModelValue() reflect.Value {
	return md.modelValue
}

func Use(model specs.Model) specs.ModelDefinition {
	schema := new(modelDefinition)
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

		if schema.modelValue.IsNil() {
			schema.modelValue = reflect.New(schema.modelType)
			schema.Model = schema.modelValue.Interface().(specs.Model)
		}

		schema.modelValue = schema.modelValue.Elem()
	}

	schema.fieldByName = make(map[string]specs.FieldDefinition)

	return schema
}

func (md *modelDefinition) Copy() specs.Model {
	tmp := reflect.New(md.modelType)
	tmp.Elem().Set(md.modelValue)
	return tmp.Interface().(specs.Model)
}

func (md *modelDefinition) SetIndex(index int) specs.ModelDefinition {
	md.index = index
	md.counter = index
	return md
}

func (md *modelDefinition) GetPrimaryField() (specs.FieldDefinition, specs.ErrPrimaryFieldNotFound) {
	for _, field := range md.fields {

		if field.Model() != md {
			continue
		}

		if field.IsPrimaryKey() {
			return field, nil
		}

	}
	return nil, NewErrNoPrimaryField(md)
}

func (md *modelDefinition) GetFieldByColumn(column string) (specs.FieldDefinition, specs.ErrFieldNoFoundByColumn) {
	for _, field := range md.fields {

		if field.Model() != md {
			continue
		}

		if field.Column() == column {
			return field, nil
		}

	}
	return nil, NewErrFieldNoFoundByColumn(column, md)
}

func (md *modelDefinition) Index() int {
	return md.index
}

func (md *modelDefinition) Counter() int {
	if md.FromField() != nil {
		return md.FromField().Model().Counter()
	}

	md.counter++
	return md.counter
}

func (md *modelDefinition) Fields() []specs.FieldDefinition {
	return md.fields
}

func (md *modelDefinition) FieldByName() map[string]specs.FieldDefinition {
	return md.fieldByName
}

func (md *modelDefinition) AddField(field specs.FieldDefinition) {
	defer md.Unlock()
	md.Lock()

	if field.HasEmbeddedSchema() {
		md.fields = append(md.fields, field.EmbeddedSchema().Fields()...)

		for key, value := range field.EmbeddedSchema().FieldByName() {
			md.fieldByName[key] = value
		}
		return
	}

	md.fields = append(md.fields, field)
	md.fieldByName[field.RecursiveFullName()] = field
}

func (md *modelDefinition) FromField() specs.FieldDefinition {
	return md.fromField
}

func (md *modelDefinition) SetFromField(fromField specs.FieldDefinition) specs.ModelDefinition {
	md.fromField = fromField
	return md
}

func (md *modelDefinition) GetFieldByName(name string) (specs.FieldDefinition, specs.ErrNotFoundError) {
	md.Lock()
	defer md.Unlock()

	if field, ok := md.fieldByName[name]; ok {
		return field, nil
	}

	return nil, NewErrFieldNotFound(name, md)
}

func (md *modelDefinition) Parse() specs.ModelDefinition {
	// no need to parse again normally...
	if md.parsed {
		return md
	}

	for i := 0; i < md.modelType.NumField(); i++ {
		fieldStruct := md.modelType.Field(i)
		if !fieldStruct.IsExported() {
			continue
		}

		field := md.parseField(i)
		if field == nil {
			continue
		}

		md.AddField(field)
	}

	// TODO (kitstack) : use a setter to set this value
	md.parsed = true

	return md
}
