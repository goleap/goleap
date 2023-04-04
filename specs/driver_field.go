package specs

type DriverField interface {
	Index() int
	Name() string
	NameInModel() string

	Column() (string, error)

	SetIndex(index int) DriverField
	SetName(name string) DriverField
	SetNameInModel(nameInSchema string) DriverField
}
