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

var generateInArgument = sqlx.In

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

		column, err := field.Column()
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

		// TODO Maybe add specific operator for join
		result += fmt.Sprintf("%s `%s`.`%s` AS `t%d` ON `t%d`.`%s` = `t%d`.`%s`", field.Method(), field.ToDatabase(), field.ToTable(), field.ToTableIndex(), field.ToTableIndex(), field.ToKey(), field.FromTableIndex(), field.FromKey())
	}

	return
}

func (m *Mysql) buildWhere(fields []specs.DriverWhere) (result string, args []any, err error) {
	for i, field := range fields {

		operator, err := m.buildOperator(field)
		if err != nil {
			return "", nil, err
		}

		if i > 0 {
			result += " AND "
		}

		fromField, err := field.From().Column()
		if err != nil {
			return "", nil, err
		}

		result += fmt.Sprintf("%s %s", fromField, operator)

		if field.To() != nil {
			// TODO Maybe need to interpret the like a DriverField
			args = append(args, field.To())
		}
	}

	if result != "" {
		result = fmt.Sprintf("WHERE %s", result)
	}

	return
}

func (m *Mysql) buildLimit(limit specs.DriverLimit) (result string) {
	if limit == nil {
		return
	}

	return fmt.Sprintf("LIMIT %d, %d", limit.Offset(), limit.Limit())
}

// TODO Create and trigger error when operator is not supported !
func (m *Mysql) buildOperator(field specs.DriverWhere) (string, error) {
	switch field.Operator() {
	case operators.Equal, operators.NotEqual:
		return fmt.Sprintf("%s ?", field.Operator()), nil
	case operators.In, operators.NotIn:
		return fmt.Sprintf("%s (?)", field.Operator()), nil
	case operators.IsNull, operators.IsNotNull:
		return field.Operator(), nil
	}
	return "", NewUnknownOperatorErr(field.Operator())
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

	buildLimit := m.buildLimit(payload.Limit())

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

	queryWithArgs, args, err := generateInArgument(query, args...)
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
