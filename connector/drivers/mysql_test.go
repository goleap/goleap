package drivers

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lab210-dev/dbkit/connector/config"
	"github.com/lab210-dev/dbkit/connector/drivers/joins"
	"github.com/lab210-dev/dbkit/connector/drivers/operators"
	"github.com/lab210-dev/dbkit/specs"
	"github.com/lab210-dev/dbkit/tests/mocks"
	"github.com/lab210-dev/dbkit/tests/mocks/fakesql"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"log"
	"testing"
)

type MysqlTestSuite struct {
	suite.Suite
	fakeDriver  *fakesql.FakeDriver
	fakeConn    *fakesql.FakeConn
	fakeStmt    *fakesql.FakeStmt
	fakeRows    *fakesql.FakeRows
	fakeIn      *mocks.FakeIn
	fakePayload *mocks.FakePayload

	fakeDriverLimit *mocks.FakeDriverLimit
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
	test.fakeIn = mocks.NewFakeIn(test.T())
	test.fakeDriverLimit = mocks.NewFakeDriverLimit(test.T())

	// Resetting with default function
	generateInArgument = sqlx.In

	test.fakeDriver.ExpectedCalls = nil
	test.fakeDriver.On("Open", ":@tcp(:3306)/acceptance?parseTime=true&loc=Local").Return(test.fakeConn, nil)
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

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Mapping").Return([]any{}, nil)
	test.fakePayload.On("Fields").Return([]specs.DriverField{
		NewField().SetName("id").SetIndex(0),
		NewField().SetName("email").SetIndex(0),
	})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Table").Return("users")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(test.fakeDriverLimit)

	test.fakeDriverLimit.On("Offset").Return(0)
	test.fakeDriverLimit.On("Limit").Return(1)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id`, `t0`.`email` FROM `acceptance`.`users` AS `t0` LIMIT 0, 1").Return(nil, errors.New("test")).Once()

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
}

func (test *MysqlTestSuite) TestSelectMappingErr() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Mapping").Return([]any{}, errors.New("test_mapping_err"))
	test.fakePayload.On("Fields").Return([]specs.DriverField{
		NewField().SetName("id").SetIndex(0),
		NewField().SetName("email").SetIndex(0),
	})

	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Table").Return("users")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	err = drv.Select(context.Background(), test.fakePayload)
	test.ErrorContains(err, "test_mapping_err")
}

func (test *MysqlTestSuite) TestSimpleWhere() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{
		NewField().SetName("id").SetIndex(0),
		NewField().SetName("email").SetIndex(0),
	})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		NewWhere().SetFrom(NewField().SetName("id").SetIndex(0)).SetOperator(operators.Equal).SetTo(1),
	})
	test.fakePayload.On("Mapping").Return([]any{}, nil)
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Table").Return("users")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id`, `t0`.`email` FROM `acceptance`.`users` AS `t0` WHERE `t0`.`id` = ?").Return(nil, errors.New("test")).Once()

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
}

func (test *MysqlTestSuite) TestSimpleWhereWithBadOperator() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		NewWhere().SetFrom(NewField().SetName("id").SetIndex(0)).SetOperator("").SetTo(1),
	})

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
}

func (test *MysqlTestSuite) TestSimpleWhereIsNullOperator() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{
		NewField().SetName("id").SetIndex(0),
		NewField().SetName("label").SetIndex(0),
	})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		NewWhere().SetFrom(NewField().SetName("id").SetIndex(0)).SetOperator(operators.IsNull),
	})

	test.fakePayload.On("Mapping").Return([]any{}, nil)
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id`, `t0`.`label` FROM `acceptance`.`test` AS `t0` WHERE `t0`.`id` IS NULL").Return(nil, errors.New("test")).Once()

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
}

func (test *MysqlTestSuite) TestSimpleWhereInOperator() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{
		NewField().SetName("id").SetIndex(0),
		NewField().SetName("label").SetIndex(0),
	})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		NewWhere().SetFrom(NewField().SetName("id").SetIndex(0)).SetOperator(operators.In).SetTo([]int{1, 2}),
	})
	test.fakePayload.On("Mapping").Return([]any{}, nil)
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id`, `t0`.`label` FROM `acceptance`.`test` AS `t0` WHERE `t0`.`id` IN (?, ?)").Return(nil, errors.New("test")).Once()

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
}

func (test *MysqlTestSuite) TestSimpleWhereInOperatorErr() {
	// Override the generateInArgument function to return an error
	generateInArgument = test.fakeIn.GenerateInArgument

	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{
		NewField().SetName("id").SetIndex(0),
		NewField().SetName("label").SetIndex(0),
	})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		NewWhere().SetFrom(NewField().SetName("id").SetIndex(0)).SetOperator(operators.In).SetTo([]int{1, 2}),
	})
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	fnErrorMsg := "function `GenerateInArgument` returns an error"
	test.fakeIn.On("GenerateInArgument", "SELECT `t0`.`id`, `t0`.`label` FROM `acceptance`.`test` AS `t0` WHERE `t0`.`id` IN (?)", []int{1, 2}).Return("", []any{}, errors.New(fnErrorMsg)).Once()

	err = drv.Select(context.Background(), test.fakePayload)
	test.EqualError(err, fnErrorMsg)
}

