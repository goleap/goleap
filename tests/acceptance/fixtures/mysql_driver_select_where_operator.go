package fixtures

import (
	"context"
	"github.com/lab210-dev/dbkit"
	"github.com/lab210-dev/dbkit/connector/drivers"
	"github.com/lab210-dev/dbkit/connector/drivers/operators"
	"github.com/lab210-dev/dbkit/specs"
	"github.com/lab210-dev/dbkit/tests/models"
)

func (fixture *Fixture) MysqlDriverSelectWhereEqual(ctx context.Context) (err error) {
	fields := []specs.DriverField{
		drivers.NewField().
			SetColumn("id").
			SetName("Id"),
	}
	wheres := []specs.DriverWhere{
		drivers.NewWhere().
			SetFrom(drivers.NewField().SetColumn("id")).
			SetOperator(operators.Equal).SetTo(1),
	}

	selectPayload := dbkit.NewPayload[*models.UsersModel]()
	selectPayload.SetFields(fields)
	selectPayload.SetWheres(wheres)

	err = fixture.Connector().Select(ctx, selectPayload)
	fixture.Assert().NoError(err)

	fixture.Assert().Len(selectPayload.Result(), 1)
	fixture.Assert().EqualValues(1, selectPayload.Result()[0].Id)

	return
}

func (fixture *Fixture) MysqlDriverSelectWhereNotEqual(ctx context.Context) (err error) {
	fields := []specs.DriverField{
		drivers.NewField().
			SetColumn("id").
			SetName("Id"),
	}
	wheres := []specs.DriverWhere{
		drivers.NewWhere().
			SetFrom(drivers.NewField().SetColumn("id")).
			SetOperator(operators.NotEqual).SetTo(1),
	}

	selectPayload := dbkit.NewPayload[*models.UsersModel]()
	selectPayload.SetFields(fields)
	selectPayload.SetWheres(wheres)

	err = fixture.Connector().Select(ctx, selectPayload)
	fixture.Assert().NoError(err)

	for _, user := range selectPayload.Result() {
		fixture.Assert().NotEqual(1, user.Id)
	}
	return
}

func (fixture *Fixture) MysqlDriverSelectWhereIn(ctx context.Context) (err error) {
	fields := []specs.DriverField{
		drivers.NewField().
			SetColumn("id").
			SetName("Id"),
	}
	wheres := []specs.DriverWhere{
		drivers.NewWhere().
			SetFrom(drivers.NewField().SetColumn("id")).
			SetOperator(operators.In).
			SetTo([]int{2, 999}),
	}

	payload := dbkit.NewPayload[*models.UsersModel]()
	payload.SetFields(fields)
	payload.SetWheres(wheres)

	err = fixture.Connector().Select(ctx, payload)
	fixture.Assert().NoError(err)
	fixture.Assert().Greater(len(payload.Result()), 0)

	for _, user := range payload.Result() {
		fixture.Assert().NotEqual(1, user.Id)
	}
	return
}

func (fixture *Fixture) MysqlDriverSelectWhereBetween(ctx context.Context) (err error) {
	fields := []specs.DriverField{
		drivers.NewField().
			SetColumn("id").
			SetName("Id"),
	}
	wheres := []specs.DriverWhere{
		drivers.NewWhere().
			SetFrom(drivers.NewField().SetColumn("created_at")).
			SetOperator(operators.Between).
			SetTo([]any{"2023-03-28 00:00:00", "2023-03-30 00:00:00"}),
	}

	payload := dbkit.NewPayload[*models.CommentsModel]()
	payload.SetFields(fields)
	payload.SetWheres(wheres)

	err = fixture.Connector().Select(ctx, payload)
	fixture.Assert().NoError(err)
	fixture.Assert().Greater(len(payload.Result()), 0)

	return
}

func (fixture *Fixture) MysqlDriverSelectWhereLike(ctx context.Context) (err error) {
	fields := []specs.DriverField{
		drivers.NewField().
			SetColumn("content").
			SetName("Content"),
	}
	wheres := []specs.DriverWhere{
		drivers.NewWhere().
			SetFrom(drivers.NewField().SetColumn("content")).
			SetOperator(operators.Like).
			SetTo("%pratique%"),
	}

	payload := dbkit.NewPayload[*models.CommentsModel]()
	payload.SetFields(fields)
	payload.SetWheres(wheres)

	err = fixture.Connector().Select(ctx, payload)
	fixture.Assert().NoError(err)
	fixture.Assert().Greater(len(payload.Result()), 0)

	for _, comment := range payload.Result() {
		fixture.Assert().NotEmpty(comment.Content)
		fixture.Assert().Contains(comment.Content, "pratique")
	}

	return
}

func (fixture *Fixture) MysqlDriverSelectWhereNotLike(ctx context.Context) (err error) {
	fields := []specs.DriverField{
		drivers.NewField().
			SetColumn("content").
			SetName("Content"),
	}
	wheres := []specs.DriverWhere{
		drivers.NewWhere().
			SetFrom(drivers.NewField().SetColumn("content")).
			SetOperator(operators.NotLike).
			SetTo("%pratique%"),
	}

	payload := dbkit.NewPayload[*models.CommentsModel]()
	payload.SetFields(fields)
	payload.SetWheres(wheres)

	err = fixture.Connector().Select(ctx, payload)
	fixture.Assert().NoError(err)
	fixture.Assert().Greater(len(payload.Result()), 0)

	for _, comment := range payload.Result() {
		fixture.Assert().NotEmpty(comment.Content)
		fixture.Assert().NotContains(comment.Content, "pratique")
	}

	return
}

