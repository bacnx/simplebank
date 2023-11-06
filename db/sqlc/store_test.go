package db

import (
	"context"
	"testing"

	"github.com/bacnx/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	randomCurrency := util.RandomCurrency()
	account1 := createRandomAccount(t, randomCurrency)
	account2 := createRandomAccount(t, randomCurrency)
	amount := int64(100)

	// run n concurrent transfer transactions
	n := 5

	results := make(chan TransferTxResult, n)
	errs := make(chan error, n)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			results <- result
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.NotEmpty(t, transfer.ID)
		require.NotEmpty(t, transfer.CreatedAt)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		entry1 := result.FromEntry
		require.NotEmpty(t, entry1)
		require.NotEmpty(t, entry1.ID)
		require.Equal(t, entry1.AccountID, account1.ID)
		require.Equal(t, entry1.Amount, -amount)

		_, err = store.GetEntry(context.Background(), entry1.ID)
		require.NoError(t, err)

		entry2 := result.ToEntry
		require.NotEmpty(t, entry2)
		require.NotEmpty(t, entry2.ID)
		require.NotEmpty(t, entry2.CreatedAt)
		require.Equal(t, entry2.AccountID, account2.ID)
		require.Equal(t, entry2.Amount, amount)

		_, err = store.GetEntry(context.Background(), entry2.ID)
		require.NoError(t, err)

		// TODO: test account's balance
	}
}
