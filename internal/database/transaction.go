package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	raw *sqlx.Tx
}

func Begin(context context.Context) (*Transaction, error) {
	transaction := new(Transaction)

	var errBegin error
	transaction.raw, errBegin = raw.BeginTxx(context, nil)
	if errBegin != nil {
		return nil, errBegin
	}

	return transaction, nil
}

func (transaction *Transaction) Get(context context.Context, result any, query string, arguments ...any) error {
	return transaction.raw.GetContext(context, result, query, arguments...)
}

func (transaction *Transaction) Select(context context.Context, result any, query string, arguments ...any) error {
	return transaction.raw.SelectContext(context, result, query, arguments...)
}

func (transaction *Transaction) Execute(context context.Context, query string, arguments ...any) (*Result, error) {
	rawResult, errExecute := transaction.raw.ExecContext(context, query, arguments...)
	return newResult(rawResult), errExecute
}

func (transaction *Transaction) Commit() error {
	return transaction.raw.Commit()
}

func (transaction *Transaction) Rollback() error {
	return transaction.raw.Rollback()
}
