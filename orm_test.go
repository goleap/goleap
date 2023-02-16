package goleap

import (
	"context"
	"github.com/goleap/goleap/connector/driver"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/goleap/goleap/helper"
	"github.com/stretchr/testify/suite"
)

type OrmTestSuite struct {
	suite.Suite
	fakeConnector *helper.FakeConnector
}

func (test *OrmTestSuite) SetupTest() {
	test.fakeConnector = helper.NewFakeConnector(test.T())
}

func (test *OrmTestSuite) TestGet() {
	ctx := context.Background()
	ormInstance := Use[helper.BaseModel](context.Background(), test.fakeConnector)
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

	test.Equal([]driver.Field{
		driver.NewField().SetName("id").SetIndex(0),
		driver.NewField().SetName("id").SetIndex(351),
	}, ormInstance.Payload().Fields())

	test.Equal(uint(1), baseModel.Id)
	test.Equal(uint(2), baseModel.Extra.Id)
}

func TestSchemaTestSuite(t *testing.T) {
	suite.Run(t, new(OrmTestSuite))
}
