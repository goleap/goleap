package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/goleap/goleap/connector/config"
)

type Value driver.Value

var RegisteredDriver = map[string]func() Driver{
	"Mysql": func() Driver {
		return new(Mysql)
	},
}

type Payload interface {
	Table() string
	Database() string
	Index() int
	Fields() []Field
	Join() []Join
	Where() []Where
	Mapping() []any
	OnScan([]any) error
}

type Driver interface {
	New(config config.Config) error
	Get() *sql.DB

	Select(ctx context.Context, payload Payload) error
}

var ErrDriverNotFound = errors.New("driver not found")

func Get(driver string) (Driver, error) {
	drv, ok := RegisteredDriver[driver]
	if !ok {
		return nil, ErrDriverNotFound
	}

	return drv(), nil
}

func wrapScan(rows *sql.Rows, resultType []any, onScan func([]any) error) (err error) {
	defer rows.Close()

	for rows.Next() {
		tmp := make([]any, len(resultType))
		copy(tmp, resultType)

		err = rows.Scan(tmp...)
		if err != nil {
			return
		}

		err = onScan(tmp)
		if err != nil {
			return
		}
	}
	return
}
