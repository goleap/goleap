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

func NewUnknownOperatorErr(operator string) specs.ErrUnknownOperator {
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

func NewUnknownFieldsErr(fields []string) specs.ErrUnknownFields {
	return &unknownFieldsErr{fields: fields}
}

type requiredFieldJoinErr struct {
	fields []string
}

func (e *requiredFieldJoinErr) Fields() []string {
	return e.fields
}

func (e *requiredFieldJoinErr) Error() string {
	return fmt.Sprintf("the following fields \"%s\" are mandatory to perform the join", strings.Join(e.fields, ", "))
}

func NewRequiredFieldJoinErr(fields []string) specs.ErrRequiredFieldJoin {
	return &requiredFieldJoinErr{fields: fields}
}
