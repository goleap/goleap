package drivers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/kitstack/dbkit/specs"
	"github.com/sirupsen/logrus"
)

type transactionContext struct {
	identifier string
	success    chan bool
}

func newTransactionContext() *transactionContext {
	return &transactionContext{
		success: make(chan bool),
	}
}

func (t *transactionContext) SetIdentifier(identifier string) {
	t.identifier = identifier
}

func (t *transactionContext) Identifier() string {
	return t.identifier
}

func (t *transactionContext) Success() <-chan bool {
	return t.success
}

type connectionManager struct {
	*sql.DB
	transactions map[string]specs.Transaction
}

func WrapDB(db *sql.DB) specs.ConnectionManager {
	return &connectionManager{
		DB: db,
	}
}

func (d *connectionManager) GetConnection(ctx context.Context) (connection specs.Connection, err error) {
	return WrapConnection(ctx, d.DB)
}

func (d *connectionManager) GetTransaction(ctx context.Context) error {
	trxContext := ctx.Value(&transactionContext{}).(*transactionContext)
	if trxContext == nil {
		return errors.New("transaction context not found")
	}

	identifier := trxContext.identifier

	if d.transactions[identifier] == nil {
		connection, err := d.GetConnection(ctx)
		if err != nil {
			return err
		}

		tx, err := connection.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		d.transactions[identifier] = tx

		go func() {
			select {
			case <-ctx.Done():
				logrus.Debug("transaction context done")
				// just  kill the transaction !! auto commit = false
			case success := <-trxContext.Success():
				if success {
					logrus.Debug("transaction context success")
					// commit
				} else {
					logrus.Debug("transaction context failed")
					// rollback
				}
			}
		}()
	}

	return nil
}
