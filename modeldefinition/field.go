package modeldefinition

import (
	"fmt"
	"github.com/lab210-dev/dbkit/connector/drivers"
	"github.com/lab210-dev/dbkit/specs"
	"reflect"
	"strings"
	"sync"
)

type field struct {
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

func (field *field) Join() (joins []specs.DriverJoin) {
	if field.Model().FromField() != nil {
		if !field.IsSlice() {
			join := drivers.NewJoin().
				SetFromTableIndex(field.Model().Index()).
				SetToTable(field.Model().FromField().Model().TableName()).
				SetToTableIndex(field.Model().FromField().Model().Index()).
				SetFromKey(field.Model().FromField().Tags()["column"]).
				SetToKey(field.Model().FromField().Tags()["foreignKey"])

			joins = append(joins, join)
		}
		joins = append(joins, field.Model().FromField().Join()...)
		return
	}
	return
}

func (field *field) Copy() any {
	return reflect.New(field.fieldType).Interface()
}

func (field *field) Value() reflect.Value {
	return field.fieldValue
}

func (field *field) IsSlice() bool {
	return field.isSlice
}

func (field *field) FromSchemaTypeList() (new []string) {
	if field.Model().FromField() != nil {
		new = append(new, field.Model().FromField().FromSchemaTypeList()...)
	}
	new = append(new, fmt.Sprintf("%v:%v", field.Model().ModelOrigin().Type(), field.IsSlice()))
	return
}

func (field *field) Model() specs.ModelDefinition {
	return field.schema
}

func (field *field) Get() any {
	return field.fieldValue.Interface()
}

func (field *field) Set(value any) {
	field.fieldValue.Set(reflect.ValueOf(value).Elem())

	field.Init()
}

func (field *field) Init() {
	if field.init {
		return
	}

	if field.Model().FromField() == nil {
		return
	}

	field.Model().FromField().Init()

	if field.Model().FromField().Value().Kind() == reflect.Ptr {
		field.Model().FromField().Value().Set(field.Model().ModelValue().Addr())
		return
	}

	if field.Model().FromField().Value().Kind() == reflect.Slice {
		field.Model().FromField().Value().Set(reflect.Append(field.Model().FromField().Value(), field.Model().ModelValue()))
		return
	}

	field.Model().FromField().Value().Set(field.Model().ModelValue())
	field.init = true
}

func (field *field) Tags() map[string]string {
	return field.tags
}

func (field *field) EmbeddedSchema() specs.ModelDefinition {
	return field.embeddedSchema
}

func (field *field) SetEmbeddedSchema(embeddedSchema specs.ModelDefinition) specs.ModelField {
	field.embeddedSchema = embeddedSchema
	return field
}

func (field *field) HasEmbeddedSchema() bool {
	return field.embeddedSchema != nil
}

func (field *field) SetIsSlice(isSlice bool) {
	field.isSlice = isSlice
}

func (md *modelDefinition) parseField(index int) specs.ModelField {
	fieldStruct := md.modelType.Field(index)

	field := new(field)
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

func (field *field) ParseTags() {
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

func (field *field) IsVisited() bool {
	field.Lock()
	defer field.Unlock()

	split := field.FromSchemaTypeList()
	countMap := make(map[string]int)
	for _, v := range split {
		countMap[v] = countMap[v] + 1
	}

	return countMap[fmt.Sprintf("%v:%v", field.Model().ModelOrigin().Type(), field.IsSlice())] > 2
}

func (field *field) Name() string {
	return field.name
}

func (field *field) Index() int {
	return field.schema.Index()
}

func (field *field) Column() string {
	return field.tags["column"]
}

func (field *field) IsPrimaryKey() bool {
	return field.tags["primaryKey"] == "true"
}

func (field *field) Field() specs.DriverField {
	return drivers.NewField().SetName(field.Column()).SetIndex(field.Index()).SetNameInSchema(field.RecursiveFullName())
}

func (field *field) RecursiveFullName() string {
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

func (field *field) IsSameSchemaFromField() bool {
	return field.schema.FromField() != nil &&
		fmt.Sprintf("%s/%s", field.schema.FromField().Model().ModelValue().Type(), field.schema.FromField().Name()) == fmt.Sprintf("%s/%s", field.fieldEmbeddedValue.Type(), field.Name())
}

func (field *field) RevealEmbeddedSchema() specs.ModelField {
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
