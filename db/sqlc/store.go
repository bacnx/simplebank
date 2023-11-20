package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update account's balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		var account1, account2 Account
		if arg.FromAccountID < arg.ToAccountID {
			account1, account2, err = getTwoAccounts(ctx, q, arg.FromAccountID, arg.ToAccountID)
		} else {
			account2, account1, err = getTwoAccounts(ctx, q, arg.ToAccountID, arg.FromAccountID)
		}

		if err != nil {
			return err
		}

		if account1.Currency != account2.Currency {
			return errors.New("The currencies of the two accounts are not same")
		}

		// if account1.Balance < arg.Amount {
		// 	return errors.New("From account's balance is not enough")
		// }

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err =
				addMoney(ctx, q, arg.FromAccountID, account1.Balance-arg.Amount, arg.ToAccountID, account2.Balance+arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err =
				addMoney(ctx, q, arg.ToAccountID, account2.Balance+arg.Amount, arg.FromAccountID, account1.Balance-arg.Amount)
		}
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
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
