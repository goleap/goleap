package fixtures

import (
	"context"
	"errors"
	"github.com/lab210-dev/dbkit"
	"github.com/lab210-dev/dbkit/connector/drivers"
	"github.com/lab210-dev/dbkit/connector/drivers/operators"
	"github.com/lab210-dev/dbkit/specs"
	"github.com/lab210-dev/dbkit/tests/acceptance/fixtures/models"
)

func (f *Fixture) MysqlDriverSelectWhereEqual(ctx context.Context) (err error) {
	var joins []specs.DriverJoin
	fields := []specs.DriverField{
		drivers.NewField().
			SetName("id").
			SetNameInSchema("Id"),
	}
	wheres := []specs.DriverWhere{
		drivers.NewWhere().
			SetFrom(drivers.NewField().SetName("id")).
			SetOperator(operators.Equal).SetTo(1),
	}

	selectPayload := dbkit.NewPayload[*models.UserModel]()
	selectPayload.SetFields(fields)
	selectPayload.SetJoins(joins)
	selectPayload.SetWheres(wheres)

	err = f.Connector().Select(ctx, selectPayload)
	if err != nil {
		panic(err)
	}

	if selectPayload.Result()[0].Id != 1 {
		return errors.New("result is not equal to 1")
	}
	return
}

func (f *Fixture) MysqlDriverSelectWhereNotEqual(ctx context.Context) (err error) {
	var joins []specs.DriverJoin
	fields := []specs.DriverField{
		drivers.NewField().
			SetName("id").
			SetNameInSchema("Id"),
	}
	wheres := []specs.DriverWhere{
		drivers.NewWhere().
			SetFrom(drivers.NewField().SetName("id")).
			SetOperator(operators.NotEqual).SetTo(1),
	}

	selectPayload := dbkit.NewPayload[*models.UserModel]()
	selectPayload.SetFields(fields)
	selectPayload.SetJoins(joins)
	selectPayload.SetWheres(wheres)

	err = f.Connector().Select(ctx, selectPayload)
	if err != nil {
		panic(err)
	}

	if selectPayload.Result()[0].Id != 2 {
		return errors.New("result is not equal to 2")
	}
	return
}

func (f *Fixture) MysqlDriverSelectWhereIn(ctx context.Context) (err error) {
	var joins []specs.DriverJoin
	fields := []specs.DriverField{
		drivers.NewField().
			SetName("id").
			SetNameInSchema("Id"),
	}
	wheres := []specs.DriverWhere{
		drivers.NewWhere().
			SetFrom(drivers.NewField().SetName("id")).
			SetOperator(operators.In).SetTo([]int{2}),
	}

	selectPayload := dbkit.NewPayload[*models.UserModel]()
	selectPayload.SetFields(fields)
	selectPayload.SetJoins(joins)
	selectPayload.SetWheres(wheres)

	err = f.Connector().Select(ctx, selectPayload)
	if err != nil {
		panic(err)
	}

	if selectPayload.Result()[0].Id != 2 {
		return errors.New("result is not equal to 2")
	}
	return
}
