package main

import (
	"context"
	"github.com/lab210-dev/dbkit/tests/acceptance/fixtures"
	"github.com/sirupsen/logrus"
	"log"
	"reflect"
	"strings"
	"time"
)

var ctx context.Context

func init() {
	log.SetFlags(0)
	logrus.SetLevel(logrus.DebugLevel)
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

		log.Print(strings.Repeat("-", 100))
		log.Printf("Running fixture : %s", typeOf.Method(i).Name)
		log.Print("Debug :")
		log.Println()

		timer := time.Now()
		args := []reflect.Value{reflect.ValueOf(ctx)}
		result := method.Call(args)

		log.Println()
		log.Printf("Ending fixture `%s` in %s", typeOf.Method(i).Name, time.Since(timer))
		log.Print(strings.Repeat("-", 100))

		if errVal := result[0].Interface(); errVal != nil {
			err := errVal.(error)
			log.Printf("Returned an error: %s", err)
			continue
		}

		log.Println()
	}
}
