package specs

import "reflect"

type ModelDefinition interface {
	Model

	Fields() []FieldDefinition
	FieldByName() map[string]FieldDefinition

	GetFieldByName(name string) (FieldDefinition, FieldNotFoundError)
	GetPrimaryField() (FieldDefinition, PrimaryFieldNotFoundError)
	GetFieldByColumn(column string) (FieldDefinition, error)

	SetFromField(fromField FieldDefinition) ModelDefinition
	FromField() FieldDefinition

	SetIndex(index int) ModelDefinition
	Index() int
	Counter() int

	Parse() ModelDefinition

	ModelValue() reflect.Value
	ModelOrigin() reflect.Value

	Copy() Model
}
