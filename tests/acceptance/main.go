package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/lab210-dev/dbkit/tests/acceptance/fixtures"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"reflect"
	"time"
)

var ctx context.Context

func init() {
	log.SetFlags(0)
	logrus.SetLevel(logrus.DebugLevel)
	tmp := bytes.NewBuffer([]byte{})
	logrus.SetOutput(tmp)

	ctx = context.Background()
}

func main() {
	fx := fixtures.Fixture{}
	rf := reflect.ValueOf(&fx)
	typeOf := reflect.TypeOf(&fx)

	testsCount := 0
	failedTestCount := 0
	passedTestCount := 0

	globalTimer := time.Now()

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

		testsCount++

		state := "\x1b[32mPASS\x1b[0m"
		timerTest := time.Now()
		args := []reflect.Value{reflect.ValueOf(ctx)}
		result := method.Call(args)

		if errVal := result[0].Interface(); errVal != nil {
			state = "\x1b[31mFAILED\x1b[0m"
			failedTestCount++
		} else {
			passedTestCount++
		}

		fmt.Printf("%s %s (%s)\n", state, typeOf.Method(i).Name, time.Since(timerTest))
	}

	var color string
	if failedTestCount > 0 {
		color = "\x1b[31m" // Rouge
	} else {
		color = "\x1b[32m" // Vert
	}

	fmt.Printf("\n%sDONE %d tests in %s | Passed: %d Failed: %d\x1b[0m\n", color, testsCount, time.Since(globalTimer), passedTestCount, failedTestCount)

	if failedTestCount > 0 {
		os.Exit(1)
	}
}
