package db

import (
	"context"
	"testing"

	"github.com/jasonLuFa/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fromAccount, toAccount Account) Transfer {
		arg := CreateTransferParams{
			FromAccountID: fromAccount.ID,
			ToAccountID: toAccount.ID,
			Amount: util.RandomMoney(),
		}

		transfer , err := testQueries.CreateTransfer(context.Background(),arg)
		require.NoError(t, err)
		require.NotEmpty(t, transfer)

		require.Equal(t, transfer.FromAccountID,arg.FromAccountID)
		require.Equal(t, transfer.ToAccountID,arg.ToAccountID)
		require.Equal(t, transfer.Amount, arg.Amount)

		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		return transfer
}


func TestCreateTransfer(t *testing.T){
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	createRandomTransfer(t, fromAccount, toAccount)
}

func TestGetTransfer(t *testing.T){
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	transferCreating  := createRandomTransfer(t, fromAccount, toAccount)

	transferGetting ,err := testQueries.GetTransfer(context.Background(),transferCreating.ID)
	require.NoError(t,err)
	require.NotEmpty(t,transferGetting)

	require.Equal(t, transferCreating.ID, transferGetting.ID)
	require.Equal(t, transferCreating.FromAccountID, transferGetting.FromAccountID)
	require.Equal(t, transferCreating.ToAccountID, transferGetting.ToAccountID)
	require.Equal(t, transferCreating.Amount, transferGetting.Amount)
}

func TestListTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, account1, account2)
		createRandomTransfer(t, account2, account1)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}
}