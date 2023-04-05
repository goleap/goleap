package drivers

import (
	"fmt"
	"github.com/lab210-dev/dbkit/specs"
	"strings"
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

type unknownFieldsErr struct {
	fields []string
}

func (e *unknownFieldsErr) Fields() []string {
	return e.fields
}

func (e *unknownFieldsErr) Error() string {
	return fmt.Sprintf("unknown fields: %s", strings.Join(e.Fields(), ", "))
}

func NewUnknownFieldsErr(fields []string) specs.UnknownFieldsErr {
	return &unknownFieldsErr{fields: fields}
}
