package fixtures

import (
	"fmt"
	"github.com/lab210-dev/dbkit/connector"
	"github.com/lab210-dev/dbkit/connector/config"
	"github.com/lab210-dev/dbkit/specs"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
)

type Fixture struct {
	connector        specs.Connector
	assert           *assert.Assertions
	assertErrorCount int
}

func (fixture *Fixture) AssertErrorCount() int {
	return fixture.assertErrorCount
}

func (fixture *Fixture) Reset() {
	fixture.assertErrorCount = 0
}

func (fixture *Fixture) Assert() *assert.Assertions {
	if fixture.assert != nil {
		return fixture.assert
	}
	fixture.assert = assert.New(fixture)
	return fixture.assert
}

func (fixture *Fixture) Errorf(format string, args ...any) {
	fixture.assertErrorCount++
	fmt.Printf(format, args...)
}

func (fixture *Fixture) Connector() specs.Connector {
	if fixture.connector != nil {
		return fixture.connector
	}
	var err error
	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		port = 3306
	}

	fixture.connector, err = connector.New("acceptance",
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

	return fixture.connector
}
