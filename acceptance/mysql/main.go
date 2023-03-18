package main

import (
	"context"
	"github.com/lab210-dev/dbkit"
	"github.com/lab210-dev/dbkit/acceptance/mysql/model"
	"github.com/lab210-dev/dbkit/connector"
	"github.com/lab210-dev/dbkit/connector/config"
	log "github.com/sirupsen/logrus"
)

var ctx context.Context

func init() {
	log.SetLevel(log.DebugLevel)
	ctx = context.Background()
}

func main() {
	acceptanceConnector, err := connector.New("acceptance",
		config.New().
			SetDriver("mysql").
			SetHost("localhost").
			SetUser("root").
			SetPassword("onlyfordev").
			SetDatabase("acceptance").
			SetPort(3333),
	)
	if err != nil {
		panic(err)
	}

	// TODO (Lab210-dev) : Wrap with interface for testing.
	err = acceptanceConnector.Get().Ping()
	if err != nil {
		panic(err)
	}

	// Test With Automate Field Builder
	users, err := dbkit.Use[*model.UserModel](ctx, acceptanceConnector).
		Fields("Id").
		Get(1)

	if err != nil {
		panic(err)
	}

	log.Println(users.Id == 1)

	// Test Without Automate Field Builder
	// TODO (Lab210-dev) : Maybe schema is conditionally required.
	/*
		sch := schema.New(&model.UserModel{}).Parse()
		selectPayload := dbkit.NewPayload[*model.UserModel](sch)
		err = acceptanceConnector.Select(ctx, selectPayload)
		if err != nil {
			panic(err)
		}
	*/
}
