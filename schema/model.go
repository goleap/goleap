package schema

type Model interface {
	DatabaseName() string
	TableName() string
}
