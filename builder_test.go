package dbkit

import (
	"context"
	"errors"
	"github.com/kitstack/dbkit/definitions"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/dbkit/tests/mocks"
	"github.com/kitstack/dbkit/tests/models"
	"github.com/kitstack/depkit"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/suite"
)

// TODO - Try to use more injection dependencies like definition and payload

type BuilderTestSuite struct {
	suite.Suite
	context.Context

	fakeConnector       *mocks.FakeConnector
	fakeModelDefinition *mocks.FakeModelDefinition
	fakeFieldDefinition *mocks.FakeFieldDefinition
	fakeDriverField     *mocks.FakeDriverField

	fakeUseModelDefinition *mocks.FakeUseModelDefinition
	fakePayloadConstruct   *mocks.FakePayloadConstruct[*models.CommentsModel]
	fakePayloadAugmented   *mocks.FakePayloadAugmented[*models.CommentsModel]
}

func (test *BuilderTestSuite) SetupTest() {
	test.Context = context.Background()
	test.fakeConnector = mocks.NewFakeConnector(test.T())
	test.fakeModelDefinition = mocks.NewFakeModelDefinition(test.T())
	test.fakeFieldDefinition = mocks.NewFakeFieldDefinition(test.T())
	test.fakeUseModelDefinition = mocks.NewFakeUseModelDefinition(test.T())
	test.fakeDriverField = mocks.NewFakeDriverField(test.T())
	test.fakePayloadConstruct = mocks.NewFakePayloadConstruct[*models.CommentsModel](test.T())
	test.fakePayloadAugmented = mocks.NewFakePayloadAugmented[*models.CommentsModel](test.T())

	depkit.Reset()
	depkit.Register[specs.UseModelDefinition](test.fakeUseModelDefinition.Use)
	depkit.Register[specs.NewPayload[*models.CommentsModel]](test.fakePayloadConstruct.NewPayload)
}

func (test *BuilderTestSuite) TestGetWithNoPrimaryKeyErr() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.fakeModelDefinition.On("GetPrimaryField").Return(nil, definitions.NewErrNoPrimaryField(nil)).Once()

	_, err := builderInstance.Get("Primary")

	primaryErr := &definitions.ErrPrimaryFieldNotFound{}
	test.True(errors.As(err, &primaryErr))
}

func (test *BuilderTestSuite) TestBuildFieldsErr() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id")
	test.fakeModelDefinition.On("GetFieldByName", "unknown").Return(nil, errors.New("test")).Once()

	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	_, err := builderInstance.Fields("unknown").Get("Primary")
	test.Error(err)
}

func (test *BuilderTestSuite) TestValideRequiredFieldErr() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id")

	_, err := builderInstance.Get("Primary")
	test.Error(err)
	test.ErrorContains(err, "the method `Get` requires the selection of one or more fields")
}

func (test *BuilderTestSuite) TestBuildWhereErr() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id").Once()

	test.fakeFieldDefinition.On("FromSlice").Return(false).Once()
	test.fakeModelDefinition.On("GetFieldByName", "unknown").Return(test.fakeFieldDefinition, nil).Once()

	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(nil, errors.New("test")).Once()

	_, err := builderInstance.Fields("unknown").Get("Primary")
	test.Error(err)
}

func (test *BuilderTestSuite) TestGetWithNotFoundErr() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build fields
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build where

	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()
	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()

	test.fakeFieldDefinition.On("Join").Return([]specs.DriverJoin{}).Once()

	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id").Once()
	test.fakeFieldDefinition.On("FromSlice").Return(false).Once()

	test.fakeModelDefinition.On("TypeName").Return("mock")

	test.fakePayloadConstruct.On("NewPayload", (*models.CommentsModel)(nil)).Return(test.fakePayloadAugmented)
	test.fakePayloadAugmented.On("SetFields", mock.Anything).Return(test.fakePayloadAugmented)
	test.fakePayloadAugmented.On("SetWheres", mock.Anything).Return(test.fakePayloadAugmented)
	test.fakePayloadAugmented.On("SetJoins", mock.Anything).Return(test.fakePayloadAugmented)
	test.fakePayloadAugmented.On("Result").Return([]*models.CommentsModel{})

	test.fakeConnector.On("Select", test.Context, mock.Anything).Return(nil)

	_, err := builderInstance.Fields("Id").Get("Primary")
	test.Error(err)

	test.ErrorContains(err, "empty result for mock")
}

func (test *BuilderTestSuite) TestGetSelectErr() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id").Once()
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build fields
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build where
	test.fakeFieldDefinition.On("FromSlice").Return(false).Once()

	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()
	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()

	test.fakeFieldDefinition.On("Join").Return([]specs.DriverJoin{}).Once()

	test.fakePayloadConstruct.On("NewPayload", (*models.CommentsModel)(nil)).Return(test.fakePayloadAugmented)
	test.fakePayloadAugmented.On("SetFields", mock.Anything).Return(test.fakePayloadAugmented)
	test.fakePayloadAugmented.On("SetWheres", mock.Anything).Return(test.fakePayloadAugmented)
	test.fakePayloadAugmented.On("SetJoins", mock.Anything).Return(test.fakePayloadAugmented)

	test.fakeConnector.On("Select", test.Context, mock.Anything).Return(errors.New("select_return_err"))

	_, err := builderInstance.Fields("Id").Get("Primary")
	test.Error(err)

	test.ErrorContains(err, "select_return_err")
}

func (test *BuilderTestSuite) TestGet() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id").Once()
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build fields
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build where
	test.fakeFieldDefinition.On("FromSlice").Return(false).Once()
	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()
	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()

	test.fakeFieldDefinition.On("Join").Return([]specs.DriverJoin{}).Once()

	test.fakePayloadConstruct.On("NewPayload", (*models.CommentsModel)(nil)).Return(test.fakePayloadAugmented)
	test.fakePayloadAugmented.On("SetFields", mock.Anything).Return(test.fakePayloadAugmented)
	test.fakePayloadAugmented.On("SetWheres", mock.Anything).Return(test.fakePayloadAugmented)
	test.fakePayloadAugmented.On("SetJoins", mock.Anything).Return(test.fakePayloadAugmented)
	test.fakePayloadAugmented.On("Result").Return([]*models.CommentsModel{{Id: 1}})

	test.fakeConnector.On("Select", test.Context, mock.Anything).Return(nil)

	comment, err := builderInstance.Fields("Id").Get("Primary")
	if !test.Empty(err) {
		return
	}

	test.Equal(comment.Id, uint(1))
}

func (test *BuilderTestSuite) TestDelete() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		_ = builderInstance.Delete("Primary")
	})
}

func (test *BuilderTestSuite) TestCreate() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		_ = builderInstance.Create()
	})
}

func (test *BuilderTestSuite) TestUpdate() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		_ = builderInstance.Update()
	})
}

func (test *BuilderTestSuite) TestLimit() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		_ = builderInstance.Limit(0)
	})
}

func (test *BuilderTestSuite) TestOffset() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		_ = builderInstance.Offset(0)
	})
}

func (test *BuilderTestSuite) TestOrderBy() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)
	builderInstance := Use[*models.CommentsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		// @TODO Maybe pass OrderBy struct instead of string
		_ = builderInstance.OrderBy("Id", "ASC")
	})
}

func (test *BuilderTestSuite) TestCount() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)
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
