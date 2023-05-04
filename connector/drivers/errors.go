package drivers

import (
	"fmt"
	"github.com/kitstack/dbkit/specs"
	"strings"
)

type ErrUnknownOperator struct {
	operator string
}

func (e *ErrUnknownOperator) Operator() string {
	return e.operator
}

func (e *ErrUnknownOperator) Error() string {
	return fmt.Sprintf("unknown operator: %s", e.Operator())
}

func NewErrUnknownOperator(operator string) specs.ErrUnknownOperator {
	return &ErrUnknownOperator{operator: operator}
}

type ErrUnknownFields struct {
	fields []string
}

func (e *ErrUnknownFields) Fields() []string {
	return e.fields
}

func (e *ErrUnknownFields) Error() string {
	return fmt.Sprintf("unknown fields: %s", strings.Join(e.Fields(), ", "))
}

func NewErrUnknownFields(fields []string) specs.ErrUnknownFields {
	return &ErrUnknownFields{fields: fields}
}

type ErrRequiredFieldJoin struct {
	fields []string
}

func (e *ErrRequiredFieldJoin) Fields() []string {
	return e.fields
}

func (e *ErrRequiredFieldJoin) Error() string {
	return fmt.Sprintf("the following fields \"%s\" are mandatory to perform the join", strings.Join(e.fields, ", "))
}

func NewErrRequiredFieldJoin(fields []string) specs.ErrRequiredFieldJoin {
	return &ErrRequiredFieldJoin{fields: fields}
}
