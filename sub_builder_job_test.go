package dbkit

import (
	"context"
	"errors"
	"github.com/kitstack/dbkit/connector/drivers/operators"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/dbkit/tests/mocks"
	"github.com/kitstack/dbkit/tests/models"
	"github.com/kitstack/depkit"
	structKitSpecs "github.com/kitstack/structkit/specs"
	structKitMocks "github.com/kitstack/structkit/tests/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SubBuilderJobTestSuite struct {
	suite.Suite
	context.Context
	fakeModelDefinition    *mocks.FakeModelDefinition
	fakeFieldDefinition    *mocks.FakeFieldDefinition
	fakePayloadAugmented   *mocks.FakePayloadAugmented[specs.Model]
	fakeBuilder            *mocks.FakeBuilder[specs.Model]
	fakeConnector          *mocks.FakeConnector
	fakeUseModelDefinition *mocks.FakeUseModelDefinition
	fakeBuilderUse         *mocks.FakeBuilderUse[specs.Model]

	fakeGet *structKitMocks.FakeGet
	fakeSet *structKitMocks.FakeSet
}

func (test *SubBuilderJobTestSuite) SetupTest() {
	test.Context = context.Background()
	test.fakeModelDefinition = mocks.NewFakeModelDefinition(test.T())
	test.fakeFieldDefinition = mocks.NewFakeFieldDefinition(test.T())
	test.fakePayloadAugmented = mocks.NewFakePayloadAugmented[specs.Model](test.T())
	test.fakeBuilder = mocks.NewFakeBuilder[specs.Model](test.T())
	test.fakeConnector = mocks.NewFakeConnector(test.T())
	test.fakeUseModelDefinition = mocks.NewFakeUseModelDefinition(test.T())
	test.fakeBuilderUse = mocks.NewFakeBuilderUse[specs.Model](test.T())
	test.fakeGet = structKitMocks.NewFakeGet(test.T())
	test.fakeSet = structKitMocks.NewFakeSet(test.T())

	depkit.Reset()
	depkit.Register[specs.UseModelDefinition](test.fakeUseModelDefinition.Use)
	depkit.Register[specs.BuilderUse[specs.Model]](test.fakeBuilderUse.Use)
	depkit.Register[specs.Builder[specs.Model]](test.fakeBuilder)
	depkit.Register[structKitSpecs.Get](test.fakeGet.Execute)
	depkit.Register[structKitSpecs.Set](test.fakeSet.Execute)
}

func (test *SubBuilderJobTestSuite) TestNewSubBuilderJobFromErr() {
	newSubBuilderJob := newSubBuilderJob[specs.Model](test.fakeBuilder, "Fundamental", test.fakeModelDefinition)
	test.NotNil(newSubBuilderJob)

	test.fakeModelDefinition.On("FromField").Return(test.fakeFieldDefinition, nil).Once()
	test.fakeFieldDefinition.On("GetByColumn").Return(nil, errors.New("GetByColumn")).Once()

	err := newSubBuilderJob.Execute()
	test.Error(err)
	test.Equal("GetByColumn", err.Error())
}

func (test *SubBuilderJobTestSuite) TestNewSubBuilderJobToErr() {
	newSubBuilderJob := newSubBuilderJob[specs.Model](test.fakeBuilder, "Fundamental", test.fakeModelDefinition)
	test.NotNil(newSubBuilderJob)

	test.fakeModelDefinition.On("FromField").Return(test.fakeFieldDefinition, nil).Once()
	test.fakeFieldDefinition.On("GetByColumn").Return(test.fakeFieldDefinition, nil).Once()
	test.fakeFieldDefinition.On("GetToColumn").Return(nil, errors.New("GetToColumn")).Once()

	err := newSubBuilderJob.Execute()
	test.Error(err)
	test.Equal("GetToColumn", err.Error())
}

