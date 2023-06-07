# 💻 [golang migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

## [常用指令](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#usage)

- `create [-ext E] [-dir D] [-seq] [-digits N] [-format] NAME`
- `migrate up` : migrate up all pending migration
- `migrate up {N}` : migrate up N pending migration
- `migrate down` : migrate down all version
- `migrate down {N}` : migrate down N version base on the current version

## 範例

1. `mkdir -p db/migration` :
2. `migrate -help` :
3. `migrate create -ext sql -dir db/migration -seq init_schema` : 會在 db/migration 下自動生成兩個檔案 000001_init_schema.down.sql 和 000001_init_schema.sql
4. `migrate -path db/migration -database "postgresql://admin:admin@localhost:5433/simple_bank?sslmode=disable" -verbose up`  : 會執行 db/migration 裡所有 XXX.up.sql 裡的所有 sql 指令
   - postgres container doesn't enable SSL by defaule, so we need to disable sslmode
- `migrate -path db/migration -database "postgresql://admin:admin@localhost:5433/simple_bank?sslmode=disable" -verbose down` : 會執行 db/migration 裡所有 XXX.down.sql 裡的所有 sql 指令
- ( 常用指令都可寫在 [Makefile](https://github.com/jasonLuFa/simplebank/blob/master/Document/Makefile.md) 裡 )
