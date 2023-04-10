package drivers

import (
	"fmt"
	"github.com/lab210-dev/dbkit/specs"
	"regexp"
	"strings"
)

// field is a struct that implements the specs.DriverField interface
type field struct {
	index    int
	column   string
	database string
	table    string
	name     string

	custom     string
	customArgs []specs.DriverField
}

// Table returns the table name
func (f *field) Table() string {
	return f.table
}

// SetTable sets the table name
func (f *field) SetTable(name string) specs.DriverField {
	f.table = strings.TrimSpace(name)
	return f
}

// Database returns the database name
func (f *field) Database() string {
	return f.database
}

// SetDatabase sets the database name
func (f *field) SetDatabase(name string) specs.DriverField {
	f.database = strings.TrimSpace(name)
	return f
}

// Name returns the name
func (f *field) Name() string {
	return f.name
}

// SetName sets the name
func (f *field) SetName(name string) specs.DriverField {
	f.name = strings.TrimSpace(name)
	return f
}

// Index returns the index
func (f *field) Index() int {
	return f.index
}

// SetIndex sets the index
func (f *field) SetIndex(index int) specs.DriverField {
	f.index = index
	return f
}

// Column returns the column name
func (f *field) Column() string {
	return f.column
}

// SetColumn sets the column name
func (f *field) SetColumn(name string) specs.DriverField {
	f.column = strings.TrimSpace(name)
	return f
}

// SetCustom sets the custom function
func (f *field) SetCustom(fn string, args []specs.DriverField) specs.DriverField {
	f.custom = fmt.Sprintf("(%s)", fn)
	f.customArgs = args
	return f
}

// IsCustom returns true if the field is a custom function
func (f *field) IsCustom() bool {
	return f.custom != ""
}

// Custom returns the custom function
func (f *field) Custom() string {
	return f.custom
}

// CustomArgs returns the custom function arguments
func (f *field) CustomArgs() []specs.DriverField {
	return f.customArgs
}

// fnProcess processes the custom method for replacing arguments
func (f *field) fnProcess() (fn string, err error) {
	re := regexp.MustCompile(`\${([a-zA-Z_]+)}`)
	matches := re.FindAllStringSubmatch(f.Custom(), -1)

	replaceFieldCount := 0
	for _, match := range matches {
		for _, arg := range f.CustomArgs() {
			if arg.Name() != match[1] {
				continue
			}
			column, err := arg.Formatted()
			if err != nil {
				return "", err
			}

			fn = strings.Replace(f.Custom(), match[0], column, -1)
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

// Formatted returns the formatted field
func (f *field) Formatted() (string, error) {
	if f.Custom() != "" {
		return f.fnProcess()
	}
	return fmt.Sprintf("`t%d`.`%s`", f.Index(), f.Column()), nil
}

// NewField returns a new field
func NewField() specs.DriverField {
	return new(field)
}