func (fixture *Fixture) MysqlDriverSelectWhereGreaterAndLess(ctx context.Context) (err error) {
	fields := []specs.DriverField{
		drivers.NewField().
			SetColumn("id").
			SetName("Id"),
	}
	wheres := []specs.DriverWhere{
		drivers.NewWhere().
			SetFrom(drivers.NewField().SetColumn("id")).
			SetOperator(operators.Greater).
			SetTo(1),

		drivers.NewWhere().
			SetFrom(drivers.NewField().SetColumn("id")).
			SetOperator(operators.Less).
			SetTo(3),
	}

	payload := dbkit.NewPayload[*models.CommentsModel]()
	payload.SetFields(fields)
	payload.SetWheres(wheres)

	err = fixture.Connector().Select(ctx, payload)
	fixture.Assert().NoError(err)
	fixture.Assert().Len(payload.Result(), 1)

	fixture.Assert().Equal(payload.Result()[0].Id, uint(2))

	return
}

func (fixture *Fixture) MysqlDriverSelectWhereGreaterOrEqualAndLessOrEqual(ctx context.Context) (err error) {
	fields := []specs.DriverField{
		drivers.NewField().
			SetColumn("id").
			SetName("Id"),
	}
	wheres := []specs.DriverWhere{
		drivers.NewWhere().
			SetFrom(drivers.NewField().SetColumn("id")).
			SetOperator(operators.GreaterOrEqual).
			SetTo(1),
		drivers.NewWhere().
			SetFrom(drivers.NewField().SetColumn("id")).
			SetOperator(operators.LessOrEqual).
			SetTo(1),
	}

	payload := dbkit.NewPayload[*models.CommentsModel]()
	payload.SetFields(fields)
	payload.SetWheres(wheres)

	err = fixture.Connector().Select(ctx, payload)
	fixture.Assert().NoError(err)
	fixture.Assert().Len(payload.Result(), 1)

	fixture.Assert().Equal(payload.Result()[0].Id, uint(1))

	return
}

func (fixture *Fixture) MysqlDriverSelectWhereNotBetween(ctx context.Context) (err error) {
	fields := []specs.DriverField{
		drivers.NewField().
			SetColumn("id").
			SetName("Id"),
	}
	wheres := []specs.DriverWhere{
		drivers.NewWhere().
			SetFrom(drivers.NewField().SetColumn("created_at")).
			SetOperator(operators.NotBetween).
			SetTo([]any{"2023-03-28 00:00:00", "2023-03-30 00:00:00"}),
	}

	payload := dbkit.NewPayload[*models.CommentsModel]()
	payload.SetFields(fields)
	payload.SetWheres(wheres)

	err = fixture.Connector().Select(ctx, payload)
	fixture.Assert().NoError(err)
	fixture.Assert().Greater(len(payload.Result()), 0)

	return
}

func (fixture *Fixture) MysqlDriverSelectWhereGreaterWithFn(ctx context.Context) (err error) {
	fields := []specs.DriverField{
		drivers.NewField().
			SetColumn("id").
			SetName("Id"),
	}
	wheres := []specs.DriverWhere{
		drivers.NewWhere().
			SetTo(2).
			SetOperator(operators.Equal).
			SetFrom(drivers.NewField().SetCustom("SELECT COUNT(id) FROM posts WHERE posts.c_user_id = ${Id}", []specs.DriverField{
				drivers.NewField().SetColumn("id").SetIndex(0).SetName("Id"),
			})),
	}

	payload := dbkit.NewPayload[*models.UsersModel]()
	payload.SetFields(fields)
	payload.SetWheres(wheres)

	err = fixture.Connector().Select(ctx, payload)
	fixture.Assert().NoError(err)

	fixture.Assert().Len(payload.Result(), 1)
	fixture.Assert().EqualValues(1, payload.Result()[0].Id)

	return
}

func (fixture *Fixture) MysqlDriverSelectWhereTwoField(ctx context.Context) (err error) {
	fields := []specs.DriverField{
		drivers.NewField().
			SetColumn("id").
			SetName("Id"),
	}
	wheres := []specs.DriverWhere{
		drivers.NewWhere().
			SetFrom(drivers.NewField().SetColumn("id")).
			SetTo(drivers.NewField().SetColumn("parent_id")).
			SetOperator(operators.NotEqual),
	}

	payload := dbkit.NewPayload[*models.CommentsModel]()
	payload.SetFields(fields)
	payload.SetWheres(wheres)

	err = fixture.Connector().Select(ctx, payload)
	fixture.Assert().NoError(err)

	fixture.Assert().Len(payload.Result(), 2)

	fixture.Assert().EqualValues(7, payload.Result()[0].Id)
	fixture.Assert().EqualValues(8, payload.Result()[1].Id)
	return
}
