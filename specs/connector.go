package specs

type Connector interface {
	Driver

	Config() Config

	Name() string
	SetName(name string) Connector
}
