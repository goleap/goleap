package specs

import "reflect"

type ModelDefinition interface {
	Model

	Fields() []ModelField
	FieldByName() map[string]ModelField

	GetFieldByName(name string) ModelField

	SetFromField(fromField ModelField) ModelDefinition
	FromField() ModelField

	SetIndex(index int) ModelDefinition
	Index() int
	Counter() int

	Parse() ModelDefinition

	ModelValue() reflect.Value
	ModelOrigin() reflect.Value

	GetPrimaryKeyField() ModelField
	Get() Model
	Copy() ModelDefinition
}
