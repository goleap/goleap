package dbkit

import (
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/depkit"
	"sync"
)

type subBuilder[T specs.Model] struct {
	sync.Mutex
	jobs map[string]specs.SubBuilderJob[T]
}

func newSubBuilder[T specs.Model]() specs.SubBuilder[T] {
	return &subBuilder[T]{
		jobs: make(map[string]specs.SubBuilderJob[T]),
	}
}

func (o *subBuilder[T]) AddJob(builder specs.Builder[T], fundamentalName string, model specs.ModelDefinition) specs.SubBuilder[T] {
	o.Lock()
	defer o.Unlock()

	if _, ok := o.jobs[fundamentalName]; ok {
		return o
	}

	depkit.Register[specs.NewSubBuilderJob[T]](newSubBuilderJob[T])

	o.jobs[fundamentalName] = depkit.Get[specs.NewSubBuilderJob[T]]()(builder, fundamentalName, model)

	return o

}

func (o *subBuilder[T]) Execute() (err error) {
	for _, job := range o.jobs {
		if err := job.Execute(); err != nil {
			return err
		}
	}
	return
}
