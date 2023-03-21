package main

import (
	"context"
	"github.com/lab210-dev/dbkit/tests/acceptance/fixtures"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"time"
)

var ctx context.Context

func init() {
	log.SetLevel(log.DebugLevel)
	ctx = context.Background()
}

func main() {
	fx := fixtures.Fixture{}
	rf := reflect.ValueOf(&fx)
	typeOf := reflect.TypeOf(&fx)

	for i := 0; i < rf.NumMethod(); i++ {

		method := rf.Method(i)

		if !method.IsValid() {
			continue
		}

		if method.Type().NumIn() != 1 || method.Type().NumOut() != 1 {
			continue
		}

		if method.Type().In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() {
			continue
		}

		if method.Type().Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}

		log.Println()
		log.Printf("Running fixture : %s", typeOf.Method(i).Name)
		log.Print(strings.Repeat("-", 100))

		timer := time.Now()
		args := []reflect.Value{reflect.ValueOf(ctx)}
		result := method.Call(args)

		log.Print(strings.Repeat("-", 100))
		log.Printf("Ending fixture `%s` in %s", typeOf.Method(i).Name, time.Since(timer))

		if errVal := result[0].Interface(); errVal != nil {
			err := errVal.(error)
			log.Errorf("Returned an error: %s", err)
			continue
		}

		log.Println()
	}
}
