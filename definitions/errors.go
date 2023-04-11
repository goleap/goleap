package definitions

import (
	"fmt"
	"github.com/lab210-dev/dbkit/specs"
)

type FieldNotFoundError struct {
	field           string
	modelDefinition specs.ModelDefinition
}

func (e *FieldNotFoundError) Error() string {
	return fmt.Sprintf("field `%s` not found in model `%s`", e.field, e.modelDefinition.ModelValue().Type().Name())
}

func NewFieldNotFoundError(field string, modelDefinition specs.ModelDefinition) *FieldNotFoundError {
	return &FieldNotFoundError{
		field:           field,
		modelDefinition: modelDefinition,
	}
}

type ErrNoPrimaryField struct {
	modelDefinition specs.ModelDefinition
}

func (e *ErrNoPrimaryField) Error() string {
	return fmt.Sprintf("no primary field found in model `%s`", e.modelDefinition.ModelValue().Type().Name())
}

func NewErrNoPrimaryField(modelDefinition specs.ModelDefinition) *ErrNoPrimaryField {
	return &ErrNoPrimaryField{
		modelDefinition: modelDefinition,
	}
}
