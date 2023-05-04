package fixtures

import (
	"fmt"
	"github.com/kitstack/dbkit/connector"
	"github.com/kitstack/dbkit/connector/config"
	"github.com/kitstack/dbkit/connectors"
	"github.com/kitstack/dbkit/specs"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"strings"
)

type Fixture struct {
	connector        specs.Connector
	assert           *assert.Assertions
	assertErrorCount int
	assertCount      int
}

func (fixture *Fixture) AssertErrorCount() int {
	return fixture.assertErrorCount
}

func (fixture *Fixture) AssertCount() int {
	return fixture.assertCount
}

func (fixture *Fixture) Reset() {
	fixture.assertErrorCount = 0
}

func (fixture *Fixture) Assert() *assert.Assertions {
	fixture.assertCount++
	if fixture.assert != nil {
		return fixture.assert
	}
	fixture.assert = assert.New(fixture)
	return fixture.assert
}

func (fixture *Fixture) Errorf(format string, args ...any) {
	fixture.assertErrorCount++

	format = strings.ReplaceAll(format, "\n", "\n│")
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
			SetName("acceptance").
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

	err = connectors.Instance().Add(fixture.connector)

	if err != nil {
		panic(err)
	}

	return fixture.connector
}

func NewFixture() *Fixture {
	fixture := new(Fixture)
	fixture.Connector()

	return fixture
}
