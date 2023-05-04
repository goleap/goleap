package specs

type Model interface {
	DatabaseName() string
	TableName() string
	ConnectorName() string
}
