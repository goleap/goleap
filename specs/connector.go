package specs

import (
	"context"
)

type Connector interface {
	Driver

	GetCnx(ctx context.Context)

	Config() Config

	Name() string
	SetName(name string) Connector
}
