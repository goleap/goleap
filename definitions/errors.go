package definitions

import (
	"fmt"
	"github.com/lab210-dev/dbkit/specs"
)

type ErrNotFoundError struct {
	field           string
	modelDefinition specs.ModelDefinition
}

func (e *ErrNotFoundError) Error() string {
	return fmt.Sprintf("field `%s` not found in model `%s`", e.field, e.modelDefinition.ModelValue().Type().Name())
}

func NewErrFieldNotFound(field string, modelDefinition specs.ModelDefinition) specs.ErrNotFoundError {
	return &ErrNotFoundError{
		field:           field,
		modelDefinition: modelDefinition,
	}
}

type ErrPrimaryFieldNotFound struct {
	modelDefinition specs.ModelDefinition
}

func (e *ErrPrimaryFieldNotFound) Error() string {
	return fmt.Sprintf("no primary field found in model `%s`", e.modelDefinition.ModelValue().Type().Name())
}

func NewErrNoPrimaryField(modelDefinition specs.ModelDefinition) specs.ErrPrimaryFieldNotFound {
	return &ErrPrimaryFieldNotFound{
		modelDefinition: modelDefinition,
	}
}

type ErrFieldNoFoundByColumn struct {
	column          string
	modelDefinition specs.ModelDefinition
}

func (e *ErrFieldNoFoundByColumn) Error() string {
	return fmt.Sprintf("field with column `%s` not found in model `%s`", e.column, e.modelDefinition.ModelValue().Type().Name())
}

func (e *ErrFieldNoFoundByColumn) Column() string {
	return e.column
}

func (e *ErrFieldNoFoundByColumn) ModelDefinition() specs.ModelDefinition {
	return e.modelDefinition
}

func NewErrFieldNoFoundByColumn(column string, modelDefinition specs.ModelDefinition) specs.ErrFieldNoFoundByColumn {
	return &ErrFieldNoFoundByColumn{
		column:          column,
		modelDefinition: modelDefinition,
	}
}
