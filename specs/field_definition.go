package specs

import "reflect"

// FieldDefinition is the interface that describes a field
type FieldDefinition interface {
	// Init allows you to initialise the field with its default value recursively (to avoid `nil`)
	Init()
	// Name returns the name of the field
	Name() string
	// Model returns the schema of the field
	Model() ModelDefinition
	// Tags returns the tags of the field
	Tags() map[string]string
	FromSchemaTypeList() []string
	RecursiveFullName() string
	FundamentalName() string
	Column() string
	ForeignKey() string
	Index() int

	GetByColumn() (FieldDefinition, error)
	GetToColumn() (FieldDefinition, error)

	Join() []DriverJoin
	Field() DriverField

	Value() reflect.Value

	Copy() any

	Set(value any)
	Get() any

	HasEmbeddedSchema() bool
	EmbeddedSchema() ModelDefinition

	IsSlice() bool
	FromSlice() bool

	IsPrimaryKey() bool
}
