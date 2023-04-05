package drivers

import (
	"errors"
	"github.com/lab210-dev/dbkit/specs"
	"github.com/lab210-dev/dbkit/tests/mocks"
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
	suite.field.SetName("name")
	suite.field.SetNameInModel("name_in_model")

	suite.fakeField = mocks.NewFakeDriverField(suite.T())
}

func (suite *FieldTestSuite) TestColumn() {
	column, err := suite.field.Column()

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "`t1`.`name`", column)
}

func (suite *FieldTestSuite) TestColumnWithFn() {
	suite.field.SetFn("CONCAT('%', %Name%, '%')", []specs.DriverField{NewField().SetNameInModel("Name").SetName("name")})

	column, err := suite.field.Column()

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "CONCAT('%', `t0`.`name`, '%')", column)
}

func (suite *FieldTestSuite) TestColumnWithFnErrNoMatch() {
	suite.field.SetFn("CONCAT('%', %Name%, '%')", []specs.DriverField{NewField().SetNameInModel("unknown").SetName("name")})

	column, err := suite.field.Column()

	assert.Equal(suite.T(), "", column)
	assert.Error(suite.T(), err)

	assert.IsType(suite.T(), &unknownFieldsErr{}, err)
	assert.EqualValues(suite.T(), "unknown fields: Name", err.Error())
	assert.Equal(suite.T(), []string{"Name"}, err.(specs.UnknownFieldsErr).Fields())
}

func (suite *FieldTestSuite) TestColumnWithFnColumnErr() {
	suite.fakeField.On("NameInModel").Return("Name")
	suite.fakeField.On("Column").Return("", errors.New("column_error"))

	suite.field.SetFn("CONCAT('%', %Name%, '%')", []specs.DriverField{suite.fakeField})

	column, err := suite.field.Column()

	assert.Equal(suite.T(), "", column)
	assert.Error(suite.T(), err)
	assert.EqualValues(suite.T(), "column_error", err.Error())
}

func TestFieldTestSuite(t *testing.T) {
	suite.Run(t, new(FieldTestSuite))
}
