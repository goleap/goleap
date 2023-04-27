package connector

import (
	"github.com/kitstack/dbkit/connector/drivers"
	"github.com/kitstack/dbkit/specs"
)

type connector struct {
	specs.Driver

	name   string
	config specs.Config
}

func New(name string, config specs.Config) (specs.Connector, error) {
	connector := new(connector)
	connector.name = name
	connector.config = config

	err := connector.create()
	if err != nil {
		return nil, err
	}

	return connector, nil
}

func (c *connector) create() (err error) {

	c.Driver, err = drivers.Get(c.Config().Driver())
	if err != nil {
		return err
	}

	return c.New(c.Config())
}

func (c *connector) Config() specs.Config {
	return c.config
}

func (c *connector) Name() string {
	return c.name
}

func (c *connector) SetName(name string) specs.Connector {
	c.name = name
	return c
}
