package joins

const (
	Default = iota
	Inner
	Left
	Right
)

var Method = [...]string{
	"JOIN",
	"INNER JOIN",
	"LEFT JOIN",
	"RIGHT JOIN",
}
