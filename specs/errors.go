package specs

type FieldNotFoundError error
type PrimaryFieldNotFoundError error
type FieldRequiredError error

type UnknownOperatorErr interface {
	error
	Operator() string
}

type UnknownFieldsErr interface {
	error
	Fields() []string
}

type RequiredFieldJoinErr interface {
	error
	Fields() []string
}
