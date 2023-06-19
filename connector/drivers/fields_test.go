package drivers

import (
	"errors"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/dbkit/tests/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FieldTestSuite struct {
	suite.Suite
	field     *field
	fakeField *mocks.FakeDriverField
}

func (suite *FieldTestSuite) SetupTest() {
	suite.field = NewField().(*field)
	suite.field.SetIndex(1)
	suite.field.SetColumn("name")
	suite.field.SetName("name_in_model")

	suite.fakeField = mocks.NewFakeDriverField(suite.T())
}

func (suite *FieldTestSuite) TestTable() {
	suite.field.SetTable("t1")
	suite.Equal("t1", suite.field.Table())
}

func (suite *FieldTestSuite) TestDatabase() {
	suite.field.SetDatabase("db1")
	suite.Equal("db1", suite.field.Database())
}

func (suite *FieldTestSuite) TestColumn() {
	column, err := suite.field.Formatted()

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "`t1`.`name`", column)
}

func (suite *FieldTestSuite) TestColumnWithFn() {
	suite.field.SetCustom("CONCAT('%', ${Name}, '%')", []specs.DriverField{NewField().SetName("Name").SetColumn("name")})

	column, err := suite.field.Formatted()

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "(CONCAT('%', `t0`.`name`, '%'))", column)
	assert.True(suite.T(), suite.field.IsCustom())
}

func (suite *FieldTestSuite) TestColumnWithFnErrNoMatch() {
	suite.field.SetCustom("CONCAT('%', ${Name}, '%')", []specs.DriverField{NewField().SetName("unknown").SetColumn("name")})

	column, err := suite.field.Formatted()

	assert.Equal(suite.T(), "", column)
	assert.Error(suite.T(), err)

	assert.IsType(suite.T(), &ErrUnknownFields{}, err)
	assert.EqualValues(suite.T(), "unknown fields: Name", err.Error())
	assert.Equal(suite.T(), []string{"Name"}, err.(specs.ErrUnknownFields).Fields())
}

func (suite *FieldTestSuite) TestColumnWithFnColumnErr() {
	suite.fakeField.On("Name").Return("Name")
	suite.fakeField.On("Formatted").Return("", errors.New("column_error"))

	suite.field.SetCustom("CONCAT('%', ${Name}, '%')", []specs.DriverField{suite.fakeField})

	column, err := suite.field.Formatted()

	assert.Equal(suite.T(), "", column)
	assert.Error(suite.T(), err)
	assert.EqualValues(suite.T(), "column_error", err.Error())
}

func TestFieldTestSuite(t *testing.T) {
	suite.Run(t, new(FieldTestSuite))
}
