package dbkit

import (
	"context"
	"errors"
	"github.com/kitstack/dbkit/definitions"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/dbkit/tests/mocks"
	"github.com/kitstack/dbkit/tests/models"
	"github.com/kitstack/depkit"
	structKitSpecs "github.com/kitstack/structkit/specs"
	structKitMocks "github.com/kitstack/structkit/tests/mocks"
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

	fakeUseModelDefinition      *mocks.FakeUseModelDefinition
	fakeCommentPayloadConstruct *mocks.FakePayloadConstruct[*models.CommentsModel]
	fakeCommentPayloadAugmented *mocks.FakePayloadAugmented[*models.CommentsModel]

	fakePostPayloadConstruct *mocks.FakePayloadConstruct[*models.PostsModel]
	fakePostPayloadAugmented *mocks.FakePayloadAugmented[*models.PostsModel]

	fakeModelPayloadConstruct *mocks.FakePayloadConstruct[specs.Model]
	fakeModelPayloadAugmented *mocks.FakePayloadAugmented[specs.Model]

	fakeStructGet *structKitMocks.FakeGet
	fakeStructSet *structKitMocks.FakeSet
}

func (test *BuilderTestSuite) SetupTest() {
	test.Context = context.Background()
	test.fakeConnector = mocks.NewFakeConnector(test.T())
	test.fakeModelDefinition = mocks.NewFakeModelDefinition(test.T())
	test.fakeFieldDefinition = mocks.NewFakeFieldDefinition(test.T())
	test.fakeUseModelDefinition = mocks.NewFakeUseModelDefinition(test.T())
	test.fakeDriverField = mocks.NewFakeDriverField(test.T())

	test.fakeCommentPayloadConstruct = mocks.NewFakePayloadConstruct[*models.CommentsModel](test.T())
	test.fakeCommentPayloadAugmented = mocks.NewFakePayloadAugmented[*models.CommentsModel](test.T())

	test.fakePostPayloadConstruct = mocks.NewFakePayloadConstruct[*models.PostsModel](test.T())
	test.fakePostPayloadAugmented = mocks.NewFakePayloadAugmented[*models.PostsModel](test.T())

	test.fakeModelPayloadConstruct = mocks.NewFakePayloadConstruct[specs.Model](test.T())
	test.fakeModelPayloadAugmented = mocks.NewFakePayloadAugmented[specs.Model](test.T())

	test.fakeStructGet = structKitMocks.NewFakeGet(test.T())
	test.fakeStructSet = structKitMocks.NewFakeSet(test.T())

	depkit.Reset()
	depkit.Register[specs.UseModelDefinition](test.fakeUseModelDefinition.Use)
	depkit.Register[specs.NewPayload[*models.CommentsModel]](test.fakeCommentPayloadConstruct.NewPayload)
	depkit.Register[specs.NewPayload[*models.PostsModel]](test.fakePostPayloadConstruct.NewPayload)
	depkit.Register[specs.NewPayload[specs.Model]](test.fakeModelPayloadConstruct.NewPayload)
	depkit.Register[structKitSpecs.Get](test.fakeStructGet.Execute)
	depkit.Register[structKitSpecs.Set](test.fakeStructSet.Execute)
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

	test.fakeCommentPayloadConstruct.On("NewPayload", (*models.CommentsModel)(nil)).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetFields", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetWheres", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetJoins", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("Result").Return([]*models.CommentsModel{})

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

	test.fakeCommentPayloadConstruct.On("NewPayload", (*models.CommentsModel)(nil)).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetFields", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetWheres", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetJoins", mock.Anything).Return(test.fakeCommentPayloadAugmented)

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

	test.fakeCommentPayloadConstruct.On("NewPayload", (*models.CommentsModel)(nil)).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetFields", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetWheres", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("SetJoins", mock.Anything).Return(test.fakeCommentPayloadAugmented)
	test.fakeCommentPayloadAugmented.On("Result").Return([]*models.CommentsModel{{Id: 1}})

	test.fakeConnector.On("Select", test.Context, mock.Anything).Return(nil)

	comment, err := builderInstance.Fields("Id").Get("Primary")
	if !test.Empty(err) {
		return
	}

	test.Equal(comment.Id, uint(1))
}

