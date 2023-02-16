package driver

type JoinMethod int

const (
	BasicJoin = iota
	InnerJoin
	LeftJoin
	RightJoin
)

var method = [...]string{
	"JOIN",
	"INNER JOIN",
	"LEFT JOIN",
	"RIGHT JOIN",
}

type Join interface {
	FromTable() string
	FromTableIndex() int
	ToTable() string
	ToTableIndex() int
	FromKey() string
	ToKey() string
	Method() string

	SetMethod(method JoinMethod) Join
	SetFromTable(fromTable string) Join
	SetFromTableIndex(fromTableIndex int) Join
	SetToTable(toTable string) Join
	SetToTableIndex(toTableIndex int) Join
	SetFromKey(fromKey string) Join
	SetToKey(toKey string) Join
}

type join struct {
	fromTable      string
	fromTableIndex int
	toTable        string
	toTableIndex   int
	fromKey        string
	toKey          string
	method         JoinMethod
}

func (j *join) Method() string {
	return method[j.method]
}

func (j *join) SetMethod(method JoinMethod) Join {
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

func (j *join) SetFromTable(fromTable string) Join {
	j.fromTable = fromTable
	return j
}

func (j *join) SetFromTableIndex(fromTableIndex int) Join {
	j.fromTableIndex = fromTableIndex
	return j
}

func (j *join) SetToTable(toTable string) Join {
	j.toTable = toTable
	return j
}

func (j *join) SetToTableIndex(toTableIndex int) Join {
	j.toTableIndex = toTableIndex
	return j
}

func (j *join) SetFromKey(fromKey string) Join {
	j.fromKey = fromKey
	return j
}

func (j *join) SetToKey(toKey string) Join {
	j.toKey = toKey
	return j
}

func NewJoin() Join {
	return new(join)
}
