package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kitstack/dbkit/tests/acceptance/fixtures"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

var isDebug bool
var test string

var ctx context.Context
var debugLog *bytes.Buffer

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	log.SetFlags(0)
	logrus.SetLevel(logrus.DebugLevel)

	debugLog = bytes.NewBuffer([]byte{})
	logrus.SetOutput(debugLog)

	ctx = context.Background()

	flag.BoolVar(&isDebug, "debug", os.Getenv("DEBUG") == "true", "Debug mode")
	flag.StringVar(&test, "test", "", "Run test by name")
	flag.Parse()
}

func main() {
	fx := fixtures.NewFixture()
	rf := reflect.ValueOf(&fx).Elem()
	typeOf := reflect.TypeOf(&fx).Elem()

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

		if test != "" && test != typeOf.Method(i).Name {
			continue
		}

		testsCount++

		state := "\x1b[32mPASS\x1b[0m"
		timerTest := time.Now()

		logrus.WithFields(logrus.Fields{
			"name": typeOf.Method(i).Name,
		}).Debug("Start test")

		fx.Reset()
		args := []reflect.Value{reflect.ValueOf(ctx)}

		// recover
		try := func() (result []reflect.Value) {
			/*		defer func() {
					if r := recover(); r != nil {
						result = make([]reflect.Value, 1)
						result[0] = reflect.ValueOf(fmt.Errorf("%v", r))
						debug.PrintStack()
					}
				}()*/
			result = method.Call(args)
			return
		}

		result := try()

		logrus.WithFields(logrus.Fields{
			"name": typeOf.Method(i).Name,
		}).Debug("End test")

		if testErr := result[0].Interface(); testErr != nil || fx.AssertErrorCount() > 0 {

			if testErr != nil {
				logrus.WithFields(logrus.Fields{
					"name": typeOf.Method(i).Name,
				}).Error(testErr)
			}

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

	fmt.Printf("\n%sDONE %d tests with %d assertions in %s | Passed: %d Failed: %d\x1b[0m\n", color, testsCount, fx.AssertCount(), time.Since(globalTimer), passedTestCount, failedTestCount)

	if isDebug || failedTestCount > 0 {
		fmt.Println("\nDebug log:")
		fmt.Println(strings.Repeat("-", 50))
		fmt.Println(debugLog.String())
	}

	if failedTestCount > 0 {
		os.Exit(1)
	}
}
