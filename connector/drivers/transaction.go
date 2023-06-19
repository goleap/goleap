package drivers

import (
	"context"
	"database/sql"
	"time"
)

type transaction struct {
	tx         *sql.Tx
	connection *connection
}

func wrapTransaction(connection *connection, tx *sql.Tx) *transaction {
	return &transaction{
		tx:         tx,
		connection: connection.SetIsTransaction(true),
	}
}

func (t *transaction) Tx() *sql.Tx {
	return t.tx
}

func (t *transaction) Connection() *connection {
	return t.connection
}

func (t *transaction) Commit() error {
	defer t.Connection().debugQuery("Commit", "", time.Now(), nil)
	return t.Tx().Commit()
}

func (t *transaction) Rollback() error {
	defer t.Connection().debugQuery("Rollback", "", time.Now(), nil)
	return t.Tx().Rollback()
}

func (t *transaction) ExecContext(ctx context.Context, query string, args ...any) (result sql.Result, err error) {
	deferFn := func() {
		t.Connection().debugQuery("QueryRowContext", query, time.Now(), err, args...)
	}

	defer func() {
		if err == nil {
			deferFn()
		}
	}()

	result, err = t.Tx().ExecContext(ctx, query, args...)

	if err != nil {
		return nil, t.catchTransactionError(err, deferFn)
	}

	return result, nil
}

func (t *transaction) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	var err error
	defer t.Connection().debugQuery("QueryRowContext", query, time.Now(), err, args...)
	rows, err := t.Tx().QueryContext(ctx, query, args...)

	if err != nil {
		return nil, t.Connection().catchError(err)
	}

	return rows, nil
}

func (t *transaction) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	var err error
	defer t.Connection().debugQuery("QueryRowContext", query, time.Now(), err, args...)

	rows := t.Tx().QueryRowContext(ctx, query, args...)
	err = rows.Err()

	if err != nil {
		_ = t.Connection().catchError(err)
		return nil
	}

	return rows
}

func (t *transaction) catchTransactionError(err error, fnDefer func()) error {
	fnDefer()

	if err == context.Canceled || err == context.DeadlineExceeded {
		return t.Connection().catchError(err)
	}

	_ = t.Rollback()

	return err
}
