postgres: 
	docker run --name postgres12 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -p 5433:5432 -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=admin --owner=admin simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank --username=admin

migrateup:
	migrate -path db/migration -database "postgresql://admin:admin@localhost:5433/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://admin:admin@localhost:5433/simple_bank?sslmode=disable" -verbose down

sqlc:
	docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

server :
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/jasonLuFa/simplebank/db/sqlc Store 

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock