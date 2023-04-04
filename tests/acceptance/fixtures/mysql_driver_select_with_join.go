package fixtures

import (
	"context"
	"errors"
	"github.com/lab210-dev/dbkit"
	"github.com/lab210-dev/dbkit/connector/drivers"
	"github.com/lab210-dev/dbkit/specs"
	"github.com/lab210-dev/dbkit/tests/models"
)

func (f *Fixture) MysqlDriverSelectWithJoin(ctx context.Context) (err error) {
	joins := []specs.DriverJoin{
		drivers.NewJoin().
			SetToTable("posts").
			SetToTableIndex(1).
			SetToKey("id").
			SetFromKey("post_id").
			SetToDatabase("acceptance"),
	}
	fields := []specs.DriverField{
		drivers.NewField().SetName("id").SetNameInModel("Id"),
	}

	selectPayload := dbkit.NewPayload[*models.CommentsModel]()
	selectPayload.SetFields(fields)
	selectPayload.SetJoins(joins)

	err = f.Connector().Select(ctx, selectPayload)
	if err != nil {
		return err
	}

	if len(selectPayload.Result()) == 0 {
		return errors.New("result is empty")
	}

	if total := len(selectPayload.Result()); total != 8 {
		return errors.New("result is not equal to 7")
	}

	sum := 0
	for _, result := range selectPayload.Result() {
		sum += int(result.Id)
	}

	// Validate join return 7 rows with id 1, 2, 3, 4, 5, 6, 7, 8
	if sum != 36 {
		return errors.New("result is not equal to 36")
	}
	return
}
