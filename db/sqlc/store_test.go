package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/jasonLuFa/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
  store := NewStore(db)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// run n concurrent transfer transactions
	n := 5
	amount := "10.0"
	
	errs := make(chan error)
	results := make(chan TransferTxResult)
	
	for i := 0; i < n; i++ {
		go func ()  {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})

			errs <- err
			fmt.Println(">>>> Transfer :",result)
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t,result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t,transfer)
		require.Equal(t,account1.ID,transfer.FromAccountID)
		require.Equal(t,account2.ID,transfer.ToAccountID)
		require.Equal(t,amount,transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		fmt.Println("transfer id:",transfer.ID)
		_, err = store.GetTransfer(context.Background(),transfer.ID)
		require.NoError(t,err)

		// check fromEntry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID,fromEntry.AccountID)

		require.Equal(t, fmt.Sprintf("%.2f",-util.StringToFloat64(amount)),fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(),fromEntry.ID)
		require.NoError(t,err)

		// check toEntry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID,toEntry.AccountID)
		require.Equal(t, fmt.Sprintf("%.2f",util.StringToFloat64(amount)),toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(),toEntry.ID)
		require.NoError(t,err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t,fromAccount)
		require.Equal(t, account1.ID,fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t,toAccount)
		require.Equal(t, account2.ID,toAccount.ID)
		
		// check accounts' balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		fromAccountBalanceDiff := util.StringToFloat64(account1.Balance) - util.StringToFloat64(fromAccount.Balance)
		toAccountBalanceDiff := util.StringToFloat64(toAccount.Balance) - util.StringToFloat64(account2.Balance)
		require.Equal(t, fromAccountBalanceDiff, toAccountBalanceDiff)
		require.True(t,fromAccountBalanceDiff > 0)
		txCount := int(fromAccountBalanceDiff / util.StringToFloat64(amount))
		require.True(t,txCount >= 1 && txCount <= n )

	}

		// check the final updated balances
		updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
		require.NoError(t, err)

		updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
		require.NoError(t, err)

		fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
		
		require.Equal(t, util.StringToFloat64(account1.Balance) - float64(n)*util.StringToFloat64(amount), util.StringToFloat64(updatedAccount1.Balance))
		require.Equal(t, util.StringToFloat64(account2.Balance) + float64(n)*util.StringToFloat64(amount), util.StringToFloat64(updatedAccount2.Balance))
}