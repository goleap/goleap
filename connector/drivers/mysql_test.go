package drivers

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/kitstack/dbkit/connector/config"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/dbkit/tests/mocks"
	"github.com/kitstack/dbkit/tests/mocks/fakesql"
	"github.com/kitstack/depkit"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"strings"
	"testing"
)

type MysqlTestSuite struct {
	suite.Suite
	fakeDriver *fakesql.FakeDriver
	fakeConn   *fakesql.FakeConn
	fakeStmt   *fakesql.FakeStmt
	fakeRows   *fakesql.FakeRows

	fakeIn          *mocks.FakeIn
	fakePayload     *mocks.FakePayload
	fakeDriverField *mocks.FakeDriverField

	fakeDriverLimit *mocks.FakeDriverLimit
	fakeDriverJoin  *mocks.FakeDriverJoin
	fakeDriverWhere *mocks.FakeDriverWhere
	fakeSqlIn       *mocks.FakeSqlIn
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
	test.fakeDriverField = mocks.NewFakeDriverField(test.T())
	test.fakeDriverJoin = mocks.NewFakeDriverJoin(test.T())
	test.fakeDriverWhere = mocks.NewFakeDriverWhere(test.T())
	test.fakeSqlIn = mocks.NewFakeSqlIn(test.T())

	depkit.Reset()
	depkit.Register[specs.SqlIn](test.fakeSqlIn.Execute)

	test.fakeDriver.ExpectedCalls = nil
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

func (test *MysqlTestSuite) TestBuildFieldsErr() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.fakeDriverField.On("Formatted").Return("", errors.New("build_field_column_err"))

	_, err = drv.(*Mysql).buildFields([]specs.DriverField{test.fakeDriverField})
	test.Error(err)
	test.EqualValues("build_field_column_err", err.Error())
}

func (test *MysqlTestSuite) TestSelectBuildFieldsErr() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{test.fakeDriverField})

	test.fakeDriverField.On("Formatted").Return("", errors.New("build_field_column_err"))

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
	test.EqualValues("build_field_column_err", err.Error())
}

func (test *MysqlTestSuite) TestBuildLimit() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.fakeDriverLimit.On("Formatted").Return("LIMIT 0, 1", nil)

	limitValue, err := drv.(*Mysql).buildLimit(test.fakeDriverLimit)
	test.NoError(err)
	test.EqualValues("LIMIT 0, 1", limitValue)
}

func (test *MysqlTestSuite) TestBuildJoin() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.fakeDriverJoin.On("Validate").Return(nil)
	test.fakeDriverJoin.On("Formatted").Return("", errors.New("build_join_validate_err"))

	_, err = drv.(*Mysql).buildJoin([]specs.DriverJoin{test.fakeDriverJoin})
	test.Error(err)
	test.EqualValues("build_join_validate_err", err.Error())
}

