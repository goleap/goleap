package fixtures

import (
	"context"
	"github.com/lab210-dev/dbkit"
	"github.com/lab210-dev/dbkit/tests/models"
)

func (fixture *Fixture) BuilderGet(ctx context.Context) (err error) {

	user, err := dbkit.Use[*models.UsersModel](ctx, fixture.Connector()).Fields("Id").Get(1)

	fixture.Assert().NoError(err)
	fixture.Assert().EqualValues(1, user.Id)

	return
}

func (fixture *Fixture) BuilderGetWithJoin(ctx context.Context) (err error) {

	post, err := dbkit.Use[*models.PostsModel](ctx, fixture.Connector()).Fields("Id", "Creator.Id", "Creator.Validated", "Editor.Id", "Comments.Id").Get(1)

	fixture.Assert().NoError(err)
	fixture.Assert().EqualValues(1, post.Id)
	fixture.Assert().EqualValues(1, post.Creator.Id)
	fixture.Assert().EqualValues(true, post.Creator.Validated)

	return
}
