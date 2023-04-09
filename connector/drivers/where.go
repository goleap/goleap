package drivers

import (
	"fmt"
	"github.com/lab210-dev/dbkit/connector/drivers/operators"
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

func (w *where) buildOperator() (string, error) {
	switch w.Operator() {
	case operators.Equal, operators.NotEqual:
		return fmt.Sprintf("%s ?", w.Operator()), nil
	case operators.In, operators.NotIn:
		return fmt.Sprintf("%s (?)", w.Operator()), nil
	case operators.IsNull, operators.IsNotNull:
		return w.Operator(), nil
	}
	return "", NewUnknownOperatorErr(w.Operator())
}

func (w *where) Formatted() (string, any, error) {
	operator, err := w.buildOperator()
	if err != nil {
		return "", nil, err
	}

	from, err := w.From().Formatted()
	if err != nil {
		return "", nil, err
	}

	drvField, ok := w.To().(specs.DriverField)
	if ok {
		to, err := drvField.Formatted()
		if err != nil {
			return "", nil, err
		}
		return fmt.Sprintf("%s %s", from, strings.ReplaceAll(operator, "?", to)), nil, nil
	}

	return fmt.Sprintf("%s %s", from, operator), w.To(), nil
}

func NewWhere() specs.DriverWhere {
	return new(where)
}
