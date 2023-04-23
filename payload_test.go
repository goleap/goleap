package dbkit

import (
	"context"
	"errors"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/dbkit/tests/mocks"
	"github.com/kitstack/dbkit/tests/models"
	"github.com/kitstack/depkit"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PayloadTestSuite struct {
	suite.Suite
	context.Context
	fakeModelDefinition    *mocks.FakeModelDefinition
	fakeFieldDefinition    *mocks.FakeFieldDefinition
	fakeDriverField        *mocks.FakeDriverField
	fakeDriverLimit        *mocks.FakeDriverLimit
	fakeDriverWhere        *mocks.FakeDriverWhere
	fakeUseModelDefinition *mocks.FakeUseModelDefinition
}

func (test *PayloadTestSuite) SetupTest() {
	test.Context = context.Background()
	test.fakeModelDefinition = mocks.NewFakeModelDefinition(test.T())
	test.fakeFieldDefinition = mocks.NewFakeFieldDefinition(test.T())
	test.fakeDriverField = mocks.NewFakeDriverField(test.T())
	test.fakeDriverLimit = mocks.NewFakeDriverLimit(test.T())
	test.fakeDriverWhere = mocks.NewFakeDriverWhere(test.T())
	test.fakeUseModelDefinition = mocks.NewFakeUseModelDefinition(test.T())

	depkit.Reset()
	depkit.Register[specs.UseModelDefinition](test.fakeUseModelDefinition.Use)
}

func (test *PayloadTestSuite) TestMappingWithGetFieldByNameErr() {
	newPayload := NewPayload[specs.Model]()
	tmp := newPayload.(*payload[specs.Model])

	tmp.modelDefinition = test.fakeModelDefinition

	newPayload.SetFields([]specs.DriverField{
		test.fakeDriverField,
	})

	test.fakeDriverField.On("Name").Return("Test").Once()
	test.fakeModelDefinition.On("GetFieldByName", "Test").Return(nil, errors.New("test")).Once()

	_, err := newPayload.Mapping()
	test.Error(err)
}

func (test *PayloadTestSuite) TestMappingSuccessful() {
	newPayload := NewPayload[specs.Model]()
	tmp := newPayload.(*payload[specs.Model])

	tmp.modelDefinition = test.fakeModelDefinition

	newPayload.SetFields([]specs.DriverField{
		test.fakeDriverField,
	})

	test.fakeDriverField.On("Name").Return("Test").Once()
	test.fakeModelDefinition.On("GetFieldByName", "Test").Return(test.fakeFieldDefinition, nil).Once()

	test.fakeFieldDefinition.On("Copy").Return(nil).Once()

	values, err := newPayload.Mapping()
	if !test.NoError(err) {
		return
	}
	test.Len(values, 1)
}

func (test *PayloadTestSuite) TestOnScanGetFieldByNameErr() {
	newPayload := NewPayload[specs.Model]()
	tmp := newPayload.(*payload[specs.Model])

	tmp.modelDefinition = test.fakeModelDefinition

	newPayload.SetFields([]specs.DriverField{
		test.fakeDriverField,
	})

	test.fakeDriverField.On("Name").Return("Test").Once()
	test.fakeModelDefinition.On("GetFieldByName", "Test").Return(nil, errors.New("test")).Once()

	err := newPayload.OnScan([]any{})
	test.Error(err)
}

func (test *PayloadTestSuite) TestOnScan() {
	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition)
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition)

	newPayload := NewPayload[*models.CommentsModel]()

	newPayload.SetFields([]specs.DriverField{
		test.fakeDriverField,
	})

	test.fakeDriverField.On("Name").Return("Test").Once()
	test.fakeModelDefinition.On("GetFieldByName", "Test").Return(test.fakeFieldDefinition, nil).Once()
	test.fakeFieldDefinition.On("Set", "test").Return(nil).Once()

	comment := &models.CommentsModel{Content: "test"}
	test.fakeModelDefinition.On("Copy").Return(comment).Once()

	err := newPayload.OnScan([]any{"test"})
	test.NoError(err)
	test.Equal([]*models.CommentsModel{comment}, newPayload.Result())
}

func (test *PayloadTestSuite) TestJoin() {
	newPayload := NewPayload[specs.Model]()

	var t []specs.DriverJoin
	newPayload.SetJoins(t)
	test.Equal(newPayload.Join(), t)
}

func (test *PayloadTestSuite) TestLimit() {
	newPayload := NewPayload[specs.Model]()
	newPayload.SetLimit(test.fakeDriverLimit)

	test.Equal(newPayload.Limit(), test.fakeDriverLimit)
}

func (test *PayloadTestSuite) TestWhere() {
	newPayload := NewPayload[specs.Model]()
	wheres := []specs.DriverWhere{test.fakeDriverWhere}
	newPayload.SetWheres(wheres)

	test.Equal(newPayload.Where(), wheres)
}

func (test *PayloadTestSuite) TestNew() {
	comment := models.CommentsModel{}

	test.fakeUseModelDefinition.On("Use", &comment).Return(test.fakeModelDefinition).Once()
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition).Once()

	test.fakeModelDefinition.On("DatabaseName").Return("acceptance").Once()
	test.fakeModelDefinition.On("TableName").Return("comments").Once()
	test.fakeModelDefinition.On("Index").Return(0).Once()

	newPayload := NewPayload[specs.Model](&comment)
	test.Equal("acceptance", newPayload.Database())
	test.Equal("comments", newPayload.Table())
	test.Equal(0, newPayload.Index())

	test.fakeUseModelDefinition.On("Use", (*models.CommentsModel)(nil)).Return(test.fakeModelDefinition).Once()
	test.fakeModelDefinition.On("Parse").Return(test.fakeModelDefinition).Once()

	test.fakeModelDefinition.On("DatabaseName").Return("acceptance").Once()
	test.fakeModelDefinition.On("TableName").Return("comments").Once()
	test.fakeModelDefinition.On("Index").Return(0).Once()

	newWithType := NewPayload[*models.CommentsModel]()
	test.Equal("acceptance", newWithType.Database())
	test.Equal("comments", newWithType.Table())
	test.Equal(0, newWithType.Index())
}

func TestPayloadTestSuite(t *testing.T) {
	suite.Run(t, new(PayloadTestSuite))
}
