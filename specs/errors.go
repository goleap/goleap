package specs

type ErrNotFoundError error
type ErrPrimaryFieldNotFound error
type ErrFieldRequired error

type ErrUnknownOperator interface {
	error
	Operator() string
}

type ErrUnknownFields interface {
	error
	Fields() []string
}

type ErrRequiredFieldJoin interface {
	error
	Fields() []string
}

type ErrFieldNoFoundByColumn interface {
	error
	Column() string
	ModelDefinition() ModelDefinition
}
