package specs

type NewSubBuilderJob[T Model] func(builder Builder[T], fundamentalName string, model ModelDefinition) SubBuilderJob[T]
type SubBuilderJob[T Model] interface {
	Execute() error
}
