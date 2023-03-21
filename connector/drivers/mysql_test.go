package drivers

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/lab210-dev/dbkit/connector/config"
	"github.com/lab210-dev/dbkit/specs"
	"github.com/lab210-dev/dbkit/tests/mocks"
	fakesql2 "github.com/lab210-dev/dbkit/tests/mocks/fakesql"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"testing"
)

type MysqlTestSuite struct {
	suite.Suite
	fakeDriver  *fakesql2.FakeDriver
	fakeConn    *fakesql2.FakeConn
	fakeStmt    *fakesql2.FakeStmt
	fakeRows    *fakesql2.FakeRows
	fakePayload *mocks.FakePayload
}

func (test *MysqlTestSuite) SetupSuite() {
	test.fakeDriver = fakesql2.NewDriver(test.T())

	RegisteredDriver = map[string]func() specs.Driver{
		"test": func() specs.Driver {
			return new(Mysql)
		},
	}

	sql.Register("test", test.fakeDriver)
}

func (test *MysqlTestSuite) SetupTest() {
	test.fakeConn = fakesql2.NewFakeConn(test.T())
	test.fakeStmt = fakesql2.NewFakeStmt(test.T())
	test.fakeRows = fakesql2.NewFakeRows(test.T())
	test.fakePayload = mocks.NewFakePayload(test.T())

	test.fakeDriver.ExpectedCalls = nil
	test.fakeDriver.On("Open", ":@tcp(:0)/?parseTime=true&loc=Local").Return(test.fakeConn, nil)
}

func (test *MysqlTestSuite) TestNew() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver(""))
	test.NotEmpty(err)
}

func (test *MysqlTestSuite) TestCreate() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.NotEmpty(drv.Get())
}

func (test *MysqlTestSuite) TestSelectErr() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test"))
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

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
}

func (test *MysqlTestSuite) TestSimpleWhere() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{
		NewField().SetName("id").SetIndex(0),
		NewField().SetName("label").SetIndex(0),
	})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		NewWhere().SetFrom(NewField().SetName("id").SetIndex(0)).SetOperator(EqualOperator).SetTo(1),
	})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Index").Return(0)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id`, `t0`.`label` FROM `test` AS `t0` WHERE `t0`.`id` = ?").Return(nil, errors.New("test")).Once()

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
}

func (test *MysqlTestSuite) TestSimpleWhereWithBadOperator() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{
		NewField().SetName("id").SetIndex(0),
		NewField().SetName("label").SetIndex(0),
	})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		NewWhere().SetFrom(NewField().SetName("id").SetIndex(0)).SetOperator("").SetTo(1),
	})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Index").Return(0)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id`, `t0`.`label` FROM `test` AS `t0`").Return(nil, errors.New("test")).Once()

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
}

func (test *MysqlTestSuite) TestSimpleWhereIsNullOperator() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{
		NewField().SetName("id").SetIndex(0),
		NewField().SetName("label").SetIndex(0),
	})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		NewWhere().SetFrom(NewField().SetName("id").SetIndex(0)).SetOperator(IsNullOperator),
	})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Index").Return(0)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id`, `t0`.`label` FROM `test` AS `t0` WHERE `t0`.`id` IS NULL").Return(nil, errors.New("test")).Once()

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
}

func (test *MysqlTestSuite) TestSimpleWhereInOperator() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{
		NewField().SetName("id").SetIndex(0),
		NewField().SetName("label").SetIndex(0),
	})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		NewWhere().SetFrom(NewField().SetName("id").SetIndex(0)).SetOperator(InOperator).SetTo([]int{1, 2}),
	})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Index").Return(0)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id`, `t0`.`label` FROM `test` AS `t0` WHERE `t0`.`id` IN (?)").Return(nil, errors.New("test")).Once()

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
}

func (test *MysqlTestSuite) TestMultiSimpleWhere() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{
		NewField().SetName("id").SetIndex(0),
		NewField().SetName("label").SetIndex(0),
	})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		NewWhere().SetFrom(NewField().SetName("id").SetIndex(0)).SetOperator(EqualOperator).SetTo(1),
		NewWhere().SetFrom(NewField().SetName("label").SetIndex(0)).SetOperator(EqualOperator).SetTo("test"),
	})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Index").Return(0)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id`, `t0`.`label` FROM `test` AS `t0` WHERE `t0`.`id` = ? AND `t0`.`label` = ?").Return(nil, errors.New("test")).Once()

	err = drv.Select(context.Background(), test.fakePayload)
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
