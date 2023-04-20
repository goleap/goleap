package drivers

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/kitstack/dbkit/specs"
)

type Value driver.Value

// RegisteredDriver is a map of registered driver, in test mode, this map will be replaced.
var RegisteredDriver = map[string]func() specs.Driver{
	"mysql": func() specs.Driver {
		return new(Mysql)
	},
}

// ErrDriverNotFound is an error when the driver is not found.
var ErrDriverNotFound = errors.New("driver not found")

// Get is a helper function to get the driver.
func Get(driver string) (specs.Driver, error) {
	drv, ok := RegisteredDriver[driver]
	if !ok {
		return nil, ErrDriverNotFound
	}

	return drv(), nil
}

// wrapScan is a helper function to wrap the sql.Rows.Scan() function.
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
