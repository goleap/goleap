package specs

import "context"

type ConnectionManager interface {
	GetConnection(ctx context.Context) (connection Connection, err error)
}
