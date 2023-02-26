package drivers

import (
	"context"
	"database/sql"
	"fmt"
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

func (m *Mysql) Select(ctx context.Context, payload specs.Payload) (err error) {
	query := fmt.Sprintf("SELECT %s FROM `%s` AS `t%d`", m.buildField(payload.Fields()), payload.Table(), payload.Index())

	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return
	}

	return wrapScan(rows, payload.Mapping(), payload.OnScan)
}

func (m *Mysql) Get() *sql.DB {
	return m.db
}
