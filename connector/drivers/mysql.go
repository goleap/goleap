package drivers

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/depkit"
	log "github.com/sirupsen/logrus"
)

func init() {
	depkit.Register[specs.SqlIn](sqlx.In)
}

type Mysql struct {
	specs.Config
	db *sql.DB
}

func (m *Mysql) New(config specs.Config) (err error) {
	m.Config = config

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

func (m *Mysql) buildFields(fields []specs.DriverField) (result string, err error) {
	for i, field := range fields {
		if i > 0 {
			result += ", "
		}

		column, err := field.Formatted()
		if err != nil {
			return "", err
		}

		result += column
	}

	return result, nil
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

		formatted, err := field.Formatted()
		if err != nil {
			return "", err
		}

		result += formatted
	}

	return
}

func (m *Mysql) buildWhere(wheres []specs.DriverWhere) (result string, args []any, err error) {
	for i, where := range wheres {

		if i > 0 {
			result += " AND "
		}

		formatted, whereArgs, err := where.Formatted()
		if err != nil {
			return "", nil, err
		}

		result += formatted

		if whereArgs != nil {
			args = append(args, whereArgs...)
		}
	}

	if result != "" {
		result = fmt.Sprintf("WHERE %s", result)
	}

	return
}

func (m *Mysql) buildLimit(limit specs.DriverLimit) (result string, err error) {
	if limit == nil {
		return
	}

	return limit.Formatted()
}

func (m *Mysql) Select(ctx context.Context, payload specs.Payload) (err error) {
	buildFields, err := m.buildFields(payload.Fields())
	if err != nil {
		return
	}

	builtWhere, args, err := m.buildWhere(payload.Where())
	if err != nil {
		return
	}

	builtJoin, err := m.buildJoin(payload.Join())
	if err != nil {
		return
	}

	buildLimit, err := m.buildLimit(payload.Limit())
	if err != nil {
		return
	}

	query := fmt.Sprintf("SELECT %s FROM `%s`.`%s` AS `t%d`", buildFields, m.Database(), payload.Table(), payload.Index())

	if builtJoin != "" {
		query += fmt.Sprintf(" %s", builtJoin)
	}

	if builtWhere != "" {
		query += fmt.Sprintf(" %s", builtWhere)
	}

	if buildLimit != "" {
		query += fmt.Sprintf(" %s", buildLimit)
	}

	queryWithArgs, args, err := depkit.Get[specs.SqlIn]()(query, args...)
	if err != nil {
		return
	}

	mapping, err := payload.Mapping()
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

	return wrapScan(rows, mapping, payload.OnScan)
}

func (m *Mysql) Get() *sql.DB {
	return m.db
}
