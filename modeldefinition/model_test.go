package modeldefinition

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

func (test *SchemaTestSuite) TestValidateStruct() {
	model := &models.UsersModel{}
	modelDefinition := Use(model).Parse()

	test.Equal("acceptance", modelDefinition.DatabaseName())
	test.Equal("users", modelDefinition.TableName())

	if !test.Equal(7, len(modelDefinition.Fields())) {
		return
	}

	expectedFieldsName := []string{"Id", "Name", "Email", "Password", "CreatedAt", "UpdatedAt", "Validated", "BadTag"}
	for _, field := range modelDefinition.Fields() {
		test.Contains(expectedFieldsName, field.Name())
	}
}

func (test *SchemaTestSuite) TestFieldInfo() {
	model := &models.UsersModel{}
	modelDefinition := Use(model).Parse()

	id := modelDefinition.GetFieldByName("Id")
	test.Equal(id.Column(), "id")
	test.Equal(id.Index(), 0)

	test.Equal(id.Field(), drivers.NewField().SetIndex(0).SetName("id").SetNameInSchema("Id"))
}

func (test *SchemaTestSuite) TestGet() {
	model := &models.UsersModel{}
	modelDefinition := Use(model).Parse()

	modelDefinition.Get()

	test.Equal(modelDefinition.Get(), model)
}

func (test *SchemaTestSuite) TestCopy() {
	model := &models.UsersModel{}
	modelDefinition := Use(model).Parse()

	id := uint(1)
	modelDefinition.GetFieldByName("Id").Set(&id)

	copyOfSchema := modelDefinition.Copy()
	log.Print(copyOfSchema.GetFieldByName("Id").Get())
}

func (test *SchemaTestSuite) TestComplexeModel() {
	modelDefinition := Use(&models.CommentsModel{}).Parse()
	test.Equal("posts", modelDefinition.TableName())
	test.Equal("acceptance", modelDefinition.DatabaseName())
	test.Equal(69, len(modelDefinition.Fields()))
}

func (test *SchemaTestSuite) TestParseNilPtr() {
	var T *models.UsersModel
	modelDefinition := Use(T).Parse()
	test.Equal("users", modelDefinition.TableName())
	test.Equal("acceptance", modelDefinition.DatabaseName())
	test.Equal(7, len(modelDefinition.Fields()))
}

func (test *SchemaTestSuite) TestSet() {
	model := &models.CommentsModel{}
	modelDefinition := Use(model).Parse()
	// twice to test skip init
	modelDefinition.Parse()

	// Simple
	id1 := uint(1)
	modelDefinition.GetFieldByName("Id").Set(&id1)
	test.Equal(uint(1), model.Id)

	// Sub Embedded modelDefinition
	id2 := uint(2)
	modelDefinition.GetFieldByName("Parent.User.Id").Set(&id2)
	test.Equal(uint(2), model.Parent.User.Id)

	// Embedded modelDefinition
	id3 := uint(3)
	modelDefinition.GetFieldByName("Parent.Id").Set(&id3)
	test.Equal(uint(3), model.Parent.Id)
	test.Equal(uint(3), modelDefinition.GetFieldByName("Parent.Id").Get())

	// Sub Embedded Slice ModelDefinition
	id4 := uint(4)
	modelDefinition.GetFieldByName("Post.Comments.Id").Set(&id4)
	test.Equal(uint(4), model.Post.Comments[0].Id)

	newSchema := Use(model).Parse()

	// Two set on same field for testing skip init
	id5 := uint(5)
	newSchema.GetFieldByName("Parent.User.Id").Set(&id5)
	newSchema.GetFieldByName("Parent.User.Id").Set(&id5)
	test.Equal(uint(5), model.Parent.User.Id)

	test.Equal(new(uint), newSchema.GetFieldByName("Parent.User.Id").Copy())

	test.Equal(modelDefinition.GetFieldByName("Id"), modelDefinition.GetPrimaryKeyField())
}

func (test *SchemaTestSuite) TestJoin() {
	schemaTest := Use(&models.CommentsModel{}).Parse()
	test.Equal(schemaTest.GetFieldByName("User.Id").Join(), []specs.DriverJoin{
		drivers.NewJoin().
			SetFromTableIndex(1).
			SetToTable("posts").
			SetToTableIndex(0).
			SetFromKey("user_id").
			SetToKey("id"),
	})
	test.Equal(schemaTest.GetFieldByName("Post.Id").Join(), []specs.DriverJoin{
		drivers.NewJoin().
			SetFromTableIndex(2).
			SetToTable("posts").
			SetToTableIndex(0).
			SetFromKey("post_id").
			SetToKey("id"),
	})
}

func (test *SchemaTestSuite) TestGetPrimaryKeyField() {
	schemaTest := Use(&models.UsersModel{}).Parse()
	test.Equal(schemaTest.GetFieldByName("Id"), schemaTest.GetPrimaryKeyField())

	schemaTest = Use(models.ExtraModel{}).Parse()
	test.Equal(nil, schemaTest.GetPrimaryKeyField())
}

func TestSchemaTestSuite(t *testing.T) {
	suite.Run(t, new(SchemaTestSuite))
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Call the function you want to benchmark here
		Use(&models.UsersModel{}).Parse()
	}
}
