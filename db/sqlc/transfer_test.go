package db

import (
	"context"
	"testing"

	"github.com/bacnx/simplebank/util"
	"github.com/stretchr/testify/assert"
)

func createRandomTransfer(t *testing.T, fromAccoutID int64, toAccountID int64) Transfer {
	args := CreateTransferParams{
		FromAccountID: fromAccoutID,
		ToAccountID:   toAccountID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	assert.NoError(t, err)

	assert.NotZero(t, transfer.ID)
	assert.Equal(t, transfer.FromAccountID, args.FromAccountID)
	assert.Equal(t, transfer.ToAccountID, args.ToAccountID)
	assert.Equal(t, transfer.Amount, args.Amount)
	assert.NotZero(t, transfer.CreatedAt)

	return transfer
}

func createRandomTransfers(t *testing.T, numberOfTransfers int) ([]Transfer, int64, int64) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	var transfers []Transfer

	for i := 0; i < numberOfTransfers; i++ {
		transfer := createRandomTransfer(t, fromAccount.ID, toAccount.ID)
		transfers = append(transfers, transfer)
	}

	return transfers, fromAccount.ID, toAccount.ID
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfers(t, 1)
}

func TestGetTransfer(t *testing.T) {
	transfers, _, _ := createRandomTransfers(t, 1)
	transfer1 := transfers[0]

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	assert.NoError(t, err)

	assert.Equal(t, transfer1.ID, transfer2.ID)
	assert.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	assert.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	assert.Equal(t, transfer1.Amount, transfer2.Amount)
	assert.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
}

func TestListTransfers(t *testing.T) {
	_, fromAccountID, toAccountID := createRandomTransfers(t, 10)
	args := ListTransfersParams{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), args)
	assert.NoError(t, err)

	assert.Len(t, transfers, 5)
	assert.Equal(t, transfers[0].FromAccountID, fromAccountID)
	assert.Equal(t, transfers[0].ToAccountID, toAccountID)
}

func TestListTransfersByFromAccountId(t *testing.T) {
	_, fromAccountID, _ := createRandomTransfers(t, 10)
	args := ListTransfersByFromAccountIdParams{
		FromAccountID: fromAccountID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfersByFromAccountId(context.Background(), args)
	assert.NoError(t, err)

	assert.Len(t, transfers, 5)
	assert.Equal(t, transfers[0].FromAccountID, fromAccountID)
}

func TestListTransfersByToAccountId(t *testing.T) {
	_, _, toAccountID := createRandomTransfers(t, 10)
	args := ListTransfersByToAccountIdParams{
		ToAccountID: toAccountID,
		Limit:       5,
		Offset:      5,
	}

	transfers, err := testQueries.ListTransfersByToAccountId(context.Background(), args)
	assert.NoError(t, err)

	assert.Len(t, transfers, 5)
	assert.Equal(t, transfers[0].ToAccountID, toAccountID)
}
