package specs

type DriverField interface {
	Index() int
	Name() string

	SetIndex(index int) DriverField
	SetName(name string) DriverField
}
