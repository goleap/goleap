package dbkit

import (
	"context"
	"errors"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/dbkit/tests/mocks"
	"github.com/kitstack/dbkit/tests/models"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PayloadTestSuite struct {
	suite.Suite
	context.Context
	fakeModelDefinition *mocks.FakeModelDefinition
	fakeFieldDefinition *mocks.FakeFieldDefinition
	fakeDriverField     *mocks.FakeDriverField
	fakeDriverLimit     *mocks.FakeDriverLimit
}

func (test *PayloadTestSuite) SetupTest() {
	test.Context = context.Background()
	test.fakeModelDefinition = mocks.NewFakeModelDefinition(test.T())
	test.fakeFieldDefinition = mocks.NewFakeFieldDefinition(test.T())
	test.fakeDriverField = mocks.NewFakeDriverField(test.T())
	test.fakeDriverLimit = mocks.NewFakeDriverLimit(test.T())
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

func (test *PayloadTestSuite) TestJoin() {
	newPayload := NewPayload[specs.Model]()
	tmp := newPayload.(*payload[specs.Model])

	var t []specs.DriverJoin
	test.Equal(tmp.Join(), t)
}

func (test *PayloadTestSuite) TestLimit() {
	newPayload := NewPayload[specs.Model]()
	newPayload.SetLimit(test.fakeDriverLimit)

	test.Equal(newPayload.Limit(), test.fakeDriverLimit)
}

func (test *PayloadTestSuite) TestNew() {
	comment := models.CommentsModel{}
	newPayload := NewPayload[specs.Model](&comment)
	test.Equal("acceptance", newPayload.Database())
	test.Equal("comments", newPayload.Table())

	newWithType := NewPayload[*models.CommentsModel]()
	test.Equal("acceptance", newWithType.Database())
	test.Equal("comments", newWithType.Table())
}

func TestPayloadTestSuite(t *testing.T) {
	suite.Run(t, new(PayloadTestSuite))
}
