package dbkit

import (
	"context"
	"errors"
	"github.com/lab210-dev/dbkit/connector/drivers"
	"github.com/lab210-dev/dbkit/definitions"
	"github.com/lab210-dev/dbkit/specs"
	"github.com/lab210-dev/dbkit/tests/mocks"
	"github.com/lab210-dev/dbkit/tests/models"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BuilderTestSuite struct {
	suite.Suite
	context.Context
	fakeConnector       *mocks.FakeConnector
	fakeModelDefinition *mocks.FakeModelDefinition
	fakeFieldDefinition *mocks.FakeFieldDefinition
}

func (test *BuilderTestSuite) SetupTest() {
	test.Context = context.Background()
	test.fakeConnector = mocks.NewFakeConnector(test.T())
	test.fakeModelDefinition = mocks.NewFakeModelDefinition(test.T())
	test.fakeFieldDefinition = mocks.NewFakeFieldDefinition(test.T())
}

func (test *BuilderTestSuite) TestGetWithNoPrimaryKeyErr() {
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	// Mock the model definition to return an error
	b := builderInstance.(*builder[*models.CommentsModel])
	b.modelDefinition = test.fakeModelDefinition

	test.fakeModelDefinition.On("GetPrimaryField").Return(nil, definitions.NewErrNoPrimaryField(nil)).Once()

	_, err := builderInstance.Get("Primary")

	primaryErr := &definitions.ErrNoPrimaryField{}
	test.True(errors.As(err, &primaryErr))
}

func (test *BuilderTestSuite) TestBuildFieldsErr() {
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	// Mock the model definition to return an error
	b := builderInstance.(*builder[*models.CommentsModel])
	b.modelDefinition = test.fakeModelDefinition

	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id")
	test.fakeModelDefinition.On("GetFieldByName", "unknown").Return(nil, errors.New("test")).Once()

	_, err := builderInstance.Fields("unknown").Get("Primary")
	test.Error(err)
}

func (test *BuilderTestSuite) TestValideRequiredFieldErr() {
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	// Mock the model definition to return an error
	b := builderInstance.(*builder[*models.CommentsModel])
	b.modelDefinition = test.fakeModelDefinition

	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id")

	_, err := builderInstance.Get("Primary")
	test.Error(err)
	test.ErrorContains(err, "the method `Get` requires the selection of one or more fields")
}

func (test *BuilderTestSuite) TestBuildWhereErr() {
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	// Mock the model definition to return an error
	b := builderInstance.(*builder[*models.CommentsModel])
	b.modelDefinition = test.fakeModelDefinition

	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id").Once()

	test.fakeFieldDefinition.On("FromSlice").Return(false).Once()
	test.fakeModelDefinition.On("GetFieldByName", "unknown").Return(test.fakeFieldDefinition, nil).Once()

	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(nil, errors.New("test")).Once()

	_, err := builderInstance.Fields("unknown").Get("Primary")
	test.Error(err)
}

func (test *BuilderTestSuite) TestGetWithNotFoundErr() {
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.fakeConnector.On("Select", test.Context, mock.Anything).Return(nil)

	_, err := builderInstance.Fields("Id", "Post.Id").Get("Primary")
	test.Error(err)

	test.ErrorContains(err, "empty result for CommentsModel")
}

func (test *BuilderTestSuite) TestGetSelectErr() {
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.fakeConnector.On("Select", test.Context, mock.Anything).Return(errors.New("select_return_err"))

	_, err := builderInstance.Fields("Id", "Post.Id").Get("Primary")
	test.Error(err)

	test.ErrorContains(err, "select_return_err")
}

func (test *BuilderTestSuite) TestGet() {
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	Id := uint(1)
	ExtraId := uint(2)
	test.fakeConnector.On("Select", test.Context, mock.Anything).Run(func(args mock.Arguments) {
		err := builderInstance.Payload().OnScan([]any{&Id, &ExtraId})
		if !test.Empty(err) {
			return
		}
	}).Return(nil)

	comment, err := builderInstance.Fields("Id", "Post.Id").Get("Primary")
	if !test.Empty(err) {
		return
	}

	test.Equal(0, builderInstance.Payload().Index())
	test.Equal("acceptance", builderInstance.Payload().Database())
	test.Equal("comments", builderInstance.Payload().Table())

	test.NotEmpty(builderInstance.Payload().Where())
	for _, condition := range builderInstance.Payload().Where() {
		test.Equal(drivers.NewField().SetIndex(0).SetColumn("id").SetName("Id"), condition.From())
		test.Equal("=", condition.Operator())
		test.Equal("Primary", condition.To())
	}

	test.Equal([]specs.DriverField{
		drivers.NewField().SetColumn("id").SetIndex(0).SetName("Id"),
		drivers.NewField().SetColumn("id").SetIndex(2).SetName("Post.Id"),
	}, builderInstance.Payload().Fields())

	test.Equal(uint(1), comment.Id)
}

func (test *BuilderTestSuite) TestDelete() {
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		_ = builderInstance.Delete("Primary")
	})
}

func (test *BuilderTestSuite) TestCreate() {
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		_ = builderInstance.Create()
	})
}

func (test *BuilderTestSuite) TestUpdate() {
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		_ = builderInstance.Update()
	})
}

func (test *BuilderTestSuite) TestFindAll() {
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		_ = builderInstance.FindAll()
	})
}

func (test *BuilderTestSuite) TestLimit() {
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		_ = builderInstance.Limit(0)
	})
}

func (test *BuilderTestSuite) TestOffset() {
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		_ = builderInstance.Offset(0)
	})
}

func (test *BuilderTestSuite) TestOrderBy() {
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		_ = builderInstance.OrderBy("Id", "ASC")
	})
}

func (test *BuilderTestSuite) TestCount() {
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		_, _ = builderInstance.Count()
	})
}

func TestBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}
