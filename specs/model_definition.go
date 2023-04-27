package specs

import "reflect"

type UseModelDefinition func(model Model) ModelDefinition

type ModelDefinition interface {
	Model

	Fields() []FieldDefinition
	FieldByName() map[string]FieldDefinition

	GetFieldByName(name string) (FieldDefinition, ErrNotFoundError)
	GetPrimaryField() (FieldDefinition, ErrPrimaryFieldNotFound)
	GetFieldByColumn(column string) (FieldDefinition, ErrFieldNoFoundByColumn)

	SetFromField(fromField FieldDefinition) ModelDefinition
	FromField() FieldDefinition

	SetIndex(index int) ModelDefinition
	Index() int
	Counter() int

	Parse() ModelDefinition

	TypeName() string

	ModelValue() reflect.Value
	ModelOrigin() reflect.Value

	Copy() Model
}
