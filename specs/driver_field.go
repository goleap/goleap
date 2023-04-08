package specs

type DriverField interface {
	Index() int
	Column() string
	Database() string
	Table() string
	Name() string

	SetIndex(index int) DriverField
	SetColumn(name string) DriverField
	SetTable(name string) DriverField
	SetDatabase(name string) DriverField
	SetName(name string) DriverField

	SetFn(fn string, args []DriverField) DriverField

	Formatted() (string, error)
}
