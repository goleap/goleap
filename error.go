package goleap

type FieldNotFoundError struct {
	message string
	field   string
}

func (e *FieldNotFoundError) Error() string {
	return e.message
}
