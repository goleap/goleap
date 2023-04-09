package drivers

import (
	"errors"
	"github.com/lab210-dev/dbkit/connector/drivers/operators"
	"github.com/lab210-dev/dbkit/tests/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type WhereTestSuite struct {
	suite.Suite
	fakeDriverField *mocks.FakeDriverField
}

func (test *WhereTestSuite) SetupTest() {
	test.fakeDriverField = mocks.NewFakeDriverField(test.T())
}

func (test *WhereTestSuite) TestWhereOperator() {
	where := NewWhere()

	where.SetOperator(operators.Equal)
	test.Equal(where.Operator(), operators.Equal)
}

func (test *WhereTestSuite) TestWhereFrom() {
	where := NewWhere()

	where.SetFrom(test.fakeDriverField)
	test.NotNil(where.From())
	test.Equal(where.From(), test.fakeDriverField)
}

func (test *WhereTestSuite) TestWhereTo() {
	where := NewWhere()

	where.SetTo("test")
	test.Equal(where.To(), "test")

	where.SetTo(test.fakeDriverField)
	test.NotNil(where.To())
	test.Equal(where.To(), test.fakeDriverField)
}

func (test *WhereTestSuite) TestWhereBuildOperator() {
	where := NewWhere().(*where)

	where.SetOperator(operators.Equal)
	op, err := where.buildOperator()
	test.NoError(err)

	test.Equal(op, "= ?")

	where.SetOperator(operators.In)
	op, err = where.buildOperator()
	test.NoError(err)

	test.Equal(op, "IN (?)")

	where.SetOperator(operators.NotIn)
	op, err = where.buildOperator()
	test.NoError(err)

	test.Equal(op, "NOT IN (?)")

	where.SetOperator(operators.IsNull)
	op, err = where.buildOperator()
	test.NoError(err)

	test.Equal(op, "IS NULL")

	where.SetOperator(operators.IsNotNull)
	op, err = where.buildOperator()
	test.NoError(err)

	test.Equal(op, "IS NOT NULL")

	where.SetOperator("unknown")
	_, err = where.buildOperator()
	test.Error(err)

	test.IsType(&unknownOperatorErr{}, err)
	test.Contains(err.Error(), "unknown operator: unknown")

	_, _, err = where.Formatted()
	test.IsType(&unknownOperatorErr{}, err)
	test.Contains(err.Error(), "unknown operator: unknown")

}

func (test *WhereTestSuite) TestWhereFormatted() {
	where := NewWhere()

	test.fakeDriverField.On("Formatted").Return("`t0`.`id`", nil).Once()

	where.SetFrom(test.fakeDriverField)
	where.SetOperator(operators.Equal)
	where.SetTo("test")

	formatted, args, err := where.Formatted()

	test.NoError(err)
	test.Equal("`t0`.`id` = ?", formatted)
	test.Equal("test", args)
}

func (test *WhereTestSuite) TestWhereFormattedFromErr() {
	where := NewWhere()

	test.fakeDriverField.On("Formatted").Return("", errors.New("where_from_formatted_err")).Once()

	where.SetFrom(test.fakeDriverField)
	where.SetOperator(operators.Equal)
	where.SetTo("test")

	_, _, err := where.Formatted()

	test.Error(err)
	test.Equal(err.Error(), "where_from_formatted_err")
}

func (test *WhereTestSuite) TestWhereFormattedToErr() {
	where := NewWhere()

	test.fakeDriverField.On("Formatted").Return("`t0`.`id`", nil).Once()
	test.fakeDriverField.On("Formatted").Return("", errors.New("where_to_formatted_err")).Once()

	where.SetFrom(test.fakeDriverField)
	where.SetOperator(operators.Equal)
	where.SetTo(test.fakeDriverField)

	_, _, err := where.Formatted()

	test.Error(err)
	test.Equal(err.Error(), "where_to_formatted_err")
}

func (test *WhereTestSuite) TestWhereFormattedWithTwoField() {
	where := NewWhere()

	test.fakeDriverField.On("Formatted").Return("`t0`.`id`", nil).Once()
	test.fakeDriverField.On("Formatted").Return("`t1`.`id`", nil).Once()

	where.SetFrom(test.fakeDriverField)
	where.SetOperator(operators.Equal)
	where.SetTo(test.fakeDriverField)

	formatted, args, err := where.Formatted()

	test.NoError(err)
	test.Equal("`t0`.`id` = `t1`.`id`", formatted)
	test.Equal(nil, args)
}

func TestWhereTestSuite(t *testing.T) {
	suite.Run(t, new(WhereTestSuite))
}
