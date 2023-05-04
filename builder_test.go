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

type BuilderTestSuite struct {
	suite.Suite
	context.Context

	fakeConnector          *mocks.FakeConnector
	fakeConnectors         *mocks.FakeConnectors
	fakeConnectorsInstance *mocks.FakeConnectorsInstance

	fakeModelDefinition *mocks.FakeModelDefinition
	fakeFieldDefinition *mocks.FakeFieldDefinition
	fakeDriverField     *mocks.FakeDriverField
	fakeDriverJoin      *mocks.FakeDriverJoin

	fakeUseModelDefinition      *mocks.FakeUseModelDefinition
	fakeCommentPayloadConstruct *mocks.FakePayloadConstruct[*models.CommentsModel]
	fakeCommentPayloadAugmented *mocks.FakePayloadAugmented[*models.CommentsModel]

	fakePostPayloadConstruct *mocks.FakePayloadConstruct[*models.PostsModel]
	fakePostPayloadAugmented *mocks.FakePayloadAugmented[*models.PostsModel]

	fakeModelPayloadConstruct *mocks.FakePayloadConstruct[specs.Model]
	fakeModelPayloadAugmented *mocks.FakePayloadAugmented[specs.Model]

	fakeNewSubBuilder *mocks.FakeNewSubBuilder[*models.PostsModel]
	fakeSubBuilder    *mocks.FakeSubBuilder[*models.PostsModel]
}

func (test *BuilderTestSuite) SetupTest() {
	test.Context = context.Background()
	test.fakeConnector = mocks.NewFakeConnector(test.T())
	test.fakeConnectors = mocks.NewFakeConnectors(test.T())
	test.fakeConnectorsInstance = mocks.NewFakeConnectorsInstance(test.T())
	test.fakeModelDefinition = mocks.NewFakeModelDefinition(test.T())
	test.fakeFieldDefinition = mocks.NewFakeFieldDefinition(test.T())
	test.fakeUseModelDefinition = mocks.NewFakeUseModelDefinition(test.T())
	test.fakeDriverField = mocks.NewFakeDriverField(test.T())
	test.fakeDriverJoin = mocks.NewFakeDriverJoin(test.T())

	test.fakeCommentPayloadConstruct = mocks.NewFakePayloadConstruct[*models.CommentsModel](test.T())
	test.fakeCommentPayloadAugmented = mocks.NewFakePayloadAugmented[*models.CommentsModel](test.T())

	test.fakePostPayloadConstruct = mocks.NewFakePayloadConstruct[*models.PostsModel](test.T())
	test.fakePostPayloadAugmented = mocks.NewFakePayloadAugmented[*models.PostsModel](test.T())

	test.fakeModelPayloadConstruct = mocks.NewFakePayloadConstruct[specs.Model](test.T())
	test.fakeModelPayloadAugmented = mocks.NewFakePayloadAugmented[specs.Model](test.T())

	test.fakeNewSubBuilder = mocks.NewFakeNewSubBuilder[*models.PostsModel](test.T())
	test.fakeSubBuilder = mocks.NewFakeSubBuilder[*models.PostsModel](test.T())

	depkit.Reset()
	depkit.Register[specs.UseModelDefinition](test.fakeUseModelDefinition.Use)
	depkit.Register[specs.NewPayload[*models.CommentsModel]](test.fakeCommentPayloadConstruct.NewPayload)
	depkit.Register[specs.NewPayload[*models.PostsModel]](test.fakePostPayloadConstruct.NewPayload)
	depkit.Register[specs.NewPayload[specs.Model]](test.fakeModelPayloadConstruct.NewPayload)
	depkit.Register[specs.NewSubBuilder[*models.PostsModel]](test.fakeNewSubBuilder.NewSubBuilder)
	depkit.Register[specs.ConnectorsInstance](test.fakeConnectorsInstance.Instance)
}

