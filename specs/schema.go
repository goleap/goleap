package specs

import "reflect"

type Schema interface {
	Model

	Fields() []SchemaField
	FieldByName() map[string]SchemaField

	GetFieldByName(name string) SchemaField

	SetFromField(fromField SchemaField) Schema
	FromField() SchemaField

	SetIndex(index int) Schema
	Index() int
	Counter() int

	Parse() Schema

	ModelValue() reflect.Value
	ModelOrigin() reflect.Value

	GetPrimaryKeyField() SchemaField
	Get() Model
	Copy() Schema
}
