package specs

type Connector interface {
	Driver

	//GetCnx(ctx context.Context)

	Config() Config

	Name() string
	SetName(name string) Connector
}
