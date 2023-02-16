package connector

import (
	"database/sql"
	"github.com/goleap/goleap/connector/config"
	"github.com/goleap/goleap/connector/driver"
	"github.com/goleap/goleap/helper/fakesql"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConnectorTestSuite struct {
	suite.Suite
	fakeDriver *fakesql.FakeDriver
}

func (test *ConnectorTestSuite) SetupTest() {
	test.fakeDriver = fakesql.NewDriver(test.T())
}

func (test *ConnectorTestSuite) SetupSuite() {
	driver.RegisteredDriver = map[string]func() driver.Driver{
		"test": func() driver.Driver {
			return new(driver.Mysql)
		},
	}
	sql.Register("test", test.fakeDriver)
}

func (test *ConnectorTestSuite) TestNewConnector() {
	conn, err := New("test", config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.Equal("test", conn.Name())
	conn.SetName("conn_test")
	test.Equal("conn_test", conn.Name())
}

func (test *ConnectorTestSuite) TestFailNewConnector() {
	conn, err := New("test", config.New().SetDriver("unknown"))
	test.EqualError(err, "driver not found")
	test.Empty(conn)
}

func TestSchemaTestSuite(t *testing.T) {
	suite.Run(t, new(ConnectorTestSuite))
}
