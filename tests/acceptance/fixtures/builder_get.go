package fixtures

import (
	"context"
	"github.com/kitstack/dbkit"
	"github.com/kitstack/dbkit/connector/drivers/operators"
	"github.com/kitstack/dbkit/connectors"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/dbkit/tests/models"
)

func (fixture *Fixture) BuilderUnknownConnector(ctx context.Context) (err error) {
	connectors.Instance().Remove("acceptance")

	_, err = dbkit.Use[*models.UsersModel](ctx).SetFields("Id").Get(1)

	fixture.Assert().Error(err)
	fixture.Assert().Equal("unknown connector: acceptance", err.Error())
	fixture.Assert().Implements((*specs.ErrConnectorNotFound)(nil), err)

	return nil
}

func (fixture *Fixture) BuilderGet(ctx context.Context) (err error) {

	user, err := dbkit.Use[*models.UsersModel](ctx).SetFields("Id").Get(1)

	fixture.Assert().NoError(err)
	fixture.Assert().EqualValues(1, user.Id)

	return
}

func (fixture *Fixture) BuilderGetNotFound(ctx context.Context) (err error) {

	user, err := dbkit.Use[*models.UsersModel](ctx).SetFields("Id").Get(0)

	fixture.Assert().Error(err)
	fixture.Assert().Nil(user)
	fixture.Assert().Equal("empty result for UsersModel", err.Error())

	return nil
}

func (fixture *Fixture) BuilderGetWithMany2Many(ctx context.Context) (err error) {

	comment, err := dbkit.Use[*models.CommentsModel](ctx).SetFields("Id", "Post.Id", "Post.Comments.Id", "Post.Comments.Content").Get(1)

	fixture.Assert().NoError(err)
	fixture.Assert().EqualValues(1, comment.Id)
	fixture.Assert().EqualValues(1, comment.Post.Id)
	fixture.Assert().Len(comment.Post.Comments, 2)
	fixture.Assert().EqualValues(1, comment.Post.Comments[0].Id)
	fixture.Assert().EqualValues(2, comment.Post.Comments[1].Id)

	return
}

func (fixture *Fixture) BuilderGetWithMany2ManyFilter(ctx context.Context) (err error) {

	comment, err := dbkit.Use[*models.CommentsModel](ctx).
		SetFields("Id", "Post.Id", "Post.Comments.Id", "Post.Comments.Content", "Post.Comments.User.Id").
		SetWhere(dbkit.NewCondition().SetFrom("Post.Comments.User.Id").SetOperator(operators.Equal).SetTo(2)).
		Get(1)

	fixture.Assert().NoError(err)
	fixture.Assert().EqualValues(1, comment.Id)
	fixture.Assert().EqualValues(1, comment.Post.Id)

	fixture.Assert().Len(comment.Post.Comments, 1)

	fixture.Assert().EqualValues(1, comment.Post.Comments[0].Id)
	fixture.Assert().EqualValues(2, comment.Post.Comments[0].User.Id)

	return
}

func (fixture *Fixture) BuilderGetWithMany2One(ctx context.Context) (err error) {

	comment, err := dbkit.Use[*models.CommentsModel](ctx).SetFields("Id", "Post.Id").Get(1)

	fixture.Assert().NoError(err)
	fixture.Assert().EqualValues(1, comment.Id)
	fixture.Assert().EqualValues(1, comment.Post.Id)

	return
}
