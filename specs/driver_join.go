package specs

type JoinMethod int

type DriverJoin interface {
	Validate() error

	Method() string
	From() DriverField
	To() DriverField

	SetMethod(method JoinMethod) DriverJoin
	SetFrom(field DriverField) DriverJoin
	SetTo(field DriverField) DriverJoin

	Formatted() (string, error)
}
