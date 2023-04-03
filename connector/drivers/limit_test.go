package drivers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitOffset(t *testing.T) {
	l := &limit{limit: 10, offset: 5}

	assert.Equal(t, 5, l.Offset())
}

func TestLimitLimit(t *testing.T) {
	l := &limit{limit: 10, offset: 5}

	assert.Equal(t, 10, l.Limit())
}

func TestLimitSetOffset(t *testing.T) {
	l := &limit{limit: 10, offset: 5}

	l.SetOffset(15)

	assert.Equal(t, 15, l.Offset())
}

func TestLimitSetLimit(t *testing.T) {
	l := &limit{limit: 10, offset: 5}

	l.SetLimit(20)

	assert.Equal(t, 20, l.Limit())
}

func TestNewLimit(t *testing.T) {
	l := NewLimit()

	assert.NotNil(t, l)
	assert.IsType(t, &limit{}, l)
	assert.Zero(t, l.Limit())
	assert.Zero(t, l.Offset())
}
