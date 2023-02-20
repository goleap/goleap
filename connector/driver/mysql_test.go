package driver

import (
	"context"
	"database/sql"
	"errors"
	"github.com/goleap/goleap/connector/config"
	"github.com/goleap/goleap/helper"
	"github.com/goleap/goleap/helper/fakesql"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MysqlTestSuite struct {
	suite.Suite
	fakeDriver  *fakesql.FakeDriver
	fakeConn    *fakesql.FakeConn
	fakeStmt    *fakesql.FakeStmt
	fakeRows    *fakesql.FakeRows
	fakePayload *helper.FakePayload
}

func (test *MysqlTestSuite) SetupSuite() {
	test.fakeDriver = fakesql.NewDriver(test.T())

	RegisteredDriver = map[string]func() Driver{
		"test": func() Driver {
			return new(Mysql)
		},
	}

	sql.Register("test", test.fakeDriver)
}

func (test *MysqlTestSuite) SetupTest() {
	test.fakeConn = fakesql.NewFakeConn(test.T())
	test.fakeStmt = fakesql.NewFakeStmt(test.T())
	test.fakeRows = fakesql.NewFakeRows(test.T())
	test.fakePayload = helper.NewFakePayload(test.T())

	test.fakeDriver.ExpectedCalls = nil
	test.fakeDriver.On("Open", ":@tcp(:0)/?parseTime=true&loc=Local").Return(test.fakeConn, nil)
}

func (test *MysqlTestSuite) TestSelectErr() {
	driver, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = driver.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.fakeConn.On("Prepare", "SELECT `t0`.`id` FROM `test` AS `t0`").Return(nil, errors.New("test")).Once()

	err = driver.Select(context.Background(), test.fakePayload)
}

/*func (test *MysqlTestSuite) TestSelect() {
	driver, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = driver.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.fakeConn.On("Prepare", "SELECT `t0`.`id` FROM `test` AS `t0`").Return(test.fakeStmt, nil).Once()
	test.fakeStmt.On("NumInput").Return(0)
	test.fakeStmt.On("Close").Return(nil)
	test.fakeStmt.On("Query", []driverSql.Value{}).Return(test.fakeRows, nil)
	test.fakeRows.On("Columns").Return([]string{"id"})
	test.fakeRows.On("Close").Return(nil)

	var line = 0
	test.fakeRows.On("Next", mock.Anything).Return(func(dest []driverSql.Value) error {
		dest[0] = 100

		if line < 1 {
			line++
			return nil
		}

		return io.EOF
	})

	n := 0
	err = driver.Select(context.Background(), Payload{
		Table:      "test",
		Database:   "test",
		Index:      0,
		Fields:     []Field{NewField().SetName("id").SetIndex(0)},
		Where:      []Where{},
		Join:       []Join{},
		ResultType: []any{&n},
		OnScan: func(result []any) error {
			n = *result[0].(*int)
			return nil
		},
	})

	if !test.Empty(err) {
		return
	}

	test.Equal(0, n)
}*/

func TestMysqlTestSuite(t *testing.T) {
	suite.Run(t, new(MysqlTestSuite))
}
