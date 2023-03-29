package drivers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestField(t *testing.T) {
	suite := assert.New(t)
	f := NewField()

	suite.Equal("", f.NameInSchema())
	f.SetNameInSchema("field_name")
	suite.Equal("field_name", f.NameInSchema())

	suite.Equal(0, f.Index())
	f.SetIndex(5)
	suite.Equal(5, f.Index())

	suite.Equal("", f.Name())
	f.SetName(" field name ")
	suite.Equal("field name", f.Name())
}
