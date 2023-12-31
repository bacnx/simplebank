package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/bacnx/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	randomCurrency := util.RandomCurrency()
	account1 := createRandomAccount(t, randomCurrency)
	account2 := createRandomAccount(t, randomCurrency)
	amount := int64(10)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

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

	existed := make(map[int]bool)

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

		// test accounts
		fromAccount, err := store.GetAccount(context.Background(), account1.ID)
		require.NoError(t, err)
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount, err := store.GetAccount(context.Background(), account2.ID)
		require.NoError(t, err)
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		// check account's balances
		fmt.Println(">> tx", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, 1 <= k && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	randomCurrency := util.RandomCurrency()
	account1 := createRandomAccount(t, randomCurrency)
	account2 := createRandomAccount(t, randomCurrency)
	amount := int64(10)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// run n concurrent transfer transactions
	n := 10

	errs := make(chan error, n)

	for i := 0; i < n; i++ {
		go func(i int) {
			fromAccountID := account1.ID
			toAccountID := account2.ID

			if i%2 == 1 {
				fromAccountID = account2.ID
				toAccountID = account1.ID
			}

			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}(i)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
