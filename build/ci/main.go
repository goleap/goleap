//
// Use dagger to run acceptance tests
// https://docs.dagger.io/
//

package main

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"strings"
)

var goVersion string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	fmt.Println("Environment variables loaded")

	data, err := os.ReadFile(".gvmrc")
	if err != nil {
		panic(err)
	}
	goVersion = strings.TrimSpace(string(data))
}

func main() {
	if err := execute(context.Background()); err != nil {
		panic(err)
	}
}

func execute(ctx context.Context) (err error) {
	fmt.Println("Acceptance tests")

	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			log.Print("Received SIGINT")
			cancel()
		case <-ctx.Done():
		}
	}()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}

	defer func(client *dagger.Client) {
		err = client.Close()
	}(client)

	dir := client.Host().Directory(".", dagger.HostDirectoryOpts{
		Exclude: []string{"/vendor"},
	})

	db := client.Container().
		From("mysql:8").
		WithEnvVariable("MYSQL_USER", os.Getenv("MYSQL_USER")).
		WithEnvVariable("MYSQL_PASSWORD", os.Getenv("MYSQL_PASSWORD")).
		WithEnvVariable("MYSQL_DATABASE", os.Getenv("MYSQL_DATABASE")).
		WithEnvVariable("MYSQL_ROOT_PASSWORD", os.Getenv("MYSQL_ROOT_PASSWORD")).
		WithMountedFile("/docker-entrypoint-initdb.d/dump.sql", dir.File("/tests/acceptance/data/acceptance.sql")).
		WithExposedPort(3306).
		WithExec(nil)

	_, err = client.
		Container().
		From(fmt.Sprintf("golang:%s", goVersion)).
		WithServiceBinding("db", db).
		WithEnvVariable("MYSQL_HOST", os.Getenv("MYSQL_HOST")).
		WithEnvVariable("MYSQL_USER", os.Getenv("MYSQL_USER")).
		WithEnvVariable("MYSQL_PASSWORD", os.Getenv("MYSQL_PASSWORD")).
		WithEnvVariable("MYSQL_DATABASE", os.Getenv("MYSQL_DATABASE")).
		WithEnvVariable("MYSQL_ROOT_PASSWORD", os.Getenv("MYSQL_ROOT_PASSWORD")).
		WithEnvVariable("DEBUG", os.Getenv("DEBUG")).
		WithMountedDirectory("/src", dir).
		WithWorkdir("/src").
		WithExec(append([]string{"go", "run", "./tests/acceptance"}, os.Args[1:]...)).Stdout(ctx)

	return
}
