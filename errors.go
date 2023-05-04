package dbkit

import "fmt"

type ErrFieldRequired struct {
	queryType string
}

func (e *ErrFieldRequired) Error() string {
	return fmt.Sprintf("the method `%s` requires the selection of one or more fields", e.queryType)
}

func NewErrFieldRequired(queryType string) *ErrFieldRequired {
	return &ErrFieldRequired{
		queryType: queryType,
	}
}

type ErrNotFound struct {
	model string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("empty result for %s", e.model)
}

func NewErrNotFound(model string) *ErrNotFound {
	return &ErrNotFound{
		model: model,
	}
}
