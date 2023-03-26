package drivers

import (
	"fmt"
	"github.com/lab210-dev/dbkit/connector/drivers/joins"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoin(t *testing.T) {
	suite := assert.New(t)
	j := NewJoin()

	suite.Equal("", j.FromKey())
	j.SetFromKey("from_key")
	suite.Equal("from_key", j.FromKey())

	suite.Equal("", j.ToTable())
	j.SetToTable("to_table")
	suite.Equal("to_table", j.ToTable())

	suite.Equal("", j.ToKey())
	j.SetToKey("to_key")
	suite.Equal("to_key", j.ToKey())

	suite.Equal("", j.ToDatabase())
	j.SetToDatabase("to_database")
	suite.Equal("to_database", j.ToDatabase())

	suite.Equal(0, j.FromTableIndex())
	j.SetFromTableIndex(1)
	suite.Equal(1, j.FromTableIndex())

	suite.Equal(0, j.ToTableIndex())
	j.SetToTableIndex(2)
	suite.Equal(2, j.ToTableIndex())

	suite.Equal(joins.Method[joins.Default], j.Method())
	j.SetMethod(joins.Inner)
	suite.Equal(joins.Method[joins.Inner], j.Method())

	j.SetToKey("")

	err := j.Validate()
	suite.NotNil(err)
	suite.Equal(fmt.Errorf(`the following fields "ToKey" are mandatory to perform the join`), err)

	j.SetFromKey("from_key")
	j.SetToTable("to_table")
	j.SetToKey("to_key")

	err = j.Validate()
	suite.Nil(err)
}
