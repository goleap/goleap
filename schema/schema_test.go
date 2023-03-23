package schema

import (
	"github.com/lab210-dev/dbkit/connector/drivers"
	"github.com/lab210-dev/dbkit/specs"
	"github.com/lab210-dev/dbkit/tests/models"
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
	model := &models.BaseModel{}
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
	model := &models.BaseModel{}
	schema := New(model).Parse()

	id := schema.GetFieldByName("Id")
	test.Equal(id.Column(), "id")
	test.Equal(id.Index(), 0)

	test.Equal(id.Field(), drivers.NewField().SetIndex(0).SetName("id").SetNameInSchema("Id"))
}

func (test *SchemaTestSuite) TestGet() {
	model := &models.BaseModel{}
	schema := New(model).Parse()

	schema.Get()

	test.Equal(schema.Get(), model)
}

func (test *SchemaTestSuite) TestCopy() {
	model := &models.BaseModel{}
	schema := New(model).Parse()

	id := uint(1)
	schema.GetFieldByName("Id").Set(&id)

	copyOfSchema := schema.Copy()
	log.Print(copyOfSchema.GetFieldByName("Id").Get())
}

func (test *SchemaTestSuite) TestStruct() {
	schema := New(models.ExtraModel{}).Parse()
	test.Equal("extra", schema.TableName())
	test.Equal("test", schema.DatabaseName())
	test.Equal(268, len(schema.Fields()))
}

func (test *SchemaTestSuite) TestParseNilPtr() {
	var T *models.BaseModel
	New(T).Parse()
}

func (test *SchemaTestSuite) TestSet() {
	model := &models.BaseModel{}
	schema := New(model).Parse()
	// twice to test skip init
	schema.Parse()

	// Simple
	id1 := uint(1)
	schema.GetFieldByName("Id").Set(&id1)
	test.Equal(uint(1), model.Id)

	// Sub Embedded schema
	id2 := uint(2)
	schema.GetFieldByName("Recursive.Extra.Id").Set(&id2)
	test.Equal(uint(2), model.Recursive.Extra.Id)

	// Embedded schema
	id3 := uint(3)
	schema.GetFieldByName("Recursive.Id").Set(&id3)
	test.Equal(uint(3), model.Recursive.Id)
	test.Equal(uint(3), schema.GetFieldByName("Recursive.Id").Get())

	// Sub Embedded Slice Schema
	id4 := uint(4)
	schema.GetFieldByName("Recursive.Slice.Id").Set(&id4)
	test.Equal(uint(4), model.Recursive.Slice[0].Id)

	// With other instance of schema
	newSchema := New(model).Parse()

	// Two set on same field for testing skip init
	id5 := uint(5)
	newSchema.GetFieldByName("Recursive.Extra.Id").Set(&id5)
	newSchema.GetFieldByName("Recursive.Extra.Id").Set(&id5)
	test.Equal(uint(5), model.Recursive.Extra.Id)

	test.Equal(new(uint), newSchema.GetFieldByName("Recursive.Extra.Id").Copy())

	test.Equal(schema.GetFieldByName("Id"), schema.GetPrimaryKeyField())
}

func (test *SchemaTestSuite) TestJoin() {
	schemaTest := New(&models.BaseModel{}).Parse()
	test.Equal(schemaTest.GetFieldByName("Extra.Id").Join(), []specs.DriverJoin{
		drivers.NewJoin().
			SetFromTableIndex(72).
			SetToTable("base").
			SetToTableIndex(0).
			SetFromKey("extra_id").
			SetToKey("id"),
	})
}

func (test *SchemaTestSuite) TestGetPrimaryKeyField() {
	schemaTest := New(&models.ExtraModel{}).Parse()
	test.Equal(schemaTest.GetFieldByName("Id"), schemaTest.GetPrimaryKeyField())

	schemaTest = New(&models.ExtraJumpModel{}).Parse()
	test.Equal(nil, schemaTest.GetPrimaryKeyField())
}

func TestSchemaTestSuite(t *testing.T) {
	suite.Run(t, new(SchemaTestSuite))
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Call the function you want to benchmark here
		New(&models.BaseModel{}).Parse()
	}
}
