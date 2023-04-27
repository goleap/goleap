package specs

type DriverWhere interface {
	From() DriverField
	Operator() string
	To() any

	SetFrom(from DriverField) DriverWhere
	SetOperator(operator string) DriverWhere
	SetTo(to any) DriverWhere

	Formatted() (value string, args []any, err error)
}
