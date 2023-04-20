package drivers

import (
	"errors"
	"github.com/kitstack/dbkit/connector/drivers/joins"
	"github.com/kitstack/dbkit/tests/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type JoinTestSuite struct {
	suite.Suite
	fakeField *mocks.FakeDriverField
}

func (suite *JoinTestSuite) SetupTest() {
	suite.fakeField = mocks.NewFakeDriverField(suite.T())
}

func (suite *JoinTestSuite) TestJoinMethodDefault() {
	join := NewJoin()

	suite.Equal(joins.Method[joins.Default], join.Method())
}

func (suite *JoinTestSuite) TestJoinSetMethod() {
	join := NewJoin()

	join.SetMethod(joins.Inner)

	suite.Equal(joins.Method[joins.Inner], join.Method())
}

func (suite *JoinTestSuite) TestFromFormattedErr() {
	join := NewJoin().(*join)

	suite.fakeField.On("Formatted").Return("", errors.New("from_formatted_err")).Once()

	join.SetFrom(suite.fakeField)
	_, err := join.fromFormatted()
	suite.Error(err)

	suite.EqualValues("from_formatted_err", err.Error())
}

func (suite *JoinTestSuite) TestToFormattedErr() {
	join := NewJoin().(*join)

	suite.fakeField.On("Formatted").Return("", errors.New("to_formatted_err")).Once()

	join.SetTo(suite.fakeField)
	_, err := join.toFormatted()
	suite.Error(err)

	suite.EqualValues("to_formatted_err", err.Error())
}

func (suite *JoinTestSuite) TestJoinWithCustomField() {
	join := NewJoin()

	suite.fakeField.On("Formatted").Return("(CONCAT('%', `t0`.`name`, '%'))", nil).Once()

	join.SetFrom(suite.fakeField)

	suite.fakeField.On("IsCustom").Return(false).Once()
	suite.fakeField.On("Formatted").Return("`t1`.`id`", nil).Once()
	suite.fakeField.On("Database").Return("acceptance").Once()
	suite.fakeField.On("Table").Return("users").Once()
	suite.fakeField.On("Index").Return(1).Once()

	join.SetTo(suite.fakeField)

	formattedJoin, err := join.Formatted()
	suite.NoError(err)

	suite.Equal("JOIN `acceptance`.`users` AS `t1` ON `t1`.`id` = (CONCAT('%', `t0`.`name`, '%'))", formattedJoin)
}

func (suite *JoinTestSuite) TestJoinWithCustomFieldToErr() {
	join := NewJoin()

	suite.fakeField.On("Formatted").Return("(CONCAT('%', `t0`.`name`, '%'))", nil).Once()

	join.SetFrom(suite.fakeField)

	suite.fakeField.On("Formatted").Return("`t1`.`id`", errors.New("formatted_to_err")).Once()

	join.SetTo(suite.fakeField)

	_, err := join.Formatted()
	suite.Error(err)

	suite.EqualValues("formatted_to_err", err.Error())
}

func (suite *JoinTestSuite) TestJoinWithCustomFieldFromErr() {
	join := NewJoin()

	suite.fakeField.On("Formatted").Return("", errors.New("formatted_from_err")).Once()

	join.SetFrom(suite.fakeField)
	join.SetTo(suite.fakeField)

	_, err := join.Formatted()
	suite.Error(err)

	suite.EqualValues("formatted_from_err", err.Error())
}

func (suite *JoinTestSuite) TestJoinWithInvertedCustomField() {
	join := NewJoin()

	suite.fakeField.On("Formatted").Return("`t0`.`id`", nil).Once()

	join.SetFrom(suite.fakeField)

	suite.fakeField.On("IsCustom").Return(true).Once()
	suite.fakeField.On("Formatted").Return("(CONCAT('%', `t0`.`name`, '%'))", nil).Once()

	join.SetTo(suite.fakeField)

	formattedJoin, err := join.Formatted()
	suite.NoError(err)

	suite.Equal("JOIN (CONCAT('%', `t0`.`name`, '%')) = `t0`.`id`", formattedJoin)
}

func (suite *JoinTestSuite) TestValidateErr() {
	join := NewJoin()

	err := join.Validate()

	suite.Error(err)
	suite.IsType(&requiredFieldJoinErr{}, err)
	suite.Contains(err.Error(), "are mandatory to perform the join")

	for _, field := range []string{"From", "To"} {
		suite.Contains(err.(*requiredFieldJoinErr).Fields(), field)
	}
}

func (suite *JoinTestSuite) TestValidate() {
	join := NewJoin()

	join.SetFrom(suite.fakeField)
	join.SetTo(suite.fakeField)
	err := join.Validate()

	suite.NoError(err)
}

func TestJoinTestSuite(t *testing.T) {
	suite.Run(t, new(JoinTestSuite))
}
