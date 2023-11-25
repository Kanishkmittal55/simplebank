package db

import (
	"context"
	"github.com/kanishkmittal55/simplebank/db/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomEntry(t *testing.T) Entry {
	accountID, err := testQueries.GetRandomAccountId(context.Background())
	require.NoError(t, err)
	arg := CreateEntryParams{
		AccountID: accountID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreatEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	arg := UpdateEntryParams{
		ID:     entry2.ID,          // ID remains the same , however you can change but logically should be unchangeable
		Amount: util.RandomMoney(), // The Entry Amount has been updated
	}

	updatedEntry, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedEntry)

	require.Equal(t, entry2.ID, updatedEntry.ID)
	require.Equal(t, entry2.AccountID, updatedEntry.AccountID)
	require.NotEqual(t, entry2.Amount, updatedEntry.Amount)
	require.WithinDuration(t, entry2.CreatedAt, updatedEntry.CreatedAt, time.Second)

}

// You need to implement these tests -

//func TestDeleteAccount(t *testing.T) {
//	account1 := createRandomAccount(t)
//	_, err := testQueries.DeleteAccount(context.Background(), account1.ID)
//	require.NoError(t, err)
//
//	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
//	require.Error(t, err)
//	require.Regexp(t, "no rows in result set", err.Error())
//	require.Empty(t, account2)
//}
//
//func TestListAccounts(t *testing.T) {
//	for i := 0; i < 10; i++ {
//		createRandomAccount(t)
//	}
//
//	arg := ListAccountsParams{
//		Limit:  5,
//		Offset: 5,
//	}
//
//	accounts, err := testQueries.ListAccounts(context.Background(), arg)
//	require.NoError(t, err)
//	require.Len(t, accounts, 5)
//
//	for _, account := range accounts {
//		require.NotEmpty(t, account)
//	}
//}
