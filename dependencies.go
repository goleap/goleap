package dbkit

import (
	"github.com/kitstack/dbkit/definitions"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/depkit"
	"github.com/kitstack/structkit"
	structKitSpecs "github.com/kitstack/structkit/specs"
	log "github.com/sirupsen/logrus"
)

func init() {
	injectDependencies()
}

// injectDependencies injects dependencies when package is imported
func injectDependencies() {
	depkit.Register[structKitSpecs.Get](structkit.Get)
	depkit.Register[structKitSpecs.Set](structkit.Set)
	depkit.Register[specs.UseModelDefinition](definitions.Use)
	depkit.Register[specs.BuilderUse[specs.Model]](Use[specs.Model])

	log.WithFields(log.Fields{
		"dependencies": []string{},
	}).Debug("Dependencies injected")
}

// injectGenericDependencies injects generic dependencies for all models on demand
func injectGenericDependencies[T specs.Model]() {
	depkit.Register[specs.NewSubBuilder[T]](newSubBuilder[T])

	log.WithFields(log.Fields{
		"dependencies": []string{},
	}).Debug("Generic dependencies injected")
}
