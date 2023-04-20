package drivers

import (
	"fmt"
	"github.com/kitstack/dbkit/specs"
)

type limit struct {
	limit  int
	offset int
}

func (l *limit) Offset() int {
	return l.offset
}

func (l *limit) Limit() int {
	return l.limit
}

func (l *limit) SetOffset(index int) specs.DriverLimit {
	l.offset = index
	return l
}

func (l *limit) SetLimit(index int) specs.DriverLimit {
	l.limit = index
	return l
}

func (l *limit) Formatted() (string, error) {
	return fmt.Sprintf("LIMIT %d, %d", l.Offset(), l.Limit()), nil
}

func NewLimit() specs.DriverLimit {
	return new(limit)
}
