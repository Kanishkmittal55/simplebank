package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before: ", account1.Balance, account2.Balance)

	// Database transactions are easy to write but they must be handled with care
	// Because they can easily become a nightmare.
	// Best way to do transactions - Run it with several concurrent go routines
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {

		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// To Check K is unique
	existed := make(map[int]bool)

	// Check Results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// _____________________
		// check entries ( From and to )

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check Accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// Test : Check Account Balance
		fmt.Println(">> tx: ", fromAccount.Balance, toAccount.Balance)
		//        old balance       Post-transfer balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)     // Absolute value will be equal
		require.True(t, diff1 > 0)         //
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)  // should be unique for each tx.
		require.NotContains(t, existed, k) // Then we check the existed map should not contain k
		existed[k] = true                  // as that value of k occurred hence the corresponding boolean was set to true
	}

	// Check the final Updated Balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after: ", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before: ", account1.Balance, account2.Balance)

	// To send 5 transaction from account 1 to account 2
	// Then send 5 transactions in the reverse direction from account 2 to account 1
	n := 10
	amount := int64(10)

	errs := make(chan error)
	//results := make(chan TransferTxResult) --> Result is no longer Needed

	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID

		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: fromAccountId,
				ToAccountId:   toAccountId,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	// Check Only Errors
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	// Check the final Updated Balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after: ", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

}
