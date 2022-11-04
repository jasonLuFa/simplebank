package token

import (
	"testing"
	"time"

	"github.com/jasonLuFa/simplebank/util"
	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	issueAt := time.Now()
	expiredAt := issueAt.Add(duration)


	token, err:= maker.CreateToken(username, duration)
	require.NoError(t,err)
	require.NotEmpty(t,token)
	
	payload ,err:= maker.VerifyToken(token)
	require.NoError(t,err)
	require.NotEmpty(t,payload)

	require.NotZero(t,payload.ID)
	require.Equal(t,username,payload.Username)
	require.WithinDuration(t,issueAt,payload.IssuedAt, time.Second)
	require.WithinDuration(t,expiredAt,payload.ExpiredAt, time.Second)
}


func TestExpiredPasetoToken(t *testing.T){
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t,err)

	token, err := maker.CreateToken(util.RandomOwner(),-time.Minute)
	require.NoError(t,err)
	require.NotEmpty(t,token)

	payload, err := maker.VerifyToken(token)
	require.Error(t,err)
	require.EqualError(t,err,ErrExpiredToken.Error())
	require.Nil(t,payload)
}

func TestInvalidPasetoTokenAuth(t *testing.T){
	maker1, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	maker2, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	token, err:= maker1.CreateToken(username, duration)
	require.NoError(t,err)
	require.NotEmpty(t,token)
	
	payload ,err:= maker2.VerifyToken(token)
	require.Error(t,err)
	require.EqualError(t,err,paseto.ErrInvalidTokenAuth.Error())
	require.Nil(t,payload)
}