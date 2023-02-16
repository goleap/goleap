package driver

type Field interface {
	Index() int
	Name() string

	SetIndex(index int) Field
	SetName(name string) Field
}

type field struct {
	index int
	name  string
}

func (f *field) Index() int {
	return f.index
}

func (f *field) SetIndex(index int) Field {
	f.index = index
	return f
}

func (f *field) Name() string {
	return f.name
}

func (f *field) SetName(name string) Field {
	f.name = name
	return f
}

func NewField() Field {
	return new(field)
}
