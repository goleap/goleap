package specs

import "reflect"

type SchemaField interface {
	// Init allows you to initialise the field with its default value recursively (to avoid `nil`)
	Init()
	// Name returns the name of the field
	Name() string
	// Schema returns the schema of the field
	Schema() Schema
	// Tags returns the tags of the field
	Tags() map[string]string
	FromSchemaTypeList() []string
	RecursiveFullName() string
	Column() string
	Index() int

	Join() []DriverJoin
	Field() DriverField

	Value() reflect.Value

	Copy() any

	Set(value any)
	Get() any

	HasEmbeddedSchema() bool
	EmbeddedSchema() Schema

	IsPrimaryKey() bool
}
