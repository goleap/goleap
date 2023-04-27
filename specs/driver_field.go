package specs

// DriverField is the interface that wraps the basic methods of a field.
type DriverField interface {
	Index() int
	Column() string
	Database() string
	Table() string
	Name() string

	IsCustom() bool

	SetIndex(index int) DriverField
	SetColumn(name string) DriverField
	SetTable(name string) DriverField
	SetDatabase(name string) DriverField
	SetName(name string) DriverField

	SetCustom(fn string, args []DriverField) DriverField

	Formatted() (string, error)
}
