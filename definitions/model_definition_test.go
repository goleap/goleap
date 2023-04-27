package definitions

import (
	"errors"
	"github.com/kitstack/dbkit/connector/drivers"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/dbkit/tests/models"
	"github.com/stretchr/testify/suite"
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

	id, err := modelDefinition.GetFieldByName("Id")
	if !test.NoError(err) {
		return
	}
	test.Equal(id.Column(), "id")
	test.Equal(id.Index(), 0)

	test.Equal(id.Field(), drivers.NewField().SetIndex(0).SetColumn("id").SetName("Id"))
}

func (test *SchemaTestSuite) TestCopy() {
	model := &models.UsersModel{}
	modelDefinition := Use(model).Parse()

	test.Equal(modelDefinition.Copy(), model)

	id := uint(1)
	field, err := modelDefinition.GetFieldByName("Id")
	if !test.NoError(err) {
		return
	}
	field.Set(&id)

	snapshot := modelDefinition.Copy()

	id2 := uint(2)
	field, err = modelDefinition.GetFieldByName("Id")
	if !test.NoError(err) {
		return
	}
	field.Set(&id2)

	test.Equal(uint(1), snapshot.(*models.UsersModel).Id)
	test.Equal(uint(2), modelDefinition.Copy().(*models.UsersModel).Id)
}

func (test *SchemaTestSuite) TestSimpleFromSlice() {
	model := &models.PostsModel{}
	modelDefinition := Use(model).Parse()

	field, err := modelDefinition.GetFieldByName("Comments.Id")
	test.NoError(err)
	test.False(field.IsSlice())
	test.True(field.FromSlice())
}

func (test *SchemaTestSuite) TestComplexeFromSlice() {
	model := &models.CommentsModel{}
	modelDefinition := Use(model).Parse()

	field, err := modelDefinition.GetFieldByName("Post.Comments.Id")
	test.NoError(err)
	test.False(field.IsSlice())
	test.True(field.FromSlice())
}

func (test *SchemaTestSuite) TestComplexeModel() {
	modelDefinition := Use(&models.CommentsModel{}).Parse()
	test.Equal("comments", modelDefinition.TableName())
	test.Equal("acceptance", modelDefinition.DatabaseName())
	test.Equal(93, len(modelDefinition.Fields()))
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
	field, err := modelDefinition.GetFieldByName("Id")
	if !test.NoError(err) {
		return
	}
	field.Set(&id1)
	test.Equal(uint(1), model.Id)

	primary, err := modelDefinition.GetPrimaryField()
	if !test.NoError(err) {
		return
	}

	test.Equal(field, primary)

	// Sub Embedded modelDefinition
	id2 := uint(2)
	field, err = modelDefinition.GetFieldByName("Parent.User.Id")
	if !test.NoError(err) {
		return
	}
	field.Set(&id2)

	test.Equal(uint(2), model.Parent.User.Id)

	// Embedded modelDefinition
	id3 := uint(3)
	field, err = modelDefinition.GetFieldByName("Parent.Id")
	if !test.NoError(err) {
		return
	}
	field.Set(&id3)
	test.Equal(uint(3), model.Parent.Id)

	test.Equal(uint(3), field.Get())

	// Sub Embedded Slice ModelDefinition
	id4 := uint(4)
	field, err = modelDefinition.GetFieldByName("Post.Comments.Id")
	if !test.NoError(err) {
		return
	}
	field.Set(&id4)

	test.Equal(uint(4), model.Post.Comments[0].Id)

	newSchema := Use(model).Parse()

	// Two set on same fieldDefinition for testing skip init
	id5 := uint(5)
	field, err = newSchema.GetFieldByName("Parent.User.Id")
	if !test.NoError(err) {
		return
	}

	field.Set(&id5)
	field.Set(&id5)

	test.Equal(uint(5), model.Parent.User.Id)

	test.Equal(new(uint), field.Copy())
}

func (test *SchemaTestSuite) TestJoin() {
	schemaTest := Use(&models.CommentsModel{}).Parse()
	userIdFieldDefinition, err := schemaTest.GetFieldByName("User.Id")
	if !test.NoError(err) {
		return
	}

	postIdFieldDefinition, err := schemaTest.GetFieldByName("Post.Id")
	if !test.NoError(err) {
		return
	}

	test.Equal(userIdFieldDefinition.Join(), []specs.DriverJoin{
		drivers.NewJoin().
			SetFrom(drivers.NewField().SetIndex(0).SetColumn("user_id").SetTable("comments").SetDatabase("acceptance")).
			SetTo(drivers.NewField().SetIndex(1).SetColumn("id").SetTable("users").SetDatabase("acceptance")),
	})

	test.Equal(postIdFieldDefinition.Join(), []specs.DriverJoin{
		drivers.NewJoin().
			SetFrom(drivers.NewField().SetIndex(0).SetColumn("post_id").SetTable("comments").SetDatabase("acceptance")).
			SetTo(drivers.NewField().SetIndex(2).SetColumn("id").SetTable("posts").SetDatabase("acceptance")),
	})
}

func (test *SchemaTestSuite) TestGetFieldByName() {
	schemaTest := Use(&models.UsersModel{}).Parse()

	_, err := schemaTest.GetFieldByName("unknown")

	fieldErr := &ErrNotFoundError{}
	test.True(errors.As(err, &fieldErr))

	test.ErrorContains(err, "field `unknown` not found in model `UsersModel`")
}

func (test *SchemaTestSuite) TestGetPrimaryField() {
	schemaTest := Use(&models.UsersModel{}).Parse()

	idFieldDefinition, err := schemaTest.GetFieldByName("Id")
	if !test.NoError(err) {
		return
	}

	primaryFieldDefinition, err := schemaTest.GetPrimaryField()
	if !test.NoError(err) {
		return
	}

	test.Equal(idFieldDefinition, primaryFieldDefinition)

	schemaTest = Use(models.DebugModel{}).Parse()

	_, err = schemaTest.GetPrimaryField()

	primaryErr := &ErrPrimaryFieldNotFound{}
	test.True(errors.As(err, &primaryErr))

	test.ErrorContains(err, "no primary field found in model `DebugModel`")
}

func (test *SchemaTestSuite) TestGetToColumn() {
	schemaTest := Use(&models.CommentsModel{}).Parse()

	idFieldDefinition, err := schemaTest.GetFieldByName("Post.Comments.Content")
	if !test.NoError(err) {
		return
	}

	from, err := idFieldDefinition.Model().FromField().GetByColumn()
	if !test.NoError(err) {
		return
	}
	test.Equal("Post.Id", from.RecursiveFullName())

	to, err := idFieldDefinition.Model().FromField().GetToColumn()
	if !test.NoError(err) {
		return
	}
	test.Equal("Post.Comments.PostId", to.RecursiveFullName())

	// FundamentalName
	test.Equal("Post", idFieldDefinition.FundamentalName())

	_, err = idFieldDefinition.Model().GetFieldByColumn("unknown")
	test.Error(err)

	errFieldNoFoundByColumn := &ErrFieldNoFoundByColumn{}
	test.True(errors.As(err, &errFieldNoFoundByColumn))

	test.ErrorContains(err, "field with column `unknown` not found in model `CommentsModel`")
	test.Equal("unknown", errFieldNoFoundByColumn.Column())
	test.Equal("CommentsModel", errFieldNoFoundByColumn.ModelDefinition().TypeName())

}

func TestSchemaTestSuite(t *testing.T) {
	suite.Run(t, new(SchemaTestSuite))
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Use(&models.UsersModel{}).Parse()
	}
}
