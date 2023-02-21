package specs

type JoinMethod int

const (
	BasicJoin = iota
	InnerJoin
	LeftJoin
	RightJoin
)

var Method = [...]string{
	"JOIN",
	"INNER JOIN",
	"LEFT JOIN",
	"RIGHT JOIN",
}

type DriverJoin interface {
	FromTable() string
	FromTableIndex() int
	ToTable() string
	ToTableIndex() int
	FromKey() string
	ToKey() string
	Method() string

	SetMethod(method JoinMethod) DriverJoin
	SetFromTable(fromTable string) DriverJoin
	SetFromTableIndex(fromTableIndex int) DriverJoin
	SetToTable(toTable string) DriverJoin
	SetToTableIndex(toTableIndex int) DriverJoin
	SetFromKey(fromKey string) DriverJoin
	SetToKey(toKey string) DriverJoin
}
