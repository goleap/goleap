package specs

type DriverField interface {
	Index() int
	Name() string
	NameInSchema() string

	SetIndex(index int) DriverField
	SetName(name string) DriverField
	SetNameInSchema(nameInSchema string) DriverField
}
