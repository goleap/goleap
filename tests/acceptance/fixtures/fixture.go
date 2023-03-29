package fixtures

import (
	"github.com/lab210-dev/dbkit/connector"
	"github.com/lab210-dev/dbkit/connector/config"
	"github.com/lab210-dev/dbkit/specs"
	"os"
	"strconv"
)

type Fixture struct {
	connector specs.Connector
}

func (f *Fixture) Connector() specs.Connector {
	if f.connector != nil {
		return f.connector
	}
	var err error
	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		port = 3306
	}

	f.connector, err = connector.New("acceptance",
		config.New().
			SetDriver("mysql").
			SetHost(os.Getenv("MYSQL_HOST")).
			SetUser(os.Getenv("MYSQL_USER")).
			SetPassword(os.Getenv("MYSQL_PASSWORD")).
			SetDatabase(os.Getenv("MYSQL_DATABASE")).
			SetPort(port),
	)

	if err != nil {
		panic(err)
	}

	return f.connector
}
