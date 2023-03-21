package dbkit

import "github.com/lab210-dev/dbkit/specs"

type condition struct {
	from     string
	operator string
	to       any
}

func (w *condition) From() string {
	return w.from
}

func (w *condition) Operator() string {
	return w.operator
}

func (w *condition) To() any {
	return w.to
}

func (w *condition) SetFrom(from string) specs.WhereCondition {
	w.from = from
	return w
}

func (w *condition) SetOperator(operator string) specs.WhereCondition {
	w.operator = operator
	return w
}

func (w *condition) SetTo(to any) specs.WhereCondition {
	w.to = to
	return w
}

func NewCondition() specs.WhereCondition {
	return new(condition)
}