func (test *BuilderTestSuite) TestGetWithNoPrimaryKeyErr() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context)
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

	builderInstance := Use[*models.CommentsModel](test.Context)
	if !test.NotEmpty(builderInstance) {
		return
	}

	_, err := builderInstance.SetFields("unknown").Get("Primary")
	test.Error(err)
}

func (test *BuilderTestSuite) TestValideRequiredFieldErr() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context)
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

	builderInstance := Use[*models.CommentsModel](test.Context)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id").Once()

	test.fakeFieldDefinition.On("FromSlice").Return(false).Once()
	test.fakeModelDefinition.On("GetFieldByName", "unknown").Return(test.fakeFieldDefinition, nil).Once()

	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(nil, errors.New("test")).Once()

	_, err := builderInstance.SetFields("unknown").Get("Primary")
	test.Error(err)
}

func (test *BuilderTestSuite) TestGetWithNotFoundErr() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build fields
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build where

	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()
	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()

	test.fakeFieldDefinition.On("Join").Return([]specs.DriverJoin{}).Times(2)

	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id").Times(3)
	test.fakeFieldDefinition.On("FromSlice").Return(false).Once()

	test.fakeModelDefinition.On("TypeName").Return("mock")

	test.fakeCommentPayloadConstruct.On("NewPayload", (*models.CommentsModel)(nil)).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetFields", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetWheres", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetJoins", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("Result").Return([]*models.CommentsModel{})

	test.fakeConnectorsInstance.On("Instance").Return(test.fakeConnectors, nil)
	test.fakeConnectors.On("Get", "acceptance").Return(test.fakeConnector, nil)

	test.fakeConnector.On("Select", test.Context, mock.Anything).Return(nil)

	_, err := builderInstance.SetFields("Id").Get("Primary")
	test.Error(err)

	test.ErrorContains(err, "empty result for mock")
}

func (test *BuilderTestSuite) TestGetSelectErr() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id").Times(3)
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build fields
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build where
	test.fakeFieldDefinition.On("FromSlice").Return(false).Once()

	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()
	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()

	test.fakeFieldDefinition.On("Join").Return([]specs.DriverJoin{}).Times(2)

	test.fakeCommentPayloadConstruct.On("NewPayload", (*models.CommentsModel)(nil)).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetFields", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetWheres", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetJoins", mock.Anything).Return(test.fakeCommentPayloadAugmented)

	test.fakeConnectorsInstance.On("Instance").Return(test.fakeConnectors, nil)
	test.fakeConnectors.On("Get", "acceptance").Return(test.fakeConnector, nil)

	test.fakeConnector.On("Select", test.Context, mock.Anything).Return(errors.New("select_return_err"))

	_, err := builderInstance.SetFields("Id").Get("Primary")
	test.Error(err)
	test.ErrorContains(err, "select_return_err")
}

func (test *BuilderTestSuite) TestBuildPayloadJoinErr() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id").Times(2)
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build fields
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build where
	test.fakeFieldDefinition.On("FromSlice").Return(false).Once()

	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()
	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()

	test.fakeFieldDefinition.On("Join").Return([]specs.DriverJoin{test.fakeDriverJoin}).Once()

	test.fakeCommentPayloadConstruct.On("NewPayload", (*models.CommentsModel)(nil)).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetFields", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetWheres", mock.Anything).Return(test.fakeCommentPayloadAugmented)

	test.fakeDriverJoin.On("Formatted").Return("", errors.New("join_err")).Once()

	_, err := builderInstance.SetFields("Id").Get("Primary")
	test.Error(err)
	test.EqualValues("join_err", err.Error())
}

