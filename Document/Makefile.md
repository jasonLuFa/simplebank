# ⌨️ Makefile

- 使用 Makefile 來管理我們程式的構建，減少了大量輸入、拼寫錯誤，簡化構建項目的難度。真實線上環境配合 CI/CD 更佳
- 將一長串的的 command 縮短成自訂的 command
- ex : 建一個 Makefile 檔案

  ```
  postgres:
    docker run --name postgres12 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -p 5433:5432 -d postgres:12-alpine
  test:
    go test -v -cover ./...
  ```

  - 之後 command 只要下 `make postgres`，則就會執行 `docker run --name postgres12 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -p 5433:5432 -d postgres:12-alpine` 這一長串的指令了
