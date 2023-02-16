package connector

import (
	"context"
	"github.com/goleap/goleap/connector/config"
	"github.com/goleap/goleap/connector/driver"
)

type connector struct {
	driver.Driver

	name   string
	config config.Config
}

type Connector interface {
	driver.Driver

	GetCnx(ctx context.Context)

	Config() config.Config

	Name() string
	SetName(name string) Connector
}

func New(name string, config config.Config) (Connector, error) {
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

	c.Driver, err = driver.Get(c.Config().Driver())
	if err != nil {
		return err
	}

	return c.New(c.Config())
}

func (c *connector) Config() config.Config {
	return c.config
}

func (c *connector) Name() string {
	return c.name
}

func (c *connector) SetName(name string) Connector {
	c.name = name
	return c
}

func (c *connector) GetCnx(ctx context.Context) {

}