func (test *SubBuilderJobTestSuite) TestNewSubBuilderJob() {
	comments := &models.CommentsModel{}
	post := &models.PostsModel{Id: 1}

	test.fakeModelDefinition.On("FromField").Return(test.fakeFieldDefinition, nil).Once()
	test.fakeFieldDefinition.On("GetByColumn").Return(test.fakeFieldDefinition, nil).Once()
	test.fakeFieldDefinition.On("GetToColumn").Return(test.fakeFieldDefinition, nil).Once()
	test.fakeBuilder.On("Payload").Return(test.fakePayloadAugmented).Once()
	test.fakePayloadAugmented.On("Result").Return([]specs.Model{post, nil}).Once()

	test.fakeFieldDefinition.On("RecursiveFullName").Return("From").Once()
	test.fakeFieldDefinition.On("RecursiveFullName").Return("PostId").Once()

	test.fakeBuilder.On("Context").Return(context.Background()).Once()
	test.fakeBuilder.On("Connector").Return(test.fakeConnector).Once()

	test.fakeFieldDefinition.On("Model").Return(test.fakeModelDefinition, nil).Once()
	test.fakeModelDefinition.On("Copy").Return(comments).Once()

	test.fakeBuilderUse.On("Use", test.Context, test.fakeConnector).Return(test.fakeBuilder).Once()
	test.fakeBuilder.On("SetModel", comments).Return(test.fakeBuilder)
	test.fakeBuilder.On("Fields").Return([]string{"Fundamental.Label", "Other"})
	test.fakeBuilder.On("SetFields", "Label", "PostId").Return(test.fakeBuilder)
	test.fakeBuilder.On("SetWhere", NewCondition().SetFrom("PostId").SetOperator(operators.In).SetTo([]any{1})).Return(test.fakeBuilder)
	test.fakeBuilder.On("Wheres").Return([]specs.Condition{NewCondition().SetFrom("Fundamental.PostId").SetOperator(operators.In).SetTo([]any{1}), NewCondition().SetFrom("Other")}).Once()

	subResult := []specs.Model{&models.CommentsModel{}}
	test.fakeBuilder.On("FindAll").Return(subResult, nil)

	test.fakeGet.On("Execute", post, "From").Return(1).Once()
	test.fakeGet.On("Execute", nil, "From").Return(nil).Once()

	test.fakeGet.On("Execute", subResult[0], "PostId").Return(1).Once()
	test.fakeSet.On("Execute", post, "Fundamental.[*]", subResult[0]).Return(nil).Once()

	newSubBuilderJob := newSubBuilderJob[specs.Model](test.fakeBuilder, "Fundamental", test.fakeModelDefinition)
	test.NotNil(newSubBuilderJob)

	err := newSubBuilderJob.Execute()
	test.NoError(err)
}

func (test *SubBuilderJobTestSuite) TestNewSubBuilderJobSetError() {
	comments := &models.CommentsModel{}
	post := &models.PostsModel{Id: 1}

	test.fakeModelDefinition.On("FromField").Return(test.fakeFieldDefinition, nil).Once()
	test.fakeFieldDefinition.On("GetByColumn").Return(test.fakeFieldDefinition, nil).Once()
	test.fakeFieldDefinition.On("GetToColumn").Return(test.fakeFieldDefinition, nil).Once()
	test.fakeBuilder.On("Payload").Return(test.fakePayloadAugmented).Once()
	test.fakePayloadAugmented.On("Result").Return([]specs.Model{post, nil}).Once()

	test.fakeFieldDefinition.On("RecursiveFullName").Return("From").Once()
	test.fakeFieldDefinition.On("RecursiveFullName").Return("PostId").Once()

	test.fakeBuilder.On("Context").Return(context.Background()).Once()
	test.fakeBuilder.On("Connector").Return(test.fakeConnector).Once()

	test.fakeFieldDefinition.On("Model").Return(test.fakeModelDefinition, nil).Once()
	test.fakeModelDefinition.On("Copy").Return(comments).Once()

	test.fakeBuilderUse.On("Use", test.Context, test.fakeConnector).Return(test.fakeBuilder).Once()
	test.fakeBuilder.On("SetModel", comments).Return(test.fakeBuilder)
	test.fakeBuilder.On("Fields").Return([]string{"Fundamental.Label", "Other"})
	test.fakeBuilder.On("SetFields", "Label", "PostId").Return(test.fakeBuilder)
	test.fakeBuilder.On("SetWhere", NewCondition().SetFrom("PostId").SetOperator(operators.In).SetTo([]any{1})).Return(test.fakeBuilder)
	test.fakeBuilder.On("Wheres").Return([]specs.Condition{NewCondition().SetFrom("Fundamental.PostId").SetOperator(operators.In).SetTo([]any{1}), NewCondition().SetFrom("Other")}).Once()

	subResult := []specs.Model{&models.CommentsModel{}}
	test.fakeBuilder.On("FindAll").Return(subResult, nil)

	test.fakeGet.On("Execute", post, "From").Return(1).Once()
	test.fakeGet.On("Execute", nil, "From").Return(nil).Once()

	test.fakeGet.On("Execute", subResult[0], "PostId").Return(1).Once()
	test.fakeSet.On("Execute", post, "Fundamental.[*]", subResult[0]).Return(errors.New("set")).Once()

	newSubBuilderJob := newSubBuilderJob[specs.Model](test.fakeBuilder, "Fundamental", test.fakeModelDefinition)
	test.NotNil(newSubBuilderJob)

	err := newSubBuilderJob.Execute()
	test.Error(err)
	test.Equal("set", err.Error())
}

