package config

import (
	"github.com/kitstack/dbkit/specs"
	"net/url"
)

type config struct {
	name     string
	driver   string
	user     string
	password string
	host     string
	port     int
	database string
	locale   string
}

func (c *config) Name() string {
	return c.name
}

func (c *config) SetName(name string) specs.Config {
	c.name = name
	return c
}

func (c *config) Driver() string {
	return c.driver
}

func (c *config) SetDriver(driver string) specs.Config {
	c.driver = driver
	return c
}

func (c *config) User() string {
	return c.user
}

func (c *config) SetUser(user string) specs.Config {
	c.user = user
	return c
}

func (c *config) Password() string {
	return c.password
}

func (c *config) SetPassword(password string) specs.Config {
	c.password = password
	return c
}

func (c *config) Host() string {
	return c.host
}

func (c *config) SetHost(host string) specs.Config {
	c.host = host
	return c
}

func (c *config) Port() int {
	if c.port == 0 {
		return 3306
	}
	return c.port
}

func (c *config) SetPort(port int) specs.Config {
	c.port = port
	return c
}

func (c *config) Database() string {
	return c.database
}

func (c *config) SetDatabase(database string) specs.Config {
	c.database = database
	return c
}

func (c *config) Locale() string {
	if c.locale == "" {
		return "Local"
	}
	return url.QueryEscape(c.locale)
}

func (c *config) SetLocale(locale string) specs.Config {
	c.locale = locale
	return c
}

func New() specs.Config {
	return new(config)
}
