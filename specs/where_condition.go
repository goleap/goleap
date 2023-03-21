package specs

type WhereCondition interface {
	From() string
	Operator() string
	To() any

	SetFrom(from string) WhereCondition
	SetOperator(operator string) WhereCondition
	SetTo(to any) WhereCondition
}
