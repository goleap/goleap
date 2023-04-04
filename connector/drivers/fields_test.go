package drivers

import (
	"github.com/lab210-dev/dbkit/specs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FieldTestSuite struct {
	suite.Suite
	field *field
}

func (suite *FieldTestSuite) SetupTest() {
	suite.field = NewField().(*field)
	suite.field.SetIndex(1)
	suite.field.SetName("name")
	suite.field.SetNameInModel("name_in_model")
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

func TestFieldTestSuite(t *testing.T) {
	suite.Run(t, new(FieldTestSuite))
}
