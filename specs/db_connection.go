package specs

import (
	"context"
	"database/sql"
)

type Connection interface {
	Id() int64
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Transaction, error)
	PingContext(ctx context.Context) error
	Kill() error
	Close() error
}
