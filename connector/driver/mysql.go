package driver

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/goleap/goleap/connector/config"
)

type Mysql struct {
	db *sql.DB
}

func (m *Mysql) New(config config.Config) (err error) {
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

func (m *Mysql) Create() {
	//TODO implement me
	panic("implement me")
}

func (m *Mysql) buildField(fields []Field) (result string) {
	for i, field := range fields {
		if i > 0 {
			result += ", "
		}

		result += fmt.Sprintf("`t%d`.`%s`", field.Index(), field.Name())
	}

	return result
}

func (m *Mysql) Select(ctx context.Context, payload Payload) (err error) {
	query := fmt.Sprintf("SELECT %s FROM `%s` AS `t%d`", m.buildField(payload.Fields()), payload.Table(), payload.Index())

	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return
	}

	return wrapScan(rows, payload.Mapping(), payload.OnScan)
}

func (m *Mysql) Update() {
	//TODO implement me
	panic("implement me")
}

func (m *Mysql) Delete() {
	//TODO implement me
	panic("implement me")
}

func (m *Mysql) Count() {
	//TODO implement me
	panic("implement me")
}

func (m *Mysql) Get() *sql.DB {
	return m.db
}
