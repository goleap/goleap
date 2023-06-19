package fixtures

import (
	"context"
	"github.com/kitstack/dbkit"
	"github.com/kitstack/dbkit/connector/drivers/operators"
	"github.com/kitstack/dbkit/specs"
	"github.com/kitstack/dbkit/tests/models"
)

func (fixture *Fixture) BuilderFind(ctx context.Context) (err error) {

	user, err := dbkit.Use[*models.UsersModel](ctx).SetFields("Id").Find()

	fixture.Assert().NoError(err)
	fixture.Assert().EqualValues(1, user.Id)

	return
}

func (fixture *Fixture) BuilderFindNotFound(ctx context.Context) (err error) {

	user, err := dbkit.Use[*models.UsersModel](ctx).
		SetFields("Id").
		SetWhere(dbkit.NewCondition().SetFrom("Id").SetOperator(operators.Equal).SetTo(0)).
		Find()

	fixture.Assert().Error(err)
	fixture.Assert().Nil(user)
	fixture.Assert().Implements((*specs.ErrNotFound)(nil), err)
	fixture.Assert().Equal("empty result for UsersModel", err.Error())

	return nil
}
