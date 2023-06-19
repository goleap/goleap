package drivers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/kitstack/dbkit/specs"
	"sync"
)

type TransactionContext struct {
	identifier string
	trx        *[]specs.Transaction
}

func NewTransactionContext() *TransactionContext {
	return &TransactionContext{
		trx: &[]specs.Transaction{},
	}
}

func (t *TransactionContext) SetIdentifier(identifier string) *TransactionContext {
	t.identifier = identifier
	return t
}

func (t *TransactionContext) Identifier() string {
	return t.identifier
}

func (t *TransactionContext) Add(trx specs.Transaction) {
	*t.trx = append(*t.trx, trx)
}

func (t *TransactionContext) Done(err error) {
	for _, trx := range *t.trx {
		if err != nil {
			_ = trx.Rollback()
			return
		}
		_ = trx.Commit()
	}
}

type connectionManager struct {
	sync.Mutex
	*sql.DB
	transactions map[string]specs.Transaction
}

func WrapDB(db *sql.DB) specs.ConnectionManager {
	return &connectionManager{
		DB:           db,
		transactions: make(map[string]specs.Transaction),
	}
}

func (d *connectionManager) GetConnection(ctx context.Context) (connection specs.Connection, err error) {
	return WrapConnection(ctx, d.DB)
}

func (d *connectionManager) GetTransaction(ctx context.Context) (transaction specs.Transaction, err error) {
	trxContextAny := ctx.Value(TransactionContext{})
	if trxContextAny == nil {
		return nil, errors.New("GetTransaction required a context with instance of transactionContext")
	}

	trxContext := trxContextAny.(*TransactionContext)

	if trxContext.Identifier() == "" {
		return nil, errors.New("GetTransaction required a context with instance of transactionContext with identifier")
	}

	d.Lock()
	defer d.Unlock()

	identifier := trxContext.identifier
	if d.transactions[identifier] == nil {
		connection, err := d.GetConnection(ctx)
		if err != nil {
			return nil, err
		}

		tx, err := connection.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}

		d.transactions[identifier] = tx

		trxContext.Add(tx)
	}

	return d.transactions[identifier], nil
}
