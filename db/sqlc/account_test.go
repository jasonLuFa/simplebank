package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jasonLuFa/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account{
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner: user.Username,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t,arg.Owner,account.Owner)
	require.Equal(t,arg.Balance,account.Balance)
	require.Equal(t,arg.Currency,account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T){
	account := createRandomAccount(t)
	accountGetting, err := testQueries.GetAccount(context.Background(),account.ID)
	require.NoError(t,err)
	require.NotEmpty(t,accountGetting)

	require.Equal(t,account.ID,accountGetting.ID)
	require.Equal(t,account.Owner,accountGetting.Owner)
	require.Equal(t,account.Balance,accountGetting.Balance)
	require.Equal(t,account.Currency,accountGetting.Currency)
	require.WithinDuration(t,account.CreatedAt,accountGetting.CreatedAt,time.Second)

}

func TestUpdateAccount(t *testing.T){
	account := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID: account.ID,
		Balance: util.RandomMoney(),
	}
	
	accountUpdatting, err := testQueries.UpdateAccount(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,accountUpdatting)

	require.Equal(t,account.ID,accountUpdatting.ID)
	require.Equal(t,account.Owner,accountUpdatting.Owner)
	require.Equal(t,arg.Balance,accountUpdatting.Balance)
	require.Equal(t,account.Currency,accountUpdatting.Currency)
}


func TestDeleteAccount(t *testing.T){
	account := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(),account.ID)
	require.NoError(t, err)

	accountGetting,err := testQueries.GetAccount(context.Background(),account.ID)
	// require.Error(t,err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountGetting)
}

func TestListAccounts(t *testing.T){
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner: lastAccount.Owner,
		Limit: 5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,accounts)

	for _, account := range accounts{
		require.NotEmpty(t, account)
		require.Equal(t,lastAccount.Owner, account.Owner)
	}
}