func (test *MysqlTestSuite) TestBuildWhere() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test"))
	if !test.Empty(err) {
		return
	}

	test.fakeDriverWhere.On("Formatted").Return("", nil, errors.New("build_where_err"))

	wheres := []specs.DriverWhere{
		test.fakeDriverWhere,
	}
	// test.fakePayload.On("SetWheres").Return(wheres)

	_, _, err = drv.(*Mysql).buildWhere(wheres)
	test.Error(err)
	test.EqualValues("build_where_err", err.Error())
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

	test.fakeDriver.On("Open", ":@tcp(:3306)/acceptance?parseTime=true&loc=Local").Return(test.fakeConn, nil)

	test.fakePayload.On("Mapping").Return([]any{}, nil)

	test.fakeDriverField.On("Formatted").Return("`t0`.`id`", nil).Once()
	test.fakeDriverField.On("Formatted").Return("`t0`.`email`", nil).Once()
	test.fakePayload.On("Fields").Return([]specs.DriverField{
		test.fakeDriverField,
		test.fakeDriverField,
	})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Database").Return("acceptance")
	test.fakePayload.On("Table").Return("users")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(test.fakeDriverLimit)

	test.fakeDriverLimit.On("Formatted").Return("LIMIT 0, 1", nil)

	query := "SELECT `t0`.`id`, `t0`.`email` FROM `acceptance`.`users` AS `t0` LIMIT 0, 1"
	test.fakeSqlIn.On("Execute", query).Return(query, []any{}, nil)
	test.fakeConn.On("Prepare", query).Return(nil, errors.New("test")).Once()

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

	test.fakeDriver.On("Open", ":@tcp(:3306)/acceptance?parseTime=true&loc=Local").Return(test.fakeConn, nil)

	test.fakePayload.On("Mapping").Return([]any{}, errors.New("test_mapping_err"))

	test.fakeDriverField.On("Formatted").Return("`t0`.`id`", nil).Once()
	test.fakeDriverField.On("Formatted").Return("`t0`.`email`", nil).Once()
	test.fakePayload.On("Fields").Return([]specs.DriverField{
		test.fakeDriverField,
		test.fakeDriverField,
	})

	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Database").Return("acceptance")
	test.fakePayload.On("Table").Return("users")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	query := "SELECT `t0`.`id`, `t0`.`email` FROM `acceptance`.`users` AS `t0`"
	test.fakeSqlIn.On("Execute", query).Return(query, []any{}, nil)

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

	test.fakeDriver.On("Open", ":@tcp(:3306)/acceptance?parseTime=true&loc=Local").Return(test.fakeConn, nil)

	test.fakeDriverField.On("Formatted").Return("`t0`.`id`", nil).Once()
	test.fakePayload.On("Fields").Return([]specs.DriverField{
		test.fakeDriverField,
	})

	test.fakeDriverWhere.On("Formatted").Return("`t0`.`id` = ?", []any{1}, nil)
	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		test.fakeDriverWhere,
	})
	test.fakePayload.On("Mapping").Return([]any{}, nil)
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Table").Return("users")
	test.fakePayload.On("Database").Return("acceptance")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	query := "SELECT `t0`.`id` FROM `acceptance`.`users` AS `t0` WHERE `t0`.`id` = ?"
	test.fakeSqlIn.On("Execute", query, 1).Return(query, []any{1}, nil)

	test.fakeConn.On("Prepare", query).Return(nil, errors.New("test")).Once()

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

	test.fakePayload.On("Fields").Return([]specs.DriverField{})

	test.fakeDriverWhere.On("Formatted").Return("", nil, errors.New("build_where_err"))
	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		test.fakeDriverWhere,
	})

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
	test.EqualValues("build_where_err", err.Error())
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

	test.fakeDriver.On("Open", ":@tcp(:3306)/acceptance?parseTime=true&loc=Local").Return(test.fakeConn, nil)

	test.fakeDriverField.On("Formatted").Return("`t0`.`id`", nil).Once()
	test.fakePayload.On("Fields").Return([]specs.DriverField{
		test.fakeDriverField,
	})

	test.fakeDriverWhere.On("Formatted").Return("`t0`.`id` IS NULL", nil, nil)
	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		test.fakeDriverWhere,
	})

	test.fakePayload.On("Mapping").Return([]any{}, nil)
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Database").Return("acceptance")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	query := "SELECT `t0`.`id` FROM `acceptance`.`test` AS `t0` WHERE `t0`.`id` IS NULL"
	test.fakeSqlIn.On("Execute", query).Return(query, []any{}, nil)

	test.fakeConn.On("Prepare", "SELECT `t0`.`id` FROM `acceptance`.`test` AS `t0` WHERE `t0`.`id` IS NULL").Return(nil, errors.New("test")).Once()

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

	test.fakeDriver.On("Open", ":@tcp(:3306)/acceptance?parseTime=true&loc=Local").Return(test.fakeConn, nil)

	test.fakeDriverField.On("Formatted").Return("`t0`.`id`", nil).Once()
	test.fakePayload.On("Fields").Return([]specs.DriverField{
		test.fakeDriverField,
	})

	test.fakeDriverWhere.On("Formatted").Return("`t0`.`id` IN (?)", nil, nil)
	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		test.fakeDriverWhere,
	})

	test.fakePayload.On("Mapping").Return([]any{}, nil)
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Database").Return("acceptance")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	query := "SELECT `t0`.`id` FROM `acceptance`.`test` AS `t0` WHERE `t0`.`id` IN (?)"
	test.fakeSqlIn.On("Execute", query).Return(strings.Replace(query, "?", "?, ?", -1), []any{}, nil)

	test.fakeConn.On("Prepare", strings.Replace(query, "?", "?, ?", -1)).Return(nil, errors.New("test")).Once()

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
}

