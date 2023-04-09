package drivers

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type LimitTestSuite struct {
	suite.Suite
}

func (suite *LimitTestSuite) TestLimit() {
	limit := NewLimit()
	suite.Equal(0, limit.Offset())
	suite.Equal(0, limit.Limit())

	limit.SetOffset(1)
	suite.Equal(1, limit.Offset())

	limit.SetLimit(2)
	suite.Equal(2, limit.Limit())

	formatted, err := limit.Formatted()
	suite.NoError(err)
	suite.Equal("LIMIT 1, 2", formatted)
}

func TestLimitTestSuite(t *testing.T) {
	suite.Run(t, new(LimitTestSuite))
}
