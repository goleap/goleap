package drivers

import (
	"github.com/lab210-dev/dbkit/specs"
	"strings"
)

type field struct {
	index        int
	name         string
	nameInSchema string
}

func (f *field) NameInSchema() string {
	return f.nameInSchema
}

func (f *field) SetNameInSchema(nameInSchema string) specs.DriverField {
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

func NewField() specs.DriverField {
	return new(field)
}
