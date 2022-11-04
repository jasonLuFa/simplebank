package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jasonLuFa/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
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

func TestExpiredJWTToken(t *testing.T){
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t,err)

	token, err := maker.CreateToken(util.RandomOwner(),-time.Minute)
	require.NoError(t,err)
	require.NotEmpty(t,token)

	payload, err := maker.VerifyToken(token)
	require.Error(t,err)
	require.EqualError(t,err,ErrExpiredToken.Error())
	require.Nil(t,payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T){
	payload, err := NewPayload(util.RandomOwner(),time.Minute)
	require.NoError(t,err)
	require.NotEmpty(t,payload)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t,err)

	maker,err  := NewJWTMaker(util.RandomString(32))
	require.NoError(t,err)

	payload, err = maker.VerifyToken(token)
	require.Error(t,err)
	require.EqualError(t,err,ErrInvalidToken.Error())
	require.Nil(t, payload)
}