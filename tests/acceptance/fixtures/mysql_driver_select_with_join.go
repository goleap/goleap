package fixtures

import (
	"context"
	"github.com/lab210-dev/dbkit"
	"github.com/lab210-dev/dbkit/connector/drivers"
	"github.com/lab210-dev/dbkit/connector/drivers/joins"
	"github.com/lab210-dev/dbkit/specs"
	"github.com/lab210-dev/dbkit/tests/models"
)

func (fixture *Fixture) MysqlDriverSelectWithJoin(ctx context.Context) (err error) {
	requiredJoins := []specs.DriverJoin{
		drivers.NewJoin().SetMethod(joins.Default).
			SetFrom(
				drivers.NewField().SetTable("comments").SetColumn("parent_id"),
			).
			SetTo(
				drivers.NewField().SetDatabase("acceptance").SetIndex(1).SetTable("comments").SetColumn("id"),
			),
	}
	fields := []specs.DriverField{
		drivers.NewField().SetIndex(1).SetColumn("id").SetName("Parent.Id"),
	}

	selectPayload := dbkit.NewPayload[*models.CommentsModel]()
	selectPayload.SetFields(fields)
	selectPayload.SetJoins(requiredJoins)

	err = fixture.Connector().Select(ctx, selectPayload)
	fixture.Assert().NoError(err)

	fixture.Assert().Len(selectPayload.Result(), 2)

	fixture.Assert().EqualValues(6, selectPayload.Result()[0].Parent.Id)
	fixture.Assert().EqualValues(7, selectPayload.Result()[1].Parent.Id)
	return
}

func (fixture *Fixture) MysqlDriverSelectWithJoinCustom(ctx context.Context) (err error) {
	// TODO: Change to real scenario because this is not a real scenario...
	requiredJoins := []specs.DriverJoin{
		drivers.NewJoin().SetMethod(joins.Default).
			SetFrom(
				drivers.NewField().SetTable("users").SetColumn("id"),
			).
			SetTo(
				drivers.NewField().SetDatabase("acceptance").SetIndex(1).SetTable("posts").SetCustom("SELECT MIN(id) FROM posts WHERE posts.user_id = %Id% LIMIT 1", []specs.DriverField{
					drivers.NewField().SetColumn("id").SetIndex(0).SetName("Id"),
				}),
			),
	}
	fields := []specs.DriverField{
		drivers.NewField().SetIndex(1).SetColumn("id").SetName("Id"),
	}

	selectPayload := dbkit.NewPayload[*models.UsersModel]()
	selectPayload.SetFields(fields)
	selectPayload.SetJoins(requiredJoins)

	err = fixture.Connector().Select(ctx, selectPayload)
	fixture.Assert().NoError(err)

	fixture.Assert().Len(selectPayload.Result(), 12)
	return
}
