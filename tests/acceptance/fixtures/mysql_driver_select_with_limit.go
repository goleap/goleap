package fixtures

import (
	"context"
	"github.com/lab210-dev/dbkit"
	"github.com/lab210-dev/dbkit/connector/drivers"
	"github.com/lab210-dev/dbkit/specs"
	"github.com/lab210-dev/dbkit/tests/models"
)

func (fixture *Fixture) MysqlDriverSelectWithLimit(ctx context.Context) (err error) {
	fields := []specs.DriverField{
		drivers.NewField().
			SetName("Id").
			SetColumn("id").
			SetIndex(0),
	}

	selectPayload := dbkit.NewPayload[*models.UsersModel]()
	selectPayload.SetFields(fields)
	selectPayload.SetLimit(drivers.NewLimit().SetLimit(1))

	err = fixture.Connector().Select(ctx, selectPayload)
	if err != nil {
		return err
	}

	fixture.Assert().Len(selectPayload.Result(), 1)
	fixture.Assert().EqualValues(1, selectPayload.Result()[0].Id)

	return
}
