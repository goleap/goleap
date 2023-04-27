package specs

type NewSubBuilder[T Model] func() SubBuilder[T]

type SubBuilder[T Model] interface {
	AddJob(Builder[T], string, ModelDefinition) SubBuilder[T]
	Execute() error
}
