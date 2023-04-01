package dbkit

import "fmt"

type FieldRequiredError struct {
	queryType string
}

func (e *FieldRequiredError) Error() string {
	return fmt.Sprintf("field required for %s", e.queryType)
}

func NewFieldRequiredError(queryType string) *FieldRequiredError {
	return &FieldRequiredError{
		queryType: queryType,
	}
}

type NotFoundError struct {
	model string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("empty result for %s", e.model)
}

func NewNotFoundError(model string) *NotFoundError {
	return &NotFoundError{
		model: model,
	}
}
