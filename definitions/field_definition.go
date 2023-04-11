package definitions

import (
	"fmt"
	"github.com/lab210-dev/dbkit/connector/drivers"
	"github.com/lab210-dev/dbkit/specs"
	"reflect"
	"strings"
	"sync"
)

type fieldDefinition struct {
	sync.Mutex
	name           string
	schema         specs.ModelDefinition
	embeddedSchema specs.ModelDefinition
	tags           map[string]string

	recursiveFullName string

	fieldType          reflect.Type
	fieldValue         reflect.Value
	fieldEmbeddedValue reflect.Value
	structField        reflect.StructField
	tag                reflect.StructTag
	index              int
	isSlice            bool
	init               bool

	visitedMap map[string]bool
}

func (field *fieldDefinition) Join() (joins []specs.DriverJoin) {
	if field.Model().FromField() != nil {
		if !field.IsSlice() {
			join := drivers.NewJoin().
				SetFrom(drivers.NewField().SetIndex(field.Model().FromField().Model().Index()).SetTable(field.Model().FromField().Model().TableName()).SetColumn(field.Model().FromField().Tags()["column"]).SetDatabase(field.Model().FromField().Model().DatabaseName())).
				SetTo(drivers.NewField().SetIndex(field.Model().Index()).SetTable(field.Model().TableName()).SetColumn(field.Model().FromField().Tags()["foreignKey"]).SetDatabase(field.Model().DatabaseName()))

			joins = append(joins, join)
		}
		joins = append(joins, field.Model().FromField().Join()...)
		return
	}
	return
}

func (field *fieldDefinition) Copy() any {
	return reflect.New(field.fieldType).Interface()
}

func (field *fieldDefinition) Value() reflect.Value {
	return field.fieldValue
}

func (field *fieldDefinition) IsSlice() bool {
	return field.isSlice
}

func (field *fieldDefinition) FromSchemaTypeList() (new []string) {
	if field.Model().FromField() != nil {
		new = append(new, field.Model().FromField().FromSchemaTypeList()...)
	}
	new = append(new, fmt.Sprintf("%v:%v", field.Model().ModelOrigin().Type(), field.IsSlice()))
	return
}

func (field *fieldDefinition) Model() specs.ModelDefinition {
	return field.schema
}

func (field *fieldDefinition) Get() any {
	return field.fieldValue.Interface()
}

func (field *fieldDefinition) Set(value any) {
	field.fieldValue.Set(reflect.ValueOf(value).Elem())

	field.Init()
}

func (field *fieldDefinition) Init() {
	if field.init {
		return
	}

	if field.Model().FromField() == nil {
		return
	}

	field.Model().FromField().Init()

	if field.Model().FromField().Value().Kind() == reflect.Ptr {
		cpy := reflect.New(field.Model().ModelValue().Type())
		cpy.Elem().Set(field.Model().ModelValue())

		field.Model().FromField().Value().Set(cpy)
		return
	}

	if field.Model().FromField().Value().Kind() == reflect.Slice {
		field.Model().FromField().Value().Set(reflect.Append(field.Model().FromField().Value(), field.Model().ModelValue()))
		return
	}

	field.Model().FromField().Value().Set(field.Model().ModelValue())
	field.init = true
}

func (field *fieldDefinition) Tags() map[string]string {
	return field.tags
}

func (field *fieldDefinition) EmbeddedSchema() specs.ModelDefinition {
	return field.embeddedSchema
}

func (field *fieldDefinition) SetEmbeddedSchema(embeddedSchema specs.ModelDefinition) specs.FieldDefinition {
	field.embeddedSchema = embeddedSchema
	return field
}

func (field *fieldDefinition) HasEmbeddedSchema() bool {
	return field.embeddedSchema != nil
}

func (field *fieldDefinition) SetIsSlice(isSlice bool) {
	field.isSlice = isSlice
}

func (field *fieldDefinition) FromSlice() bool {
	if field.Model().FromField() != nil {
		if field.Model().FromField().FromSlice() {
			return true
		}
	}
	return field.IsSlice()
}

