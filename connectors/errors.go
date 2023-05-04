package connectors

import "fmt"

type ErrConnectorAlreadyAdded ErrConnectorNotFound

type ErrConnectorNotFound struct {
	name string
}

func (e *ErrConnectorNotFound) Error() string {
	return fmt.Sprintf("unknown connector: %s", e.name)
}

func (e *ErrConnectorNotFound) Name() string {
	return e.name
}

func NewConnectorNotFoundError(name string) *ErrConnectorNotFound {
	return &ErrConnectorNotFound{
		name: name,
	}
}

func NewConnectorAlreadyAddedError(name string) *ErrConnectorAlreadyAdded {
	return &ErrConnectorAlreadyAdded{
		name: name,
	}
}

func (e *ErrConnectorAlreadyAdded) Name() string {
	return e.name
}

func (e *ErrConnectorAlreadyAdded) Error() string {
	return fmt.Sprintf("connector `%s` already added", e.name)
}
