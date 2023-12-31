package db

import (
	"context"
	"testing"

	"github.com/bacnx/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, accountID int64) Entry {
	args := CreateEntryParams{
		AccountID: accountID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)

	require.NotZero(t, entry.ID)
	require.Equal(t, entry.AccountID, accountID)

	return entry
}

func createRandomEntries(t *testing.T, numberOfEntries int) ([]Entry, int64) {
	account := createRandomAccount(t)
	var entries []Entry

	for i := 0; i < numberOfEntries; i++ {
		entry := createRandomEntry(t, account.ID)
		entries = append(entries, entry)
	}

	return entries, account.ID
}

func TestCreateEntry(t *testing.T) {
	createRandomEntries(t, 1)
}

func TestGetEntry(t *testing.T) {
	entries, _ := createRandomEntries(t, 1)
	entry1 := entries[0]
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
}

func TestListEntries(t *testing.T) {
	_, accountID := createRandomEntries(t, 10)
	args := ListEntriesParams{
		AccountID: accountID,
		Limit:     5,
		Offset:    5,
	}
	entries, err := testQueries.ListEntries(context.Background(), args)
	require.NoError(t, err)

	require.Len(t, entries, 5)
	require.Equal(t, entries[0].AccountID, accountID)
}
