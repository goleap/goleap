package specs

type Builder[T Model] interface {
	Get(primaryKey any) (T, error)
	Delete(primaryKey any) error

	Create() (err error)
	Update() error

	Find() (T, error)
	FindAll() ([]T, error)

	Fields(field ...string) Builder[T]
	Where(condition Condition) Builder[T]
	Limit(limit int) Builder[T]
	Offset(offset int) Builder[T]
	OrderBy(fields ...string) Builder[T]

	Count() (total int64, err error)

	Payload() PayloadAugmented[T]

	SetModel(model T) Builder[T]
}