func (test *BuilderTestSuite) TestMany2Many() {
	test.fakeUseModelDefinition.On("Use", (*models.PostsModel)(nil)).Return(test.fakeModelDefinition).Once()
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition).Once()

	builderInstance := Use[*models.PostsModel](test.Context, test.fakeConnector)
	if !test.NotEmpty(builderInstance) {
		return
	}

	test.fakeFieldDefinition.On("Model").Return(test.fakeModelDefinition).Twice()
	test.fakeModelDefinition.On("FromField").Return(test.fakeFieldDefinition).Once()
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Comments.Id").Once()

	test.fakeModelDefinition.On("GetPrimaryField").Return(test.fakeFieldDefinition, nil)
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id").Once()
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build fields
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build where

	test.fakeModelDefinition.On("GetFieldByName", "Comments.Id").Return(test.fakeFieldDefinition, nil).Once() // for build fields
	test.fakeModelDefinition.On("GetFieldByName", "Comments.Id").Return(test.fakeFieldDefinition, nil).Once() // for build where

	test.fakeFieldDefinition.On("FromSlice").Return(false).Once()
	test.fakeFieldDefinition.On("FromSlice").Return(true).Once()

	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()
	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()

	test.fakeFieldDefinition.On("Join").Return([]specs.DriverJoin{}).Once()

	test.fakePostPayloadConstruct.On("NewPayload", (*models.PostsModel)(nil)).Return(test.fakePostPayloadAugmented).Once()
	test.fakePostPayloadAugmented.On("SetFields", mock.Anything).Return(test.fakePostPayloadAugmented).Once()
	test.fakePostPayloadAugmented.On("SetWheres", mock.Anything).Return(test.fakePostPayloadAugmented).Once()
	test.fakePostPayloadAugmented.On("SetJoins", mock.Anything).Return(test.fakePostPayloadAugmented).Once()

	test.fakeConnector.On("Select", test.Context, mock.Anything).Return(nil).Once()
	posts := []*models.PostsModel{{Id: 1}}
	test.fakePostPayloadAugmented.On("Result").Return(posts).Once()

	test.fakeModelDefinition.On("FromField").Return(test.fakeFieldDefinition).Once()
	test.fakeFieldDefinition.On("GetByColumn").Return(test.fakeFieldDefinition, nil).Once()

	test.fakeModelDefinition.On("FromField").Return(test.fakeFieldDefinition).Once()
	test.fakeFieldDefinition.On("GetToColumn").Return(test.fakeFieldDefinition, nil).Once()

	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id").Once()
	test.fakeFieldDefinition.On("Model").Return(test.fakeModelDefinition).Once()

	test.fakeFieldDefinition.On("FromSlice").Return(false).Once()
	test.fakeModelDefinition.On("GetFieldByName", "Id").Return(test.fakeFieldDefinition, nil).Once() // for build fields
	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()
	test.fakeFieldDefinition.On("Field").Return(test.fakeDriverField).Once()
	test.fakeFieldDefinition.On("Join").Return([]specs.DriverJoin{}).Once()

	test.fakeConnector.On("Select", test.Context, mock.Anything).Return(nil).Once()

	// Build One2Many/Many2Many Case
	comment := &models.CommentsModel{}
	test.fakeModelDefinition.On("Copy").Return(comment).Once()
	test.fakeUseModelDefinition.On("Use", comment).Return(test.fakeModelDefinition).Once()
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition).Once()
	test.fakeFieldDefinition.On("RecursiveFullName").Return("Id").Once()
	test.fakeStructGet.On("Execute", posts[0], "Id").Return(1).Once()

	test.fakeModelPayloadConstruct.On("NewPayload", comment).Return(test.fakeModelPayloadAugmented).Once()
	test.fakeModelPayloadAugmented.On("SetFields", mock.Anything).Return(test.fakeModelPayloadAugmented).Once()
	test.fakeModelPayloadAugmented.On("SetWheres", mock.Anything).Return(test.fakeModelPayloadAugmented).Once()
	test.fakeModelPayloadAugmented.On("SetJoins", mock.Anything).Return(test.fakeModelPayloadAugmented).Once()

	comments := []specs.Model{&models.CommentsModel{PostId: 1, Content: "Hello"}}
	test.fakeModelPayloadAugmented.On("Result").Return(comments).Once()

	test.fakeStructGet.On("Execute", comments[0], "Id").Return(1).Once()
	test.fakeStructSet.On("Execute", posts[0], "Id.[*]", comments[0]).Return(nil).Once()

	commentResult, err := builderInstance.Fields("Id", "Comments.Id").Get("Primary")
	if !test.Empty(err) {
		return
	}

	test.Equal(commentResult.Id, uint(1))
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