func (test *MysqlTestSuite) TestWhereMultiEqual() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{
		NewField().SetName("id").SetIndex(0),
		NewField().SetName("email").SetIndex(0),
	})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		NewWhere().SetFrom(NewField().SetName("id").SetIndex(0)).SetOperator(operators.Equal).SetTo(1),
		NewWhere().SetFrom(NewField().SetName("email").SetIndex(0)).SetOperator(operators.Equal).SetTo("test"),
	})

	test.fakePayload.On("Mapping").Return([]any{}, nil)
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Table").Return("users")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id`, `t0`.`email` FROM `acceptance`.`users` AS `t0` WHERE `t0`.`id` = ? AND `t0`.`email` = ?").Return(nil, errors.New("test")).Once()

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
}

func (test *MysqlTestSuite) TestSelect() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{NewField().SetName("id").SetIndex(0)})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id` FROM `acceptance`.`test` AS `t0`").Return(test.fakeStmt, nil).Once()
	test.fakeStmt.On("NumInput").Return(0)
	test.fakeStmt.On("Close").Return(nil)
	test.fakeRows.On("Columns").Return([]string{"id"})
	test.fakeRows.On("Close").Return(nil)

	mapping := []any{new(uint64)}
	test.fakePayload.On("Mapping").Return(mapping, nil)

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

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{NewField().SetName("id").SetIndex(0)})
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Table").Return("users")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id` FROM `acceptance`.`users` AS `t0`").Return(test.fakeStmt, nil).Once()
	test.fakeStmt.On("NumInput").Return(0)
	test.fakeStmt.On("Close").Return(nil)
	test.fakeRows.On("Columns").Return([]string{"id"})
	test.fakeRows.On("Close").Return(nil)

	mapping := []any{new(uint64)}
	test.fakePayload.On("Mapping").Return(mapping, nil)

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

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{NewField().SetName("id").SetIndex(0)})
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Table").Return("users")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id` FROM `acceptance`.`users` AS `t0`").Return(test.fakeStmt, nil).Once()
	test.fakeStmt.On("NumInput").Return(0)
	test.fakeStmt.On("Close").Return(nil)
	test.fakeRows.On("Columns").Return([]string{"id"})
	test.fakeRows.On("Close").Return(nil)

	mapping := []any{new(uint64)}
	test.fakePayload.On("Mapping").Return(mapping, nil)

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

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{NewField().SetName("id").SetIndex(0)})
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Table").Return("users")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id` FROM `acceptance`.`users` AS `t0`").Return(test.fakeStmt, nil).Once()
	test.fakeStmt.On("NumInput").Return(0)
	test.fakeStmt.On("Close").Return(nil)
	test.fakeRows.On("Columns").Return([]string{"id"})
	test.fakeRows.On("Close").Return(nil)

	mapping := []any{new(uint64)}
	test.fakePayload.On("Mapping").Return(mapping, nil)

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

func (test *MysqlTestSuite) TestJoin() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Mapping").Return([]any{}, nil)
	test.fakePayload.On("Fields").Return([]specs.DriverField{NewField().SetName("id").SetIndex(0).SetNameInSchema("Id")})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Join").Return([]specs.DriverJoin{
		NewJoin().
			SetMethod(joins.Default).
			SetFromKey("posts_id").
			SetFromTableIndex(0).
			SetToTable("posts").
			SetToKey("id").
			SetToTableIndex(1).
			SetToDatabase("acceptance"),
	})
	test.fakePayload.On("Limit").Return(nil)
	test.fakePayload.On("Table").Return("comments")
	test.fakePayload.On("Index").Return(0)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id` FROM `acceptance`.`comments` AS `t0` JOIN `acceptance`.`posts` AS `t1` ON `t1`.`id` = `t0`.`posts_id`").Return(nil, errors.New("test")).Once()

	err = drv.Select(context.Background(), test.fakePayload)
	test.NotEmpty(err)
}

func (test *MysqlTestSuite) TestMultiJoin() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Mapping").Return([]any{}, nil)
	test.fakePayload.On("Fields").Return([]specs.DriverField{NewField().SetName("id").SetIndex(0).SetNameInSchema("Id")})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Join").Return([]specs.DriverJoin{
		NewJoin().
			SetMethod(joins.Default).
			SetFromKey("posts_id").
			SetFromTableIndex(0).
			SetToTable("posts").
			SetToKey("id").
			SetToTableIndex(1).
			SetToDatabase("acceptance"),

		NewJoin().
			SetMethod(joins.Default).
			SetFromKey("users_id").
			SetFromTableIndex(0).
			SetToTable("users").
			SetToKey("id").
			SetToTableIndex(2).
			SetToDatabase("acceptance"),
	})
	test.fakePayload.On("Limit").Return(nil)
	test.fakePayload.On("Table").Return("comments")
	test.fakePayload.On("Index").Return(0)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id` FROM `acceptance`.`comments` AS `t0` JOIN `acceptance`.`posts` AS `t1` ON `t1`.`id` = `t0`.`posts_id` JOIN `acceptance`.`users` AS `t2` ON `t2`.`id` = `t0`.`users_id`").Return(nil, errors.New("test")).Once()

	err = drv.Select(context.Background(), test.fakePayload)
	test.NotEmpty(err)
}

func (test *MysqlTestSuite) TestJoinErr() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Join").Return([]specs.DriverJoin{
		NewJoin().
			SetMethod(joins.Default).
			SetFromTableIndex(0).
			SetToTable("comments").
			SetToKey("comments_id").
			SetToTableIndex(1).
			SetToDatabase("blog"),
	})

	err = drv.Select(context.Background(), test.fakePayload)
	log.Print(err)
	test.NotEmpty(err)
}

func TestMysqlTestSuite(t *testing.T) {
	suite.Run(t, new(MysqlTestSuite))
}
