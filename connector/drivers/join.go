package drivers

import (
	"fmt"
	"github.com/lab210-dev/dbkit/connector/drivers/joins"
	"github.com/lab210-dev/dbkit/specs"
	"strings"
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

func (j *join) Formatted() (string, error) {
	return fmt.Sprintf("%s `%s`.`%s` AS `t%d` ON `t%d`.`%s` = `t%d`.`%s`", j.Method(), j.To().Database(), j.To().Table(), j.To().Index(), j.To().Index(), j.To().Column(), j.From().Index(), j.From().Column()), nil
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
		return fmt.Errorf(`the following fields "%s" are mandatory to perform the join`, strings.Join(errList, ", "))
	}

	return nil
}

func NewJoin() specs.DriverJoin {
	return new(join)
}
