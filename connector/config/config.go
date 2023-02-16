package config

import "net/url"

type Config interface {
	Name() string
	SetName(name string)

	Driver() string
	SetDriver(driver string) Config

	User() string
	SetUser(user string)

	Password() string
	SetPassword(password string)

	Host() string
	SetHost(host string)

	Port() int
	SetPort(port int)

	Database() string
	SetDatabase(database string)

	Locale() string
	SetLocale(locale string)
}

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

func (c *config) SetName(name string) {
	c.name = name
}

func (c *config) Driver() string {
	return c.driver
}

func (c *config) SetDriver(driver string) Config {
	c.driver = driver
	return c
}

func (c *config) User() string {
	return c.user
}

func (c *config) SetUser(user string) {
	c.user = user
}

func (c *config) Password() string {
	return c.password
}

func (c *config) SetPassword(password string) {
	c.password = password
}

func (c *config) Host() string {
	return c.host
}

func (c *config) SetHost(host string) {
	c.host = host
}

func (c *config) Port() int {
	return c.port
}

func (c *config) SetPort(port int) {
	c.port = port
}

func (c *config) Database() string {
	return c.database
}

func (c *config) SetDatabase(database string) {
	c.database = database
}

func (c *config) Locale() string {
	if c.locale == "" {
		return "Local"
	}
	return url.QueryEscape(c.locale)
}

func (c *config) SetLocale(locale string) {
	c.locale = locale
}

func New() Config {
	return new(config)
}
