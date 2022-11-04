package db

import (
	"context"
	"testing"
	"time"

	"github.com/jasonLuFa/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t,err)
	
	arg := CreateUserParams{
		Username:    util.RandomOwner(),
		HashedPassword:  hashedPassword,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName,user.FullName)
	require.Equal(t, arg.Email,user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	userCreating := createRandomUser(t)
	userGetting, err := testQueries.GetUser(context.Background(), userCreating.Username)
	require.NoError(t, err)
	require.NotEmpty(t, userGetting)

	require.Equal(t, userCreating.Username, userGetting.Username)
	require.Equal(t, userCreating.HashedPassword, userGetting.HashedPassword)
	require.Equal(t, userCreating.FullName, userGetting.FullName)
	require.Equal(t, userCreating.Email, userGetting.Email)
	require.WithinDuration(t, userCreating.PasswordChangedAt, userGetting.PasswordChangedAt, time.Second)
	require.WithinDuration(t, userCreating.CreatedAt, userGetting.CreatedAt, time.Second)

}