package drivers

import (
	"github.com/lab210-dev/dbkit/specs"
	"strings"
)

type where struct {
	from     specs.DriverField
	operator string
	to       any
}

func (w *where) From() specs.DriverField {
	return w.from
}

func (w *where) Operator() string {
	return w.operator
}

func (w *where) To() any {
	return w.to
}

func (w *where) SetFrom(from specs.DriverField) specs.DriverWhere {
	w.from = from
	return w
}

func (w *where) SetOperator(operator string) specs.DriverWhere {
	w.operator = strings.TrimSpace(operator)
	return w
}

func (w *where) SetTo(to any) specs.DriverWhere {
	w.to = to
	return w
}

func NewWhere() specs.DriverWhere {
	return new(where)
}
