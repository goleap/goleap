package schema

import (
	"fmt"
	"github.com/goleap/goleap/connector/driver"
	"reflect"
	"strings"
	"sync"
)

type Field interface {
	Name() string
	Schema() Schema
	Tags() map[string]string
	fromSliceSchema(value string) string
	FromSchemaTypeList() []string
	VisitedMap() map[string]bool
	RecursiveFullName() string
	Column() string
	Index() int
	Key() string

	Join() []driver.Join
	Field() driver.Field

	Value() reflect.Value

	Copy() any

	Set(value any)
	Get() any

	HasEmbeddedSchema() bool
	EmbeddedSchema() Schema

	IsPrimaryKey() bool
}

type field struct {
	sync.Mutex
	name           string
	schema         Schema
	embeddedSchema Schema
	tags           map[string]string

	fieldType          reflect.Type
	fieldValue         reflect.Value
	fieldEmbeddedValue reflect.Value
	structField        reflect.StructField
	tag                reflect.StructTag
	index              int
	isSlice            bool

	visitedMap map[string]bool
}

func (field *field) Join() (joins []driver.Join) {
	if field.Schema().FromField() != nil {
		if !field.IsSlice() {
			join := driver.NewJoin().
				SetFromTable(field.Schema().TableName()).
				SetFromTableIndex(field.Schema().Index()).
				SetToTable(field.Schema().FromField().Schema().TableName()).
				SetToTableIndex(field.Schema().FromField().Schema().Index()).
				SetFromKey(field.Schema().FromField().Tags()["column"]).
				SetToKey(field.Schema().FromField().Tags()["foreignKey"])

			joins = append(joins, join)
		}
		joins = append(joins, field.Schema().FromField().Join()...)
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

func (field *field) fromSliceSchema(value string) (new string) {
	if field.IsSlice() {
		value = field.RecursiveFullName()
	}

	if field.Schema().FromField() != nil {
		tmp := field.Schema().FromField().fromSliceSchema(value)
		if tmp != "" {
			return tmp
		}
		return value
	}

	return value
}

func (field *field) FromSliceSchema() (new string) {
	return field.fromSliceSchema("")
}

func (field *field) FromSchemaTypeList() (new []string) {
	if field.Schema().FromField() != nil {
		new = append(new, field.Schema().FromField().FromSchemaTypeList()...)
	}
	new = append(new, fmt.Sprintf("%v:%v", field.Schema().ModelOrigin().Type(), field.IsSlice()))
	return
}

func (field *field) Schema() Schema {
	return field.schema
}

func (field *field) Get() any {
	return field.fieldValue.Addr().Interface()
}

func (field *field) Set(value any) {
	field.fieldValue.Set(reflect.ValueOf(value))

	if field.Schema().FromField() != nil {
		field.Schema().FromField().Value().Set(field.Schema().ModelValue())
	}
}

func (field *field) Tags() map[string]string {
	return field.tags
}

func (field *field) EmbeddedSchema() Schema {
	return field.embeddedSchema
}

func (field *field) SetEmbeddedSchema(embeddedSchema Schema) Field {
	if embeddedSchema == nil {
		return field
	}

	field.embeddedSchema = embeddedSchema
	return field
}

func (field *field) HasEmbeddedSchema() bool {
	return field.embeddedSchema != nil
}

func (field *field) VisitedMap() map[string]bool {
	if field.schema.FromField() == nil {
		return field.visitedMap
	}

	return field.schema.FromField().VisitedMap()
}

func (field *field) SetIsSlice(isSlice bool) {
	field.isSlice = isSlice
}

func (schema *schema) parseField(index int) Field {
	fieldStruct := schema.modelType.Field(index)

	field := new(field)
	field.name = fieldStruct.Name
	field.schema = schema
	field.fieldType = fieldStruct.Type
	field.structField = fieldStruct
	field.fieldValue = schema.modelValue.Field(index)
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
	// TODO: add support to client choice of tag name
	tags := field.tag.Get("goleap")

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

	visitedMap := field.VisitedMap()
	state := visitedMap[field.Key()]

	/*if state {
		log.Print(field.RecursiveFullName())
		log.Print("field already visited: ", field.Key())
		log.Println(field.FromSchemaTypeList())
		log.Println()
	}*/

	return state
}

func (field *field) Visited() {
	field.Lock()
	defer field.Unlock()

	visitedMap := field.VisitedMap()

	/*	log.Print("--->", field.RecursiveFullName())
		log.Print("visited: ", field.Key())
		log.Print(field.FromSchemaTypeList())
		log.Println()*/
	visitedMap[field.Key()] = true
}

func (field *field) Key() string {
	split := field.FromSchemaTypeList()
	for i := 0; i < len(split); i += 2 {
		end := i + 2
		if end > len(split) {
			end = len(split)
		}
		split[i] = strings.Join(split[i:end], ".")
	}

	unique := make(map[string]bool)
	for _, word := range split {
		unique[word] = true
	}

	schemaListLength := len(unique)
	key := fmt.Sprintf("%s@%s::%s.%s-%d", field.schema.DatabaseName(), field.schema.TableName(), field.Name(), field.tags["column"], schemaListLength)

	if field.schema.FromField() != nil {
		key = fmt.Sprintf("%s@%s::%s.%s/%s@%s::%s.%s-%d",
			field.schema.FromField().Schema().DatabaseName(),
			field.schema.FromField().Schema().TableName(),
			field.schema.FromField().Name(),
			field.schema.FromField().Tags()["column"],
			field.schema.DatabaseName(),
			field.schema.TableName(),
			field.Name(),
			field.Column(),
			schemaListLength,
		)
	}

	return strings.ToLower(key)
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

func (field *field) Field() driver.Field {
	return driver.NewField().SetName(field.Column()).SetIndex(field.Index())
}

func (field *field) RecursiveNameTest() string {
	if field.schema.FromField() == nil {
		return "nil"
	}
	return field.schema.FromField().Name()
}

func (field *field) RecursiveName() string {
	if field.schema.FromField() == nil {
		return field.Name()
	}
	return field.schema.FromField().RecursiveFullName()
}

func (field *field) RecursiveFullName() string {
	if field.schema.FromField() == nil {
		return field.Name()
	}
	return fmt.Sprintf("%s.%s", field.schema.FromField().RecursiveFullName(), field.Name())
}

func (field *field) IsSameSchemaFromField() bool {
	return field.schema.FromField() != nil &&
		fmt.Sprintf("%s/%s", field.schema.FromField().Schema().ModelValue().Type(), field.schema.FromField().Name()) == fmt.Sprintf("%s/%s", field.fieldEmbeddedValue.Type(), field.Name())
}

func (field *field) RevealEmbeddedSchema() Field {
	field.fieldEmbeddedValue = field.fieldValue

	if field.fieldEmbeddedValue.Kind() == reflect.Ptr {
		if field.fieldEmbeddedValue.IsNil() {
			field.fieldEmbeddedValue = reflect.New(field.fieldValue.Type().Elem())
			field.fieldEmbeddedValue = field.fieldEmbeddedValue.Elem()
		}
	}

	if field.fieldEmbeddedValue.Kind() == reflect.Slice {
		field.fieldEmbeddedValue = reflect.New(field.fieldValue.Type().Elem())
		field.fieldEmbeddedValue = field.fieldEmbeddedValue.Elem()
		field.SetIsSlice(true)
	}

	if field.fieldEmbeddedValue.Kind() != reflect.Struct {
		return nil
	}

	var model Model
	var ok bool
	if model, ok = field.fieldEmbeddedValue.Interface().(Model); !ok {
		return nil
	}

	field.Visited()

	if field.IsSameSchemaFromField() {
		return nil
	}

	embeddedSchema := New(model).
		SetFromField(field).
		SetIndex(field.schema.Counter()).
		Parse()

	/*if !field.IsSlice() {
		field.fieldValue.Set(embeddedSchema.ModelValue())
	}*/

	return field.SetEmbeddedSchema(embeddedSchema)
}
