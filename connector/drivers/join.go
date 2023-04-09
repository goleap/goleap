package drivers

import (
	"fmt"
	"github.com/lab210-dev/dbkit/connector/drivers/joins"
	"github.com/lab210-dev/dbkit/specs"
)

type join struct {
	from specs.DriverField
	to   specs.DriverField

	method specs.JoinMethod
}

func (j *join) Method() string {
	return joins.Method[j.method]
}

func (j *join) From() specs.DriverField {
	return j.from
}

func (j *join) To() specs.DriverField {
	return j.to
}

func (j *join) SetMethod(method specs.JoinMethod) specs.DriverJoin {
	j.method = method
	return j
}

func (j *join) SetFrom(field specs.DriverField) specs.DriverJoin {
	j.from = field
	return j
}

func (j *join) SetTo(field specs.DriverField) specs.DriverJoin {
	j.to = field
	return j
}

func (j *join) toFormatted() (string, error) {
	formatted, err := j.To().Formatted()
	if err != nil {
		return "", err
	}

	if j.From().IsCustom() {
		return formatted, nil
	}

	return fmt.Sprintf("`%s`.`%s` AS `t%d` ON %s", j.To().Database(), j.To().Table(), j.To().Index(), formatted), nil
}

func (j *join) fromFormatted() (string, error) {
	return j.From().Formatted()
}

func (j *join) Formatted() (string, error) {
	fromFormatted, err := j.fromFormatted()
	if err != nil {
		return "", err
	}

	toFormatted, err := j.toFormatted()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s = %s", j.Method(), toFormatted, fromFormatted), nil
}

func (j *join) Validate() error {
	var check = map[string]func() specs.DriverField{
		"From": j.From,
		"To":   j.To,
	}

	errList := make([]string, 0)
	for key, c := range check {
		if c() == nil {
			errList = append(errList, key)
		}
	}

	if len(errList) > 0 {
		return NewRequiredFieldJoinErr(errList)
	}

	return nil
}

func NewJoin() specs.DriverJoin {
	return new(join)
}
