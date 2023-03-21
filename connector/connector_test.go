package connector

import (
	"context"
	"database/sql"
	"github.com/lab210-dev/dbkit/connector/config"
	"github.com/lab210-dev/dbkit/connector/drivers"
	"github.com/lab210-dev/dbkit/specs"
	"github.com/lab210-dev/dbkit/tests/mocks/fakesql"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConnectorTestSuite struct {
	suite.Suite
	context.Context
	fakeDriver *fakesql.FakeDriver
}

func (test *ConnectorTestSuite) SetupTest() {
	test.fakeDriver = fakesql.NewDriver(test.T())
	test.Context = context.Background()
}

func (test *ConnectorTestSuite) SetupSuite() {
	drivers.RegisteredDriver = map[string]func() specs.Driver{
		"test": func() specs.Driver {
			return new(drivers.Mysql)
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