func (test *BuilderTestSuite) TestGet() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id").Times(3)
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build fields
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build where
	test.fakeFieldDefinition.On("FromSlice").Return(false).Once()
	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()
	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()

	test.fakeFieldDefinition.On("Join").Return([]specs.DriverJoin{test.fakeDriverJoin}).Once()
	test.fakeDriverJoin.On("Formatted").Return("JOIN `comments` ON `comments`.`id` = `posts`.`id`", nil).Once()

	test.fakeCommentPayloadConstruct.On("NewPayload", (*models.CommentsModel)(nil)).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetFields", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetWheres", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetJoins", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("Result").Return([]*models.CommentsModel{{Id: 1}})

	test.fakeConnectorsInstance.On("Instance").Return(test.fakeConnectors, nil)
	test.fakeConnectors.On("Get", "acceptance").Return(test.fakeConnector, nil)

	test.fakeConnector.On("Select", test.Context, mock.Anything).Return(nil)

	comment, err := builderInstance.SetFields("Id").Get("Primary")
	test.NoError(err)

	test.Equal(comment.Id, uint(1))
}

func (test *BuilderTestSuite) TestWheres() {
	test.fakeNewSubBuilder.On("NewSubBuilder").Return(test.fakeSubBuilder).Once()
	test.fakeUseModelDefinition.On("Use", (*models.PostsModel)(nil)).Return(test.fakeModelDefinition).Once()
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition).Once()

	builderInstance := Use[*models.PostsModel](test.Context)
	if !test.NotEmpty(builderInstance) {
		return
	}
	condition := NewCondition()
	builderInstance.SetWhere(condition)
	test.Equal(builderInstance.Wheres(), []specs.Condition{condition})
}

func (test *BuilderTestSuite) TestFields() {
	test.fakeNewSubBuilder.On("NewSubBuilder").Return(test.fakeSubBuilder).Once()
	test.fakeUseModelDefinition.On("Use", (*models.PostsModel)(nil)).Return(test.fakeModelDefinition).Once()
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition).Once()

	builderInstance := Use[*models.PostsModel](test.Context)
	if !test.NotEmpty(builderInstance) {
		return
	}

	builderInstance.SetFields("Hello", "World")
	test.Equal(builderInstance.Fields(), []string{"Hello", "World"})
}

func (test *BuilderTestSuite) TestSubBuilder() {
	test.fakeUseModelDefinition.On("Use", (*models.PostsModel)(nil)).Return(test.fakeModelDefinition).Once()
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition).Once()

	test.fakeNewSubBuilder.On("NewSubBuilder").Return(test.fakeSubBuilder).Once()

	builderInstance := Use[*models.PostsModel](test.Context)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.fakeFieldDefinition.On("Model").Return(test.fakeModelDefinition).Twice()
	test.fakeModelDefinition.On("FromField").Return(test.fakeFieldDefinition).Once()

	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id").Twice()
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Comments.Id").Once()

	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Times(2)        // for build fields
	test.fakeModelDefinition.On("GetFieldByName", "Comments.Id").Return(test.fakeFieldDefinition, nil).Once() // for build fields

	test.fakeFieldDefinition.On("FromSlice").Return(false).Once()
	test.fakeFieldDefinition.On("FromSlice").Return(true).Once()

	test.fakeFieldDefinition.On("FundamentalName").Return("Comments").Once()

	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()
	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()

	test.fakeFieldDefinition.On("Join").Return([]specs.DriverJoin{}).Times(2)

	test.fakePostPayloadConstruct.On("NewPayload", (*models.PostsModel)(nil)).Return(test.fakePostPayloadAugmented).Once()
	test.fakePostPayloadAugmented.On("SetFields", mock.Anything).Return(test.fakePostPayloadAugmented).Once()
	test.fakePostPayloadAugmented.On("SetWheres", mock.Anything).Return(test.fakePostPayloadAugmented).Once()
	test.fakePostPayloadAugmented.On("SetJoins", mock.Anything).Return(test.fakePostPayloadAugmented).Once()

	test.fakeConnectorsInstance.On("Instance").Return(test.fakeConnectors, nil)
	test.fakeConnectors.On("Get", "acceptance").Return(test.fakeConnector, nil)

	test.fakeConnector.On("Select", test.Context, mock.Anything).Return(nil).Once()
	posts := []*models.PostsModel{{Id: 1}}
	test.fakePostPayloadAugmented.On("Result").Return(posts).Once()

	test.fakeSubBuilder.On("AddJob", builderInstance, "Comments", test.fakeModelDefinition).Return(test.fakeSubBuilder).Once()
	test.fakeSubBuilder.On("Execute").Return(nil).Once()

	commentResult, err := builderInstance.SetFields("Id", "Comments.Id").Get("Primary")
	if !test.Empty(err) {
		return
	}

	test.Equal(commentResult.Id, uint(1))
}

