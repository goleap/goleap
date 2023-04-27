package drivers

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type DriverTestSuite struct {
	suite.Suite
}

func (test *DriverTestSuite) TestGetDriver() {
	drv, err := Get("mysql")
	if !test.Empty(err) {
		return
	}
	test.Equal(drv, &Mysql{})
}

func (test *DriverTestSuite) TestGetUnknownDriver() {
	_, err := Get("unknown")
	test.EqualError(err, "driver not found")
}

func TestDriverTestSuite(t *testing.T) {
	suite.Run(t, new(DriverTestSuite))
}
