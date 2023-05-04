package specs

import "context"

type BuilderUse[T Model] func(ctx context.Context) Builder[T]

type Builder[T Model] interface {
	Context() context.Context
	Connector() (Connector, error)

	Get(primaryKey any) (T, error)
	Delete(primaryKey any) error

	Create() (err error)
	Update() error

	Find() (T, error)
	FindAll() ([]T, error)

	SetFields(field ...string) Builder[T]
	SetWhere(condition Condition) Builder[T]
	SetLimit(limit int) Builder[T]
	SetOffset(offset int) Builder[T]
	SetOrderBy(fields ...string) Builder[T]

	Count() (total int64, err error)

	Payload() PayloadAugmented[T]

	SetModel(model T) Builder[T]

	Fields() []string
	Wheres() []Condition
}
