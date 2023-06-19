package drivers

import (
	"context"
	"database/sql"
	"github.com/kitstack/dbkit/specs"
	"github.com/sirupsen/logrus"
	"time"
)

type connection struct {
	db *sql.DB
	context.Context
	connectionId  int64
	conn          *sql.Conn
	isTransaction bool
}

func (connectionInstance *connection) Id() int64 {
	return connectionInstance.connectionId
}

func (connectionInstance *connection) PingContext(ctx context.Context) error {
	return connectionInstance.Conn().PingContext(ctx)
}

func (connectionInstance *connection) Conn() *sql.Conn {
	return connectionInstance.conn
}

func (connectionInstance *connection) SetIsTransaction(isTransaction bool) *connection {
	connectionInstance.isTransaction = isTransaction
	return connectionInstance
}

func WrapConnection(ctx context.Context, idb *sql.DB) (specs.Connection, error) {
	db := &connection{
		Context: ctx,
		db:      idb,
	}

	var err error
	db.conn, err = idb.Conn(ctx)
	if err != nil {
		return nil, err
	}

	db.connectionId = DbCacheConnectionInstance.GetConnectionId(db.conn)
	if db.connectionId == 0 {
		err = db.QueryRowContext(ctx, `SELECT CONNECTION_ID()`).Scan(&db.connectionId)
		if err != nil {
			return nil, err
		}

		DbCacheConnectionInstance.RegisterConn(db.conn, db.connectionId)
	}

	return db, nil
}

func (connectionInstance *connection) Close() error {
	logrus.WithFields(logrus.Fields{
		"kind":         "Close",
		"connectionId": connectionInstance.connectionId,
	}).Debug("MySQL Connection Close")
	return connectionInstance.Conn().Close()
}

func (connectionInstance *connection) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	var err error
	defer connectionInstance.debugQuery("QueryRowContext", query, time.Now(), err, args...)
	rows, err := connectionInstance.Conn().QueryContext(ctx, query, args...)

	if err != nil {
		return nil, connectionInstance.catchError(err)
	}

	return rows, nil
}

func (connectionInstance *connection) debugQuery(kind string, query string, start time.Time, err error, args ...any) {
	logrus.WithFields(logrus.Fields{
		"kind":          kind,
		"connectionId":  connectionInstance.connectionId,
		"transactional": connectionInstance.isTransaction,
		"query":         query,
		"args":          args,
		"duration":      time.Since(start),
		"err":           err,
	}).Debug("MySQL Connection Query")
}

func (connectionInstance *connection) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	var err error
	defer connectionInstance.debugQuery("QueryRowContext", query, time.Now(), err, args...)

	rows := connectionInstance.Conn().QueryRowContext(ctx, query, args...)
	err = rows.Err()
	if err != nil {
		_ = connectionInstance.catchError(err)
		return nil
	}

	return rows
}

func (connectionInstance *connection) BeginTx(ctx context.Context, opts *sql.TxOptions) (specs.Transaction, error) {
	tx, err := connectionInstance.Conn().BeginTx(ctx, opts)

	if err != nil {
		_ = connectionInstance.catchError(err)
		return nil, err
	}

	return wrapTransaction(connectionInstance, tx), nil
}

func (connectionInstance *connection) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	var err error
	defer connectionInstance.debugQuery("ExecContext", query, time.Now(), err, args...)

	result, err := connectionInstance.Conn().ExecContext(ctx, query, args...)

	if err != nil {
		return nil, connectionInstance.catchError(err)
	}

	return result, nil
}

func (connectionInstance *connection) Kill() error {
	query := "KILL ?"
	args := []any{DbCacheConnectionInstance.GetConnectionId(connectionInstance.conn)}

	defer connectionInstance.debugQuery("Kill", query, time.Now(), nil, args...)

	_, err := connectionInstance.db.Exec(query, args...)
	return err
}

func (connectionInstance *connection) catchError(err error) error {
	if err == context.Canceled || err == context.DeadlineExceeded {
		killErr := connectionInstance.Kill()
		if killErr != nil {
			return killErr
		}
		return err
	}

	return err
}
