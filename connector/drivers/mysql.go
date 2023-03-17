package drivers

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lab210-dev/dbkit/specs"
)

type Mysql struct {
	db *sql.DB
}

func (m *Mysql) New(config specs.Config) (err error) {
	dataSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=%s",
		config.User(),
		config.Password(),
		config.Host(),
		config.Port(),
		config.Database(),
		config.Locale(),
	)

	m.db, err = sql.Open(config.Driver(), dataSource)
	if err != nil {
		return
	}

	return
}

func (m *Mysql) buildField(fields []specs.DriverField) (result string) {
	for i, field := range fields {
		if i > 0 {
			result += ", "
		}

		result += fmt.Sprintf("`t%d`.`%s`", field.Index(), field.Name())
	}

	return result
}

func (m *Mysql) buildWhere(fields []specs.DriverWhere) (result string, args []any) {
	for i, field := range fields {

		operator := m.buildOperator(field)
		// This case is very strange because it does not generate an error.
		if operator == "" {
			continue
		}

		if i > 0 {
			result += " AND "
		}

		result += fmt.Sprintf("`t%d`.`%s` %s", field.From().Index(), field.From().Name(), operator)

		if field.To() != nil {
			args = append(args, field.To())
		}
	}

	if result != "" {
		result = fmt.Sprintf(" WHERE %s", result)
	}

	return
}

func (m *Mysql) buildOperator(field specs.DriverWhere) string {
	switch field.Operator() {
	case EqualOperator, NotEqualOperator:
		return fmt.Sprintf("%s ?", field.Operator())
	case InOperator, NotInOperator:
		return fmt.Sprintf("%s (?)", field.Operator())
	case IsNullOperator, IsNotNullOperation:
		return field.Operator()
	}
	return ""
}

func (m *Mysql) Select(ctx context.Context, payload specs.Payload) (err error) {
	buildWhere, _ := m.buildWhere(payload.Where())
	query := fmt.Sprintf("SELECT %s FROM `%s` AS `t%d`%s", m.buildField(payload.Fields()), payload.Table(), payload.Index(), buildWhere)

	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return
	}

	return wrapScan(rows, payload.Mapping(), payload.OnScan)
}

func (m *Mysql) Get() *sql.DB {
	return m.db
}
