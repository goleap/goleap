package fixtures

import (
	"github.com/lab210-dev/dbkit/connector"
	"github.com/lab210-dev/dbkit/connector/config"
	"github.com/lab210-dev/dbkit/specs"
)

type Fixture struct {
	connector specs.Connector
}

func (f *Fixture) Connector() specs.Connector {
	if f.connector != nil {
		return f.connector
	}

	var err error
	f.connector, err = connector.New("acceptance",
		config.New().
			SetDriver("mysql").
			SetHost("db").
			SetUser("acceptance").
			SetPassword("acceptance").
			SetDatabase("acceptance").
			SetPort(3306),
	)

	if err != nil {
		panic(err)
	}

	return f.connector
}
