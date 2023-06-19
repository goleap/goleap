package fixtures

import (
	"context"
	"time"
)

func (fixture *Fixture) MysqlConnection(ctx context.Context) (err error) {

	manager := fixture.Connector().Manager()
	cnx, err := manager.GetConnection(ctx)

	if err != nil {
		return err
	}

	defer cnx.Close()

	fixture.Assert().NotNil(cnx)

	err = cnx.PingContext(ctx)
	if err != nil {
		return err
	}

	return
}

func (fixture *Fixture) MysqlConnectionSlow(ctx context.Context) (err error) {

	manager := fixture.Connector().Manager()
	cnx, err := manager.GetConnection(ctx)

	if err != nil {
		return err
	}

	defer cnx.Close()

	fixture.Assert().NotNil(cnx)

	_, err = cnx.ExecContext(ctx, "SELECT SLEEP(1)")
	if err != nil {
		return err
	}

	return
}

func (fixture *Fixture) MysqlConnectionSlowWhenCancellingContext(ctx context.Context) (err error) {

	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	manager := fixture.Connector().Manager()
	cnx, err := manager.GetConnection(ctx)

	if err != nil {
		return err
	}

	defer cnx.Close()

	testingCnx, err := manager.GetConnection(ctx)

	if err != nil {
		return err
	}

	defer testingCnx.Close()

	fixture.Assert().NotNil(cnx)

	_, err = cnx.ExecContext(ctx, "SELECT SLEEP(1)")
	fixture.Assert().NotNil(err)
	fixture.Assert().EqualValues("context deadline exceeded", err.Error())

	var query *string
	err = testingCnx.QueryRowContext(context.Background(), "SELECT GROUP_CONCAT(info SEPARATOR ', ') AS info FROM information_schema.processlist WHERE Id = ?", cnx.Id()).Scan(&query)
	if err != nil {
		return err
	}

	fixture.Assert().Nil(query)

	err = nil

	return
}
