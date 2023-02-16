package driver

var (
	EqualOperator = "="
)

type where struct {
	from     Field
	operator string
	to       any
}

type Where interface {
	From() Field
	Operator() string
	To() any

	SetFrom(from Field) Where
	SetOperator(operator string) Where
	SetTo(to any) Where
}

func (w *where) From() Field {
	return w.from
}

func (w *where) Operator() string {
	return w.operator
}

func (w *where) To() any {
	return w.to
}

func (w *where) SetFrom(from Field) Where {
	w.from = from
	return w
}

func (w *where) SetOperator(operator string) Where {
	w.operator = operator
	return w
}

func (w *where) SetTo(to any) Where {
	w.to = to
	return w
}

func NewWhere() Where {
	return new(where)
}
