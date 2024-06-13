package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %w, rb err: %w", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func getTwoAccounts(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	accountID2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.GetAccountForUpdate(ctx, accountID1)
	if err != nil {
		return
	}
	account2, err = q.GetAccountForUpdate(ctx, accountID2)

	return
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	balance1 int64,
	accountID2 int64,
	balance2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      accountID1,
		Balance: balance1,
	})
	if err != nil {
		return
	}

	account2, err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      accountID2,
		Balance: balance2,
	})

	return
}
