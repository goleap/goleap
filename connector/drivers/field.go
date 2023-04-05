package drivers

import (
	"fmt"
	"github.com/lab210-dev/dbkit/specs"
	"regexp"
	"strings"
)

type field struct {
	index        int
	name         string
	nameInSchema string

	// @todo: found a better name for this
	fn     string
	fnArgs []specs.DriverField
}

func (f *field) NameInModel() string {
	return f.nameInSchema
}

func (f *field) SetNameInModel(nameInSchema string) specs.DriverField {
	f.nameInSchema = nameInSchema
	return f
}

func (f *field) Index() int {
	return f.index
}

func (f *field) SetIndex(index int) specs.DriverField {
	f.index = index
	return f
}

func (f *field) Name() string {
	return f.name
}

func (f *field) SetName(name string) specs.DriverField {
	f.name = strings.TrimSpace(name)
	return f
}

func (f *field) SetFn(fn string, args []specs.DriverField) specs.DriverField {
	f.fn = fn
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
			if arg.NameInModel() != match[1] {
				continue
			}
			column, err := arg.Column()
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

func (f *field) Column() (string, error) {
	if f.Fn() != "" {
		return f.fnProcess()
	}
	return fmt.Sprintf("`t%d`.`%s`", f.Index(), f.Name()), nil
}

func NewField() specs.DriverField {
	return new(field)
}
