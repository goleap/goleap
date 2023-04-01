package specs

type DriverLimit interface {
	Offset() int
	Limit() int

	SetOffset(index int) DriverLimit
	SetLimit(index int) DriverLimit
}
