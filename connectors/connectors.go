package connectors

import (
	"github.com/kitstack/dbkit/specs"
	log "github.com/sirupsen/logrus"
	"sync"
)

type connectors struct {
	sync.Mutex
	connector map[string]specs.Connector
}

var instance *connectors

func init() {
	instance = &connectors{
		connector: make(map[string]specs.Connector),
	}

	log.WithFields(log.Fields{
		"package": "connectors",
	}).Info("package initialized")
}

func Instance() specs.Connectors {
	return instance
}

func (i *connectors) Add(connector specs.Connector) specs.ErrConnectorAlreadyAdded {
	i.Lock()
	defer i.Unlock()

	if _, ok := i.connector[connector.Name()]; ok {
		return NewConnectorAlreadyAddedError(connector.Name())
	}

	i.connector[connector.Name()] = connector

	log.WithFields(log.Fields{
		"connectorName": connector.Name(),
	}).Debug("Connector added")

	return nil
}

func (i *connectors) Get(name string) (specs.Connector, specs.ErrConnectorNotFound) {
	i.Lock()
	defer i.Unlock()

	if connector, ok := i.connector[name]; ok {
		return connector, nil
	}

	return nil, NewConnectorNotFoundError(name)
}

func (i *connectors) List() (connectors []specs.Connector) {
	i.Lock()
	defer i.Unlock()

	for _, connector := range i.connector {
		connectors = append(connectors, connector)
	}

	return
}

func (i *connectors) Remove(name string) {
	i.Lock()
	defer i.Unlock()

	delete(i.connector, name)

	log.WithFields(log.Fields{
		"name": name,
	}).Info("connector removed")
}

func (i *connectors) Clear() {
	i.Lock()
	defer i.Unlock()

	i.connector = make(map[string]specs.Connector)

	log.Debug("connectors cleared")
}
