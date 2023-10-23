package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bacnx/simplebank/util"
	"github.com/stretchr/testify/assert"
)

func createRandomAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	assert.NoError(t, err)
	assert.NotEmpty(t, account)

	assert.Equal(t, account.Owner, args.Owner)
	assert.Equal(t, account.Balance, args.Balance)
	assert.Equal(t, account.Balance, args.Balance)

	assert.NotZero(t, account.ID)
	assert.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	assert.NoError(t, err)
	assert.NotEmpty(t, account2)

	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Owner, account2.Owner)
	assert.Equal(t, account1.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)
	assert.WithinDuration(t, account1.CreatedAt, account1.CreatedAt, time.Second)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	args := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), args)
	assert.NoError(t, err)
	assert.Len(t, accounts, 5)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	args := UpdateAccountParams{
		ID:      account.ID,
		Balance: 600,
	}

	account2, err := testQueries.UpdateAccount(context.Background(), args)
	assert.NoError(t, err)

	assert.Equal(t, account.ID, account2.ID)
	assert.Equal(t, account.Owner, account2.Owner)
	assert.Equal(t, args.Balance, account2.Balance)
	assert.Equal(t, account.Currency, account2.Currency)
	assert.WithinDuration(t, account.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	account2, err := testQueries.DeleteAccount(context.Background(), account.ID)
	assert.NoError(t, err)
	assert.Equal(t, account.ID, account2.ID)

	_, err = testQueries.GetAccount(context.Background(), account.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}