func (test *BuilderTestSuite) TestSubBuilderErr() {
	test.fakeUseModelDefinition.On("Use", (*models.PostsModel)(nil)).Return(test.fakeModelDefinition).Once()
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition).Once()

	test.fakeNewSubBuilder.On("NewSubBuilder").Return(test.fakeSubBuilder).Once()

	builderInstance := Use[*models.PostsModel](test.Context)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.fakeFieldDefinition.On("Model").Return(test.fakeModelDefinition).Twice()
	test.fakeModelDefinition.On("FromField").Return(test.fakeFieldDefinition).Once()

	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id").Twice()
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Comments.Id").Once()

	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Times(2)        // for build fields
	test.fakeModelDefinition.On("GetFieldByName", "Comments.Id").Return(test.fakeFieldDefinition, nil).Once() // for build fields

	test.fakeFieldDefinition.On("FromSlice").Return(false).Once()
	test.fakeFieldDefinition.On("FromSlice").Return(true).Once()

	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()
	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()

	test.fakeFieldDefinition.On("FundamentalName").Return("Comments").Once()

	test.fakeFieldDefinition.On("Join").Return([]specs.DriverJoin{}).Times(2)

	test.fakePostPayloadConstruct.On("NewPayload", (*models.PostsModel)(nil)).Return(test.fakePostPayloadAugmented).Once()
	test.fakePostPayloadAugmented.On("SetFields", mock.Anything).Return(test.fakePostPayloadAugmented).Once()
	test.fakePostPayloadAugmented.On("SetWheres", mock.Anything).Return(test.fakePostPayloadAugmented).Once()
	test.fakePostPayloadAugmented.On("SetJoins", mock.Anything).Return(test.fakePostPayloadAugmented).Once()

	test.fakeConnectorsInstance.On("Instance").Return(test.fakeConnectors, nil)
	test.fakeConnectors.On("Get", "acceptance").Return(test.fakeConnector, nil)

	test.fakeConnector.On("Select", test.Context, mock.Anything).Return(nil).Once()

	test.fakeSubBuilder.On("AddJob", builderInstance, "Comments", test.fakeModelDefinition).Return(test.fakeSubBuilder).Once()
	test.fakeSubBuilder.On("Execute").Return(errors.New("sub_builder_err")).Once()

	_, err := builderInstance.SetFields("Id", "Comments.Id").Get("Primary")
	test.Error(err)
	test.EqualValues("sub_builder_err", err.Error())
}

func (test *BuilderTestSuite) TestDelete() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context)
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
	builderInstance := Use[*models.CommentsModel](test.Context)
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
	builderInstance := Use[*models.CommentsModel](test.Context)
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

	builderInstance := Use[*models.CommentsModel](test.Context)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		_ = builderInstance.SetLimit(0)
	})
}

func (test *BuilderTestSuite) TestOffset() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	builderInstance := Use[*models.CommentsModel](test.Context)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		_ = builderInstance.SetOffset(0)
	})
}

func (test *BuilderTestSuite) TestOrderBy() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)
	builderInstance := Use[*models.CommentsModel](test.Context)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.Panics(func() {
		// @TODO Maybe pass SetOrderBy struct instead of string
		_ = builderInstance.SetOrderBy("Id", "ASC")
	})
}

func (test *BuilderTestSuite) TestCount() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)
	builderInstance := Use[*models.CommentsModel](test.Context)
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
