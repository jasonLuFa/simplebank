# 🧪 Testing
## Gomock
- [document](https://github.com/golang/mock)
	- `go install github.com/golang/mock/mockgen@v1.6.0` ( Go 1.16+ )
	- `which mockgen` ( 確認安裝成功 )
- 優點
	- 獨立測試 ( 獨立測試資料和實際資料，避免衝突 )
	- 快速測試 ( 減少和 DB 溝通時間，所有的行為會存在於記憶體 )
	- 更高覆蓋率 ( 可以簡單的模擬 edge case )

- Mock DB 和 Real DB 需實作相同介面( interface )，gomock 根據 interface 去生成相對應的 mock code

- 兩種模式 :
	- Source mode
	- Reflect mode 
      ```sh
      # mockgen <import path to db query interface> <name of db query interface>
      # -package : 修改 packaga 名稱
      # -destination : 將產生出來的 code 寫入指定位置
      mockgen -package mockdb -destination db/mock/store.go github.com/jasonLuFa/simplebank/db/sqlc Store 
      ```
### 補充 : 
- 可在預測試的 interface 上掛上 `//go:generate <mockgen 指令>`，並在當前目錄下執行 `go generate ./...`，則就會自動建立當前目錄下所有 interface 的 mock code 
	- 範例 : ( 注意這邊使用 `go generate ./...`後，-destination 的位置如果是用相對位置的話，會是以當前此 interface 的檔案位置為主 )

    ```go
    //go:generate mockgen -package mockdb -destination ../mock/store.go github.com/jasonLuFa/simplebank/db/sqlc Store
    type Store interface {
      Querier
      TransferTx(context.Context, TransferTxParams) (TransferTxResult, error)
    }
    ```


### 範例
```Go
func TestGetAccountAPT(t *testing.T){
	account := randomAccount
	
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	
	// build stubs
	store.EXPECT().
		GETAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// start test server and send request
	server := NewServer(store)
	recorder := httptest.NewReocorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, nil := http.NewRequest(http.MethodGet, url ,nil)
	require.NoError(t, err)

	server.router.ServerHTTP(recorder, request)
	require.Equal(t, http.StatusOk, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)
}
```
