package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/jasonLuFa/simplebank/db/mock"
	db "github.com/jasonLuFa/simplebank/db/sqlc"
	"github.com/jasonLuFa/simplebank/util"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok:= x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams,password string) gomock.Matcher{
	return eqCreateUserParamsMatcher{arg,password}
}


func TestCreateUserAPI(t *testing.T) {
  user, password := randomUser(t)
	testCases := []struct{
		name string
		body gin.H
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs:func (store *mockdb.MockStore)  {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusOK,recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().CreateUser(gomock.Any(),gomock.Any()).Times(1).Return(db.User{},sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusInternalServerError,recorder.Code)
			},
		},
		{
			name: "DuplicateUsernane",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore)  {
				// Error code 23505 -> ERROR:  duplicate key violates unique constraint
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{},&pq.Error{Code:"23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusForbidden,recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":  "invalid-user#1",
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Times(0) 是因為 func (server *Server) createUser(ctx *gin.Context) 中，會先在 ctx.ShouldBindJSON(&req) 中出錯( binding tag 的檢查 )，而不會執行到 server.store.CreateUser(ctx, arg)
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
				body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     "example#example.com",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "DuplicateEmail",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(),gomock.Any()).Times(1).Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username":  user.Username,
				"password":  "123",
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewTestServer(t,store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest("POST", url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomUser(t *testing.T) (db.User,string){
	password := util.RandomString(6)
		hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user := db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return user,password
}

func requireBodyMatchUser(t *testing.T,body *bytes.Buffer, user db.User){
	data,err := io.ReadAll(body)
	
	require.NoError(t,err)

	var gotUser db.User
	err = json.Unmarshal(data,&gotUser)

	require.NoError(t,err)
	require.Equal(t,gotUser.Username,user.Username)
	require.Equal(t,gotUser.FullName,user.FullName)
	require.Equal(t,gotUser.Email,user.Email)
	require.Empty(t,gotUser.HashedPassword)
}