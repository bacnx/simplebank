package db

import (
	"context"
	"errors"
)

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
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
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
			return errors.New("the currencies of the two accounts are not same")
		}

		// if account1.Balance < arg.Amount {
		// 	return errors.New("from account's balance is not enough")
		// }

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
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
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, account1.Balance-arg.Amount, arg.ToAccountID, account2.Balance+arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, account2.Balance+arg.Amount, arg.FromAccountID, account1.Balance-arg.Amount)
		}
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