func (md *modelDefinition) parseField(index int) specs.FieldDefinition {
	fieldStruct := md.modelType.Field(index)

	field := new(fieldDefinition)
	field.name = fieldStruct.Name
	field.schema = md
	field.fieldType = fieldStruct.Type
	field.structField = fieldStruct
	field.fieldValue = md.modelValue.Field(index)
	field.visitedMap = make(map[string]bool)

	field.tag = fieldStruct.Tag
	field.index = index

	field.ParseTags()

	if field.IsVisited() {
		return nil
	}

	field.RevealEmbeddedSchema()

	if field.IsSameSchemaFromField() {
		return nil
	}

	return field
}

func (field *fieldDefinition) ParseTags() {
	field.tags = make(map[string]string)
	// TODO (Lab210-dev) : add support to client choice of tag name
	tags := field.tag.Get("dbKit")

	for _, tag := range strings.Split(tags, ",") {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		tagParts := strings.Split(tag, ":")
		if len(tagParts) > 2 {
			continue
		}
		if len(tagParts) == 1 {
			tagParts = append(tagParts, "true")
		}
		field.tags[tagParts[0]] = tagParts[1]
	}
}

func (field *fieldDefinition) IsVisited() bool {
	field.Lock()
	defer field.Unlock()

	split := field.FromSchemaTypeList()
	countMap := make(map[string]int)
	for _, v := range split {
		countMap[v] = countMap[v] + 1
	}

	return countMap[fmt.Sprintf("%v:%v", field.Model().ModelOrigin().Type(), field.IsSlice())] > 2
}

func (field *fieldDefinition) Name() string {
	return field.name
}

func (field *fieldDefinition) Index() int {
	return field.schema.Index()
}

func (field *fieldDefinition) Column() string {
	return field.tags["column"]
}

func (field *fieldDefinition) IsPrimaryKey() bool {
	return field.tags["primaryKey"] == "true"
}

func (field *fieldDefinition) Field() specs.DriverField {
	return drivers.NewField().SetColumn(field.Column()).SetIndex(field.Index()).SetName(field.RecursiveFullName())
}

func (field *fieldDefinition) RecursiveFullName() string {
	// TODO (Lab210-dev) : maybe try to simplify this
	if field.recursiveFullName != "" {
		return field.recursiveFullName
	}

	if field.schema.FromField() == nil {
		field.recursiveFullName = field.Name()
		return field.recursiveFullName
	}

	field.recursiveFullName = fmt.Sprintf("%s.%s", field.schema.FromField().RecursiveFullName(), field.Name())
	return field.recursiveFullName
}

func (field *fieldDefinition) IsSameSchemaFromField() bool {
	return field.schema.FromField() != nil &&
		fmt.Sprintf("%s/%s", field.schema.FromField().Model().ModelValue().Type(), field.schema.FromField().Name()) == fmt.Sprintf("%s/%s", field.fieldEmbeddedValue.Type(), field.Name())
}

func (field *fieldDefinition) RevealEmbeddedSchema() specs.FieldDefinition {
	field.fieldEmbeddedValue = field.fieldValue

	if field.fieldEmbeddedValue.Kind() == reflect.Ptr {
		if field.fieldEmbeddedValue.IsNil() {
			field.fieldEmbeddedValue = reflect.New(field.fieldValue.Type().Elem())
		}
		field.fieldEmbeddedValue = field.fieldEmbeddedValue.Elem()
	}

	if field.fieldEmbeddedValue.Kind() == reflect.Slice {
		field.fieldEmbeddedValue = reflect.New(field.fieldValue.Type().Elem())
		field.fieldEmbeddedValue = field.fieldEmbeddedValue.Elem()
		field.SetIsSlice(true)
	}

	if field.fieldEmbeddedValue.Kind() != reflect.Struct {
		return nil
	}

	var model specs.Model
	var ok bool
	if model, ok = field.fieldEmbeddedValue.Addr().Interface().(specs.Model); !ok {
		return nil
	}

	if field.IsSameSchemaFromField() {
		return nil
	}

	embeddedSchema := Use(model).
		SetFromField(field).
		SetIndex(field.schema.Counter()).
		Parse()

	return field.SetEmbeddedSchema(embeddedSchema)
}
