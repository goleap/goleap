package dbkit

import "github.com/kitstack/dbkit/specs"

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

func (w *condition) SetFrom(from string) specs.Condition {
	w.from = from
	return w
}

func (w *condition) SetOperator(operator string) specs.Condition {
	w.operator = operator
	return w
}

func (w *condition) SetTo(to any) specs.Condition {
	w.to = to
	return w
}

func NewCondition() specs.Condition {
	return new(condition)
}
