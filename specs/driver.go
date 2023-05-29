package specs

import (
	"context"
)

type Driver interface {
	New(config Config) error
	Manager() ConnectionManager

	Select(ctx context.Context, payload Payload) error
}
