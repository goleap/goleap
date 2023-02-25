package schema

import (
	"github.com/lab210-dev/dbkit/connector/drivers"
	"github.com/lab210-dev/dbkit/specs"
	"github.com/lab210-dev/dbkit/testmodels"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

type SchemaTestSuite struct {
	suite.Suite
}

func (test *SchemaTestSuite) SetupTest() {

}

func (test *SchemaTestSuite) TestValidateStruct() {
	model := &testmodels.BaseModel{}
	schema := New(model).Parse()

	test.Equal("test", schema.DatabaseName())
	test.Equal("base", schema.TableName())

	if !test.Equal(90, len(schema.Fields())) {
		return
	}

	for _, field := range schema.FieldByName() {
		log.Print(field.RecursiveFullName())
	}
}

func (test *SchemaTestSuite) TestFieldInfo() {
	model := &testmodels.BaseModel{}
	schema := New(model).Parse()

	id := schema.GetFieldByName("Id")
	test.Equal(id.Column(), "id")
	test.Equal(id.Index(), 0)

	test.Equal(id.Field(), drivers.NewField().SetIndex(0).SetName("id"))
}

func (test *SchemaTestSuite) TestSet() {
	model := &testmodels.BaseModel{}
	schema := New(model).Parse()
	// twice to test skip init
	schema.Parse()

	// Simple
	schema.GetFieldByName("Id").Set(uint(1))
	test.Equal(uint(1), model.Id)

	// Sub Embedded schema
	schema.GetFieldByName("Recursive.Extra.Id").Set(uint(2))
	test.Equal(uint(2), model.Recursive.Extra.Id)

	// Embedded schema
	schema.GetFieldByName("Recursive.Id").Set(uint(3))
	test.Equal(uint(3), model.Recursive.Id)
	test.Equal(uint(3), schema.GetFieldByName("Recursive.Id").Get())

	// Sub Embedded Slice Schema
	schema.GetFieldByName("Recursive.Slice.Id").Set(uint(4))
	test.Equal(uint(4), model.Recursive.Slice[0].Id)

	// With other instance of schema
	newSchema := New(model).Parse()

	// Two set on same field for testing skip init
	newSchema.GetFieldByName("Recursive.Extra.Id").Set(uint(5))
	newSchema.GetFieldByName("Recursive.Extra.Id").Set(uint(5))
	test.Equal(uint(5), model.Recursive.Extra.Id)

	test.Equal(new(uint), newSchema.GetFieldByName("Recursive.Extra.Id").Copy())

	test.Equal(schema.GetFieldByName("Id"), schema.GetPrimaryKeyField())
}

func (test *SchemaTestSuite) TestJoin() {
	schemaTest := New(&testmodels.BaseModel{}).Parse()
	test.Equal(schemaTest.GetFieldByName("Extra.Id").Join(), []specs.DriverJoin{
		drivers.NewJoin().
			SetFromTable("extra").
			SetFromTableIndex(72).
			SetToTable("base").
			SetToTableIndex(0).
			SetFromKey("extra_id").
			SetToKey("id"),
	})
}

func (test *SchemaTestSuite) TestGetPrimaryKeyField() {
	schemaTest := New(&testmodels.ExtraModel{}).Parse()
	test.Equal(schemaTest.GetFieldByName("Id"), schemaTest.GetPrimaryKeyField())

	schemaTest = New(&testmodels.ExtraJumpModel{}).Parse()
	test.Equal(nil, schemaTest.GetPrimaryKeyField())
}

func TestSchemaTestSuite(t *testing.T) {
	suite.Run(t, new(SchemaTestSuite))
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Call the function you want to benchmark here
		New(&testmodels.BaseModel{}).Parse()
	}
}
