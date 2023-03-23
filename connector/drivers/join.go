package drivers

import (
	"errors"
	"fmt"
	"github.com/lab210-dev/dbkit/connector/drivers/joins"
	"github.com/lab210-dev/dbkit/specs"
	"strings"
)

type join struct {
	fromKey        string
	fromTableIndex int

	toTable      string
	toKey        string
	toSchema     string
	toTableIndex int

	method specs.JoinMethod
}

func (j *join) Validate() error {
	var check = map[string]func() string{
		"FromKey":  j.FromKey,
		"ToTable":  j.ToTable,
		"ToKey":    j.ToKey,
		"ToSchema": j.ToSchema,
	}

	errList := make([]string, 0)
	for key, c := range check {
		if c() == "" {
			errList = append(errList, key)
		}
	}

	if len(errList) > 0 {
		return errors.New(fmt.Sprintf(`The following fields "%s" are mandatory to perform the join.`, strings.Join(errList, ", ")))
	}

	return nil
}

func (j *join) ToSchema() string {
	return j.toSchema
}

func (j *join) SetToSchema(fromSchema string) specs.DriverJoin {
	j.toSchema = fromSchema
	return j
}

func (j *join) Method() string {
	return joins.Method[j.method]
}

func (j *join) SetMethod(method specs.JoinMethod) specs.DriverJoin {
	j.method = method
	return j
}

func (j *join) FromTableIndex() int {
	return j.fromTableIndex
}

func (j *join) ToTable() string {
	return j.toTable
}

func (j *join) ToTableIndex() int {
	return j.toTableIndex
}

func (j *join) FromKey() string {
	return j.fromKey
}

func (j *join) ToKey() string {
	return j.toKey
}

func (j *join) SetFromTableIndex(fromTableIndex int) specs.DriverJoin {
	j.fromTableIndex = fromTableIndex
	return j
}

func (j *join) SetToTable(toTable string) specs.DriverJoin {
	j.toTable = toTable
	return j
}

func (j *join) SetToTableIndex(toTableIndex int) specs.DriverJoin {
	j.toTableIndex = toTableIndex
	return j
}

func (j *join) SetFromKey(fromKey string) specs.DriverJoin {
	j.fromKey = fromKey
	return j
}

func (j *join) SetToKey(toKey string) specs.DriverJoin {
	j.toKey = toKey
	return j
}

func NewJoin() specs.DriverJoin {
	return new(join)
}
