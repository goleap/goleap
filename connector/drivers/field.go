package drivers

import (
	"fmt"
	"github.com/lab210-dev/dbkit/specs"
	"regexp"
	"strings"
)

type field struct {
	index    int
	column   string
	database string
	table    string
	name     string

	// @todo: found a better name for this
	fn     string
	fnArgs []specs.DriverField
}

func (f *field) Table() string {
	return f.table
}

func (f *field) SetTable(name string) specs.DriverField {
	f.table = strings.TrimSpace(name)
	return f
}

func (f *field) Database() string {
	return f.database
}

func (f *field) SetDatabase(name string) specs.DriverField {
	f.database = strings.TrimSpace(name)
	return f
}

func (f *field) Name() string {
	return f.name
}

func (f *field) SetName(name string) specs.DriverField {
	f.name = strings.TrimSpace(name)
	return f
}

func (f *field) Index() int {
	return f.index
}

func (f *field) SetIndex(index int) specs.DriverField {
	f.index = index
	return f
}

func (f *field) Column() string {
	return f.column
}

func (f *field) SetColumn(name string) specs.DriverField {
	f.column = strings.TrimSpace(name)
	return f
}

func (f *field) SetFn(fn string, args []specs.DriverField) specs.DriverField {
	f.fn = fmt.Sprintf("(%s)", fn)
	f.fnArgs = args
	return f
}

func (f *field) Fn() string {
	return f.fn
}

func (f *field) FnArgs() []specs.DriverField {
	return f.fnArgs
}

func (f *field) fnProcess() (fn string, err error) {
	re := regexp.MustCompile(`%([a-zA-Z_]+)%`)
	matches := re.FindAllStringSubmatch(f.Fn(), -1)

	replaceFieldCount := 0
	for _, match := range matches {
		for _, arg := range f.FnArgs() {
			if arg.Name() != match[1] {
				continue
			}
			column, err := arg.Formatted()
			if err != nil {
				return "", err
			}

			fn = strings.Replace(f.Fn(), match[0], column, -1)
			replaceFieldCount++
			break
		}
	}

	if replaceFieldCount != len(matches) {
		noFoundFiels := make([]string, 0)
		for _, match := range matches {
			noFoundFiels = append(noFoundFiels, match[1])
		}
		return "", NewUnknownFieldsErr(noFoundFiels)
	}

	return
}

func (f *field) Formatted() (string, error) {
	if f.Fn() != "" {
		return f.fnProcess()
	}
	return fmt.Sprintf("`t%d`.`%s`", f.Index(), f.Column()), nil
}

func NewField() specs.DriverField {
	return new(field)
}
