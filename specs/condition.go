package specs

type Condition interface {
	From() string
	Operator() string
	To() any

	SetFrom(from string) Condition
	SetOperator(operator string) Condition
	SetTo(to any) Condition
}
