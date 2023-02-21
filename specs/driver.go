package specs

import (
	"context"
	"database/sql"
)

type Driver interface {
	New(config Config) error
	Get() *sql.DB

	Select(ctx context.Context, payload Payload) error
}
