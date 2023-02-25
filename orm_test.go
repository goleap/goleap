package goleap

import (
	"context"
	"github.com/lab210-dev/dbkit/connector/drivers"
	"github.com/lab210-dev/dbkit/mocks"
	"github.com/lab210-dev/dbkit/specs"
	"github.com/lab210-dev/dbkit/testmodels"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/suite"
)

type OrmTestSuite struct {
	suite.Suite
	fakeConnector *mocks.FakeConnector
}

func (test *OrmTestSuite) SetupTest() {
	test.fakeConnector = mocks.NewFakeConnector(test.T())
}

func (test *OrmTestSuite) TestGet() {
	ctx := context.Background()
	ormInstance := Use[*testmodels.BaseModel](context.Background(), test.fakeConnector)
	if !test.NotEmpty(ormInstance) {
		return
	}

	test.fakeConnector.On("Select", ctx, ormInstance.Payload()).Run(func(args mock.Arguments) {
		err := ormInstance.Payload().OnScan([]any{uint(1), uint(2)})
		if !test.Empty(err) {
			return
		}
	}).Return(nil)

	baseModel, err := ormInstance.Fields("Id", "Extra.Id").Get("Primary")
	if !test.Empty(err) {
		return
	}

	test.Equal(0, ormInstance.Payload().Index())
	test.Equal("test", ormInstance.Payload().Database())
	test.Equal("base", ormInstance.Payload().Table())

	test.Equal([]specs.DriverField{
		drivers.NewField().SetName("id").SetIndex(0),
		drivers.NewField().SetName("id").SetIndex(72),
	}, ormInstance.Payload().Fields())

	test.Equal(uint(1), baseModel.Id)
	test.Equal(uint(2), baseModel.Extra.Id)
}

func TestSchemaTestSuite(t *testing.T) {
	suite.Run(t, new(OrmTestSuite))
}
