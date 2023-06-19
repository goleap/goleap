package specs

type ErrConnectorNotFound interface {
	error
	Name() string
}

type ErrConnectorAlreadyAdded ErrConnectorNotFound
