package dbkit

import (
	"context"
	"github.com/lab210-dev/dbkit/connector/drivers"
	"github.com/lab210-dev/dbkit/specs"
	"github.com/lab210-dev/dbkit/tests/mocks"
	"github.com/lab210-dev/dbkit/tests/models"
	"github.com/stretchr/testify/mock"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BuilderTestSuite struct {
	suite.Suite
	fakeConnector *mocks.FakeConnector
}

func (test *BuilderTestSuite) SetupTest() {
	test.fakeConnector = mocks.NewFakeConnector(test.T())
}

func (test *BuilderTestSuite) TestGet() {
	ctx := context.Background()
	builderInstance := Use[*models.BaseModel](context.Background(), test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	Id := uint(1)
	ExtraId := uint(2)
	test.fakeConnector.On("Select", ctx, mock.Anything).Run(func(args mock.Arguments) {
		err := builderInstance.Payload().OnScan([]any{&Id, &ExtraId})
		if !test.Empty(err) {
			return
		}
	}).Return(nil)

	baseModel, err := builderInstance.Fields("Id", "Extra.Id").Get("Primary")
	if !test.Empty(err) {
		return
	}

	test.Equal(0, builderInstance.Payload().Index())
	test.Equal("test", builderInstance.Payload().Database())
	test.Equal("base", builderInstance.Payload().Table())

	test.Equal([]specs.DriverField{
		drivers.NewField().SetName("id").SetIndex(0).SetNameInSchema("Id"),
		drivers.NewField().SetName("id").SetIndex(72).SetNameInSchema("Extra.Id"),
	}, builderInstance.Payload().Fields())

	log.Print(baseModel)
	test.Equal(uint(1), baseModel.Id)
	test.Equal(uint(2), baseModel.Extra.Id)
}

func TestBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}
