package specs

type Connectors interface {
	Add(Connector) ErrConnectorAlreadyAdded
	Get(name string) (Connector, ErrConnectorNotFound)
	Remove(name string)
	List() []Connector
	Clear()
}

type ConnectorsInstance func() Connectors