func (test *MysqlTestSuite) TestSimpleWhereInOperatorErr() {
	// Override the generateInArgument function to return an error
	//generateInArgument = test.fakeIn.GenerateInArgument

	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakeDriverField.On("Formatted").Return("`t0`.`id`", nil).Once()
	test.fakePayload.On("Fields").Return([]specs.DriverField{
		test.fakeDriverField,
	})

	test.fakeDriverWhere.On("Formatted").Return("`t0`.`id` IN (?)", []any{[]int{1, 2}}, nil)
	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		test.fakeDriverWhere,
	})

	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Database").Return("acceptance")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	fnErrorMsg := "function `GenerateInArgument` returns an error"
	query := "SELECT `t0`.`id` FROM `acceptance`.`test` AS `t0` WHERE `t0`.`id` IN (?)"
	test.fakeSqlIn.On("Execute", query, []int{1, 2}).Return(strings.Replace(query, "?", "?, ?", -1), []any{}, errors.New(fnErrorMsg))

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

	test.fakeDriver.On("Open", ":@tcp(:3306)/acceptance?parseTime=true&loc=Local").Return(test.fakeConn, nil)

	test.fakeDriverField.On("Formatted").Return("`t0`.`id`", nil).Once()
	test.fakePayload.On("Fields").Return([]specs.DriverField{
		test.fakeDriverField,
	})

	test.fakeDriverWhere.On("Formatted").Return("`t0`.`id` = ?", nil, nil).Once()
	test.fakeDriverWhere.On("Formatted").Return("`t0`.`email` = ?", nil, nil).Once()
	test.fakePayload.On("Where").Return([]specs.DriverWhere{
		test.fakeDriverWhere,
		test.fakeDriverWhere,
	})

	test.fakePayload.On("Mapping").Return([]any{}, nil)
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Database").Return("acceptance")
	test.fakePayload.On("Table").Return("users")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	query := "SELECT `t0`.`id` FROM `acceptance`.`users` AS `t0` WHERE `t0`.`id` = ? AND `t0`.`email` = ?"
	test.fakeSqlIn.On("Execute", query).Return(query, []any{}, nil)

	test.fakeConn.On("Prepare", query).Return(nil, errors.New("test")).Once()

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

	test.fakeDriver.On("Open", ":@tcp(:3306)/acceptance?parseTime=true&loc=Local").Return(test.fakeConn, nil)

	test.fakeDriverField.On("Formatted").Return("`t0`.`id`", nil).Once()
	test.fakePayload.On("Fields").Return([]specs.DriverField{
		test.fakeDriverField,
	})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Database").Return("acceptance")
	test.fakePayload.On("Table").Return("test")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	query := "SELECT `t0`.`id` FROM `acceptance`.`test` AS `t0`"
	test.fakeSqlIn.On("Execute", query).Return(query, []any{}, nil)
	test.fakeConn.On("Prepare", query).Return(test.fakeStmt, nil).Once()
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

	test.fakeDriver.On("Open", ":@tcp(:3306)/acceptance?parseTime=true&loc=Local").Return(test.fakeConn, nil)

	test.fakeDriverField.On("Formatted").Return("`t0`.`id`", nil).Once()
	test.fakePayload.On("Fields").Return([]specs.DriverField{
		test.fakeDriverField,
	})

	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Database").Return("acceptance")
	test.fakePayload.On("Table").Return("users")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	query := "SELECT `t0`.`id` FROM `acceptance`.`users` AS `t0`"
	test.fakeSqlIn.On("Execute", query).Return(query, []any{}, nil)

	test.fakeConn.On("Prepare", query).Return(test.fakeStmt, nil).Once()
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

	test.fakeDriver.On("Open", ":@tcp(:3306)/acceptance?parseTime=true&loc=Local").Return(test.fakeConn, nil)

	test.fakeDriverField.On("Formatted").Return("`t0`.`id`", nil).Once()
	test.fakePayload.On("Fields").Return([]specs.DriverField{
		test.fakeDriverField,
	})

	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Database").Return("acceptance")
	test.fakePayload.On("Table").Return("users")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	query := "SELECT `t0`.`id` FROM `acceptance`.`users` AS `t0`"
	test.fakeSqlIn.On("Execute", query).Return(query, []any{}, nil)

	test.fakeConn.On("Prepare", query).Return(test.fakeStmt, nil).Once()
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

	test.fakeDriver.On("Open", ":@tcp(:3306)/acceptance?parseTime=true&loc=Local").Return(test.fakeConn, nil)

	test.fakeDriverField.On("Formatted").Return("`t0`.`id`", nil).Once()
	test.fakePayload.On("Fields").Return([]specs.DriverField{
		test.fakeDriverField,
	})

	test.fakePayload.On("Join").Return([]specs.DriverJoin{})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})
	test.fakePayload.On("Database").Return("acceptance")
	test.fakePayload.On("Table").Return("users")
	test.fakePayload.On("Index").Return(0)
	test.fakePayload.On("Limit").Return(nil)

	query := "SELECT `t0`.`id` FROM `acceptance`.`users` AS `t0`"
	test.fakeSqlIn.On("Execute", query).Return(query, []any{}, nil)

	test.fakeConn.On("Prepare", query).Return(test.fakeStmt, nil).Once()
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

	test.fakeDriver.On("Open", ":@tcp(:3306)/acceptance?parseTime=true&loc=Local").Return(test.fakeConn, nil)

	test.fakeDriverField.On("Formatted").Return("`t0`.`id`", nil).Once()
	test.fakePayload.On("Fields").Return([]specs.DriverField{
		test.fakeDriverField,
	})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})

	test.fakeDriverJoin.On("Validate").Return(nil)
	test.fakeDriverJoin.On("Formatted").Return("JOIN `acceptance`.`posts` AS `t1` ON `t1`.`id` = `t0`.`posts_id`", nil)
	test.fakePayload.On("Join").Return([]specs.DriverJoin{
		test.fakeDriverJoin,
	})

	test.fakePayload.On("Database").Return("acceptance")
	test.fakePayload.On("Limit").Return(nil)
	test.fakePayload.On("Table").Return("comments")
	test.fakePayload.On("Index").Return(0)

	query := "SELECT `t0`.`id` FROM `acceptance`.`comments` AS `t0` JOIN `acceptance`.`posts` AS `t1` ON `t1`.`id` = `t0`.`posts_id`"
	test.fakeSqlIn.On("Execute", query).Return(query, []any{}, nil)

	test.fakePayload.On("Mapping").Return([]any{}, nil)

	test.fakeConn.On("Prepare", query).Return(nil, errors.New("test")).Once()

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

	test.fakeDriver.On("Open", ":@tcp(:3306)/acceptance?parseTime=true&loc=Local").Return(test.fakeConn, nil)

	test.fakePayload.On("Mapping").Return([]any{}, nil)

	test.fakeDriverField.On("Formatted").Return("`t0`.`id`", nil).Once()
	test.fakePayload.On("Fields").Return([]specs.DriverField{
		test.fakeDriverField,
	})

	test.fakePayload.On("Where").Return([]specs.DriverWhere{})

	test.fakeDriverJoin.On("Validate").Return(nil)
	test.fakeDriverJoin.On("Formatted").Return("JOIN `acceptance`.`posts` AS `t1` ON `t1`.`id` = `t0`.`posts_id`", nil).Once()
	test.fakeDriverJoin.On("Formatted").Return("JOIN `acceptance`.`users` AS `t2` ON `t2`.`id` = `t0`.`users_id`", nil).Once()

	test.fakePayload.On("Join").Return([]specs.DriverJoin{
		test.fakeDriverJoin,
		test.fakeDriverJoin,
	})

	test.fakePayload.On("Database").Return("acceptance")
	test.fakePayload.On("Limit").Return(nil)
	test.fakePayload.On("Table").Return("comments")
	test.fakePayload.On("Index").Return(0)

	query := "SELECT `t0`.`id` FROM `acceptance`.`comments` AS `t0` JOIN `acceptance`.`posts` AS `t1` ON `t1`.`id` = `t0`.`posts_id` JOIN `acceptance`.`users` AS `t2` ON `t2`.`id` = `t0`.`users_id`"
	test.fakeSqlIn.On("Execute", query).Return(query, []any{}, nil)

	test.fakeConn.On("Prepare", query).Return(nil, errors.New("test")).Once()

	err = drv.Select(context.Background(), test.fakePayload)
	test.NotEmpty(err)
}

