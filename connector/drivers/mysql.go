package drivers

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lab210-dev/dbkit/connector/drivers/operators"
	"github.com/lab210-dev/dbkit/specs"
	log "github.com/sirupsen/logrus"
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

func (m *Mysql) buildFields(fields []specs.DriverField) (result string) {
	for i, field := range fields {
		if i > 0 {
			result += ", "
		}

		result += fmt.Sprintf("`t%d`.`%s`", field.Index(), field.Name())
	}

	return result
}

func (m *Mysql) buildJoin(joins []specs.DriverJoin) (result string, err error) {
	for i, field := range joins {
		if i > 0 {
			result += " "
		}

		err := field.Validate()
		if err != nil {
			return "", err
		}

		// TODO Maybe add specific operator for join
		result += fmt.Sprintf("%s `%s`.`%s` AS `t%d` ON `t%d`.`%s` = `t%d`.`%s`", field.Method(), field.ToDatabase(), field.ToTable(), field.ToTableIndex(), field.ToTableIndex(), field.ToKey(), field.FromTableIndex(), field.FromKey())
	}

	return
}

func (m *Mysql) buildWhere(fields []specs.DriverWhere) (result string, args []any, err error) {
	for i, field := range fields {

		operator := m.buildOperator(field)
		// This case is very strange because it does not generate an error.
		if operator == "" {
			return "", nil, fmt.Errorf("unknown operator: %s", field.Operator())
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
		result = fmt.Sprintf("WHERE %s", result)
	}

	return
}

func (m *Mysql) buildOperator(field specs.DriverWhere) string {
	switch field.Operator() {
	case operators.Equal, operators.NotEqual:
		return fmt.Sprintf("%s ?", field.Operator())
	case operators.In, operators.NotIn:
		return fmt.Sprintf("%s (?)", field.Operator())
	case operators.IsNull, operators.IsNotNull:
		return field.Operator()
	}
	return ""
}

func (m *Mysql) Select(ctx context.Context, payload specs.Payload) (err error) {
	builtWhere, args, err := m.buildWhere(payload.Where())
	if err != nil {
		return
	}

	builtJoin, err := m.buildJoin(payload.Join())
	if err != nil {
		return
	}

	builtFields := m.buildFields(payload.Fields())

	query := fmt.Sprintf("SELECT %s FROM `%s` AS `t%d`", builtFields, payload.Table(), payload.Index())

	if builtJoin != "" {
		query += fmt.Sprintf(" %s", builtJoin)
	}

	if builtWhere != "" {
		query += fmt.Sprintf(" %s", builtWhere)
	}

	queryWithArgs, args, err := sqlx.In(query, args...)
	if err != nil {
		return
	}

	log.WithFields(log.Fields{
		"type":  "select",
		"query": queryWithArgs,
		"args":  args,
	}).Debug("Execute: Select()")

	rows, err := m.db.QueryContext(ctx, queryWithArgs, args...)
	if err != nil {
		return
	}

	return wrapScan(rows, payload.Mapping(), payload.OnScan)
}

func (m *Mysql) Get() *sql.DB {
	return m.db
}
