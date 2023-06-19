package fixtures

import (
	"context"
	"github.com/kitstack/dbkit"
	"github.com/kitstack/dbkit/tests/models"
)

func (fixture *Fixture) BuilderFindAll(ctx context.Context) (err error) {

	users, err := dbkit.Use[*models.UsersModel](ctx).SetFields("Id").FindAll()

	fixture.Assert().NoError(err)
	fixture.Assert().Len(users, 3)

	fixture.Assert().EqualValues(1, users[0].Id)
	fixture.Assert().EqualValues(2, users[1].Id)
	fixture.Assert().EqualValues(3, users[2].Id)

	return
}
