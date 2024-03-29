package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bacnx/simplebank/util"
	"github.com/stretchr/testify/require"
)

// When pass argument
//
// first arg is string to custom account.currency
func createRandomAccount(t *testing.T, arg ...interface{}) Account {
	user := createRandomUser(t)
	args := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	if len(arg) >= 1 {
		if currency, ok := arg[0].(string); ok {
			args.Currency = currency
		}
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.Owner, args.Owner)
	require.Equal(t, account.Balance, args.Balance)
	require.Equal(t, account.Balance, args.Balance)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account1.CreatedAt, time.Second)
}

func TestListAccounts(t *testing.T) {
	account := createRandomAccount(t)

	args := ListAccountsParams{
		Owner:  account.Owner,
		Limit:  3,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(accounts), 1)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	args := UpdateAccountParams{
		ID:      account.ID,
		Balance: 600,
	}

	account2, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)

	require.Equal(t, account.ID, account2.ID)
	require.Equal(t, account.Owner, account2.Owner)
	require.Equal(t, args.Balance, account2.Balance)
	require.Equal(t, account.Currency, account2.Currency)
	require.WithinDuration(t, account.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	account2, err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, account.ID, account2.ID)

	_, err = testQueries.GetAccount(context.Background(), account.ID)
	require.ErrorIs(t, err, sql.ErrNoRows)
}