func (test *SubBuilderJobTestSuite) TestNewSubBuilderJobErr() {
	comments := &models.CommentsModel{}
	post := &models.PostsModel{Id: 1}

	test.fakeModelDefinition.On("FromField").Return(test.fakeFieldDefinition, nil).Once()
	test.fakeFieldDefinition.On("GetByColumn").Return(test.fakeFieldDefinition, nil).Once()
	test.fakeFieldDefinition.On("GetToColumn").Return(test.fakeFieldDefinition, nil).Once()
	test.fakeBuilder.On("Payload").Return(test.fakePayloadAugmented).Once()
	test.fakePayloadAugmented.On("Result").Return([]specs.Model{post, nil}).Once()

	test.fakeFieldDefinition.On("RecursiveFullName").Return("From").Once()
	test.fakeFieldDefinition.On("RecursiveFullName").Return("To").Once()

	test.fakeBuilder.On("Context").Return(context.Background()).Once()
	test.fakeBuilder.On("Connector").Return(test.fakeConnector).Once()

	test.fakeFieldDefinition.On("Model").Return(test.fakeModelDefinition, nil).Once()
	test.fakeModelDefinition.On("Copy").Return(comments).Once()

	test.fakeBuilderUse.On("Use", test.Context, test.fakeConnector).Return(test.fakeBuilder).Once()
	test.fakeBuilder.On("SetModel", comments).Return(test.fakeBuilder)
	test.fakeBuilder.On("Fields").Return([]string{"Fundamental.Label"})
	test.fakeBuilder.On("SetFields", "Label", "To").Return(test.fakeBuilder)
	test.fakeBuilder.On("SetWhere", NewCondition().SetFrom("To").SetOperator(operators.In).SetTo([]any{1})).Return(test.fakeBuilder)
	test.fakeBuilder.On("Wheres").Return([]specs.Condition{NewCondition().SetFrom("Fundamental.To").SetOperator(operators.In).SetTo([]any{1})}).Once()
	test.fakeBuilder.On("FindAll").Return(nil, errors.New("FindAll"))

	test.fakeGet.On("Execute", post, "From").Return(1).Once()
	test.fakeGet.On("Execute", nil, "From").Return(nil).Once()

	newSubBuilderJob := newSubBuilderJob[specs.Model](test.fakeBuilder, "Fundamental", test.fakeModelDefinition)
	test.NotNil(newSubBuilderJob)

	err := newSubBuilderJob.Execute()
	test.Error(err)
	test.Equal("FindAll", err.Error())
}

func TestSubBuilderJobTestSuite(t *testing.T) {
	suite.Run(t, new(SubBuilderJobTestSuite))
}
