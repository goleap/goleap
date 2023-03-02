package drivers

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/lab210-dev/dbkit/connector/config"
	"github.com/lab210-dev/dbkit/mocks"
	"github.com/lab210-dev/dbkit/mocks/fakesql"
	"github.com/lab210-dev/dbkit/specs"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"testing"
)

type MysqlTestSuite struct {
	suite.Suite
	fakeDriver  *fakesql.FakeDriver
	fakeConn    *fakesql.FakeConn
	fakeStmt    *fakesql.FakeStmt
	fakeRows    *fakesql.FakeRows
	fakePayload *mocks.FakePayload
}

func (test *MysqlTestSuite) SetupSuite() {
	test.fakeDriver = fakesql.NewDriver(test.T())

	RegisteredDriver = map[string]func() specs.Driver{
		"test": func() specs.Driver {
			return new(Mysql)
		},
	}

	sql.Register("test", test.fakeDriver)
}

func (test *MysqlTestSuite) SetupTest() {
	test.fakeConn = fakesql.NewFakeConn(test.T())
	test.fakeStmt = fakesql.NewFakeStmt(test.T())
	test.fakeRows = fakesql.NewFakeRows(test.T())
	test.fakePayload = mocks.NewFakePayload(test.T())

	test.fakeDriver.ExpectedCalls = nil
	test.fakeDriver.On("Open", ":@tcp(:0)/?parseTime=true&loc=Local").Return(test.fakeConn, nil)
}

func (test *MysqlTestSuite) TestNew() {
	driver, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = driver.New(config.New().SetDriver(""))
	test.NotEmpty(err)
}

func (test *MysqlTestSuite) TestCreate() {
	driver, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = driver.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.NotEmpty(driver.Get())
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

	test.fakePayload.On("Fields").Return([]specs.DriverField{
		NewField().SetName("id").SetIndex(0),
		NewField().SetName("label").SetIndex(0),
	})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Index").Return(0)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id`, `t0`.`label` FROM `test` AS `t0`").Return(nil, errors.New("test")).Once()

	err = driver.Select(context.Background(), test.fakePayload)
	test.Error(err)
}

func (test *MysqlTestSuite) TestSelect() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{NewField().SetName("id").SetIndex(0)})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Index").Return(0)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id` FROM `test` AS `t0`").Return(test.fakeStmt, nil).Once()
	test.fakeStmt.On("NumInput").Return(0)
	test.fakeStmt.On("Close").Return(nil)
	test.fakeRows.On("Columns").Return([]string{"id"})
	test.fakeRows.On("Close").Return(nil)

	mapping := []any{new(uint64)}
	test.fakePayload.On("Mapping").Return(mapping)

	test.fakePayload.On("OnScan", mapping).Return(nil)

	var line = 0
	test.fakeRows.On("Next", mock.Anything).Return(func(dest []driver.Value) error {
		dest[0] = 1

		if line < 1 {
			line++
			return nil
		}

		return io.EOF
	})

	test.fakeStmt.On("Query", []driver.Value{}).Return(test.fakeRows, nil)

	err = drv.Select(context.Background(), test.fakePayload)
	test.Empty(err)
}

func (test *MysqlTestSuite) TestSelectWithNativeScanErr() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{NewField().SetName("id").SetIndex(0)})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Index").Return(0)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id` FROM `test` AS `t0`").Return(test.fakeStmt, nil).Once()
	test.fakeStmt.On("NumInput").Return(0)
	test.fakeStmt.On("Close").Return(nil)
	test.fakeRows.On("Columns").Return([]string{"id"})
	test.fakeRows.On("Close").Return(nil)

	mapping := []any{new(uint64)}
	test.fakePayload.On("Mapping").Return(mapping)

	//test.fakePayload.On("OnScan", mapping).Return(errors.New("test"))
	var line = 0
	test.fakeRows.On("Next", mock.Anything).Return(func(dest []driver.Value) error {
		dest[0] = "test"

		if line < 1 {
			line++
			return nil
		}

		return io.EOF
	})

	test.fakeStmt.On("Query", []driver.Value{}).Return(test.fakeRows, nil)

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
}

func (test *MysqlTestSuite) TestSelectWithNativeWrapScanErr() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{NewField().SetName("id").SetIndex(0)})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Index").Return(0)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id` FROM `test` AS `t0`").Return(test.fakeStmt, nil).Once()
	test.fakeStmt.On("NumInput").Return(0)
	test.fakeStmt.On("Close").Return(nil)
	test.fakeRows.On("Columns").Return([]string{"id"})
	test.fakeRows.On("Close").Return(nil)

	mapping := []any{new(uint64)}
	test.fakePayload.On("Mapping").Return(mapping)

	test.fakeStmt.On("Query", []driver.Value{}).Return(test.fakeRows, nil)

	var line = 0
	test.fakeRows.On("Next", mock.Anything).Return(func(dest []driver.Value) error {
		dest[0] = 1

		if line < 1 {
			line++
			return nil
		}

		return io.EOF
	})

	test.fakePayload.On("OnScan", mapping).Return(errors.New("test"))

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
}

func (test *MysqlTestSuite) TestSelectWithWhere() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{NewField().SetName("id").SetIndex(0)})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Index").Return(0)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id` FROM `test` AS `t0`").Return(test.fakeStmt, nil).Once()
	test.fakeStmt.On("NumInput").Return(0)
	test.fakeStmt.On("Close").Return(nil)
	test.fakeRows.On("Columns").Return([]string{"id"})
	test.fakeRows.On("Close").Return(nil)

	mapping := []any{new(uint64)}
	test.fakePayload.On("Mapping").Return(mapping)

	test.fakePayload.On("OnScan", mapping).Return(nil)

	var line = 0
	test.fakeRows.On("Next", mock.Anything).Return(func(dest []driver.Value) error {
		dest[0] = 1

		if line < 1 {
			line++
			return nil
		}

		return io.EOF
	})

	test.fakeStmt.On("Query", []driver.Value{}).Return(test.fakeRows, nil)

	err = drv.Select(context.Background(), test.fakePayload)
	test.Empty(err)
}

func TestMysqlTestSuite(t *testing.T) {
	suite.Run(t, new(MysqlTestSuite))
}
