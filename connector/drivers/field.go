package drivers

import "github.com/lab210-dev/dbkit/specs"

type field struct {
	index int
	name  string
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
	f.name = name
	return f
}

func NewField() specs.DriverField {
	return new(field)
}
