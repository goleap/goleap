package drivers

import "github.com/lab210-dev/dbkit/specs"

type join struct {
	fromTable      string
	fromTableIndex int
	toTable        string
	toTableIndex   int
	fromKey        string
	toKey          string
	method         specs.JoinMethod
}

func (j *join) Method() string {
	return specs.Method[j.method]
}

func (j *join) SetMethod(method specs.JoinMethod) specs.DriverJoin {
	j.method = method
	return j
}

func (j *join) FromTable() string {
	return j.fromTable
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

func (j *join) SetFromTable(fromTable string) specs.DriverJoin {
	j.fromTable = fromTable
	return j
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
