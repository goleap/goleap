package drivers

import (
	"fmt"
	"github.com/lab210-dev/dbkit/specs"
)

type unknownOperatorErr struct {
	operator string
}

func (e *unknownOperatorErr) Operator() string {
	return e.operator
}

func (e *unknownOperatorErr) Error() string {
	return fmt.Sprintf("unknown operator: %s", e.Operator())
}

func NewUnknownOperatorErr(operator string) specs.UnknownOperatorErr {
	return &unknownOperatorErr{operator: operator}
}