func (test *MysqlTestSuite) TestJoinFormattedErr() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})

	test.fakeDriverJoin.On("Validate").Return(nil)
	test.fakeDriverJoin.On("Formatted").Return("", errors.New("join_formatted_err"))

	test.fakePayload.On("Join").Return([]specs.DriverJoin{
		test.fakeDriverJoin,
	})

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
	test.Contains(err.Error(), "join_formatted_err")
}

func (test *MysqlTestSuite) TestJoinValidateErr() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})

	test.fakeDriverJoin.On("Validate").Return(errors.New("select_join_validate_err"))

	test.fakePayload.On("Join").Return([]specs.DriverJoin{
		test.fakeDriverJoin,
	})

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
	test.Contains(err.Error(), "select_join_validate_err")
}

func (test *MysqlTestSuite) TestLimitErr() {
	drv, err := Get("test")
	if !test.Empty(err) {
		return
	}

	err = drv.New(config.New().SetDriver("test").SetDatabase("acceptance"))
	if !test.Empty(err) {
		return
	}

	test.fakePayload.On("Fields").Return([]specs.DriverField{})
	test.fakePayload.On("Where").Return([]specs.DriverWhere{})

	test.fakeDriverJoin.On("Validate").Return(nil)
	test.fakeDriverJoin.On("Formatted").Return("", nil)

	test.fakePayload.On("Join").Return([]specs.DriverJoin{
		test.fakeDriverJoin,
	})

	test.fakePayload.On("Limit").Return(test.fakeDriverLimit)
	test.fakeDriverLimit.On("Formatted").Return("", errors.New("select_limit_formatted_err"))

	err = drv.Select(context.Background(), test.fakePayload)
	test.Error(err)
	test.Contains(err.Error(), "select_limit_formatted_err")
}

func TestMysqlTestSuite(t *testing.T) {
	suite.Run(t, new(MysqlTestSuite))
}
