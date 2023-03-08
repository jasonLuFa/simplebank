# 💻 [golang migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

## [常用指令](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#usage)

- `create [-ext E] [-dir D] [-seq] [-digits N] [-format] NAME`
- `migrate up` : migrate up all pending migration
- `migrate up {N}` : migrate up N pending migration
- `migrate down` : migrate down all version
- `migrate down {N}` : migrate down N version base on the current version

## 範例

1. `mkdir -p db/migration`
2. `migrate create -ext sql -dir db/migration -seq init_schema`
   - 會在 db/migration 下自動生成兩個檔案 000001_init_schema.down.sql 和 000001_init_schema.sql
3. 使用 `migrate -path db/migration -database "postgresql://admin:admin@localhost:5433/simple_bank?sslmode=disable" -verbose up` 可執行 000001_init_schema.up.sql 此檔案裡的所有 sql 指令，`migrate -path db/migration -database "postgresql://admin:admin@localhost:5433/simple_bank?sslmode=disable" -verbose down` 則執行 000001_init_schema.down.sql 此檔案裡的所有 sql 指令 ( 這兩個指令也可寫在 Makefile 裡 )
