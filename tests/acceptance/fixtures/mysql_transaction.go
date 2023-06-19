package fixtures

import (
	"context"
	"github.com/kitstack/dbkit/connector/drivers"
)

func (fixture *Fixture) MysqlTransactionWithoutTransactionContext(ctx context.Context) (err error) {

	manager := fixture.Connector().Manager()
	tx, err := manager.GetTransaction(ctx)

	fixture.Assert().Nil(tx)
	fixture.Assert().NotNil(err)

	fixture.Assert().Equal("GetTransaction required a context with instance of transactionContext", err.Error())

	err = nil

	return
}

func (fixture *Fixture) MysqlTransactionWithTransactionContextWithoutIdentifier(ctx context.Context) (err error) {

	ctx = context.WithValue(ctx, drivers.TransactionContext{}, drivers.NewTransactionContext())

	manager := fixture.Connector().Manager()
	tx, err := manager.GetTransaction(ctx)

	fixture.Assert().Nil(tx)
	fixture.Assert().NotNil(err)

	fixture.Assert().Equal("GetTransaction required a context with instance of transactionContext with identifier", err.Error())

	err = nil

	return
}

func (fixture *Fixture) MysqlTransactionWithErr(ctx context.Context) (err error) {
	trxContext := drivers.NewTransactionContext().SetIdentifier("test")
	ctx = context.WithValue(ctx, drivers.TransactionContext{}, trxContext)

	defer trxContext.Done(err)

	manager := fixture.Connector().Manager()
	tx, err := manager.GetTransaction(ctx)

	fixture.Assert().NotNil(tx)

	_, err = tx.ExecContext(ctx, "DELETE FROM likes;")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM unknown;")
	err = nil

	return
}

func (fixture *Fixture) MysqlTransactionWithoutErr(ctx context.Context) (err error) {
	trxContext := drivers.NewTransactionContext().SetIdentifier("test")
	ctx = context.WithValue(ctx, drivers.TransactionContext{}, trxContext)

	defer trxContext.Done(err)

	manager := fixture.Connector().Manager()
	tx, err := manager.GetTransaction(ctx)

	fixture.Assert().NotNil(tx)

	_, err = tx.ExecContext(ctx, "INSERT INTO likes (user_id, post_id) VALUES (1, 2);")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM likes WHERE user_id = 1 AND post_id = 2;")
	err = nil

	cnx, err := manager.GetConnection(ctx)
	if err != nil {
		return err
	}
	defer cnx.Close()

	var found bool
	err = cnx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = 1 AND post_id = 2);").Scan(&found)

	fixture.Assert().Nil(err)
	fixture.Assert().False(found)

	return
}

func (fixture *Fixture) MysqlTransactionWithSameContextErr(ctx context.Context) (err error) {
	trxContext := drivers.NewTransactionContext().SetIdentifier("test")
	ctx = context.WithValue(ctx, drivers.TransactionContext{}, trxContext)

	defer func() {
		trxContext.Done(err)
		err = nil
	}()

	manager := fixture.Connector().Manager()
	tx, err := manager.GetTransaction(ctx)

	fixture.Assert().NotNil(tx)

	_, err = tx.ExecContext(ctx, "INSERT INTO likes (user_id, post_id) VALUES (1, 2);")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM likes WHERE user_id = 1 AND post_id = 2;")
	if err != nil {
		return err
	}

	managerExtend := fixture.Connector().Manager()
	txExtend, err := managerExtend.GetTransaction(ctx)

	fixture.Assert().NotNil(tx)

	_, err = txExtend.ExecContext(ctx, "DELETE FROM unknown;")

	cnx, err := manager.GetConnection(ctx)
	if err != nil {
		return err
	}
	defer cnx.Close()

	var found bool
	err = cnx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = 1 AND post_id = 2);").Scan(&found)

	fixture.Assert().Nil(err)
	fixture.Assert().False(found)

	return
}
