package connectors

import (
	"context"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/dbkit/tests/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConnectorsTestSuite struct {
	suite.Suite
	context.Context

	fakeConnector *mocks.FakeConnector
}

func (test *ConnectorsTestSuite) SetupTest() {
	test.Context = context.Background()
	test.fakeConnector = mocks.NewFakeConnector(test.T())

	Instance().Clear()

	test.fakeConnector.On("Name").Return("test")
	err := Instance().Add(test.fakeConnector)
	if !test.Empty(err) {
		return
	}
}

func (test *ConnectorsTestSuite) TestAddConnectorTwice() {
	err := Instance().Add(test.fakeConnector)
	test.Error(err)
	test.Implements((*specs.ErrConnectorAlreadyAdded)(nil), err)
	test.Equal("connector `test` already added", err.Error())
	test.Equal("test", err.Name())
}

func (test *ConnectorsTestSuite) TestGetConnector() {
	conn, err := Instance().Get("test")
	if !test.Empty(err) {
		return
	}

	test.Equal("test", conn.Name())
}

func (test *ConnectorsTestSuite) TestGetConnectorNotFound() {
	conn, err := Instance().Get("unknown")
	test.Error(err)
	test.Implements((*specs.ErrConnectorNotFound)(nil), err)
	test.Equal("unknown connector: unknown", err.Error())
	test.Equal("unknown", err.Name())
	test.Empty(conn)
}

func (test *ConnectorsTestSuite) TestListConnectors() {
	connectors := Instance().List()
	test.Len(connectors, 1)
	test.Equal("test", connectors[0].Name())
}

func (test *ConnectorsTestSuite) TestClearConnectors() {
	Instance().Clear()
	test.Len(Instance().List(), 0)
}

func (test *ConnectorsTestSuite) TestRemoveConnector() {
	Instance().Remove("test")
	test.Len(Instance().List(), 0)
}

func TestConnectorsTestSuite(t *testing.T) {
	suite.Run(t, new(ConnectorsTestSuite))
}
