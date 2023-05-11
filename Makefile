DB_URL=postgresql://admin:admin@localhost:5433/simple_bank?sslmode=disable
# OSFLAG :=
# ifeq ($(OS),Windows_NT)
# 	OSFLAG = WIN
# else
# 	UNAME_S := $(uname -s)
# 	ifeq ($(UNAME_S),Linux)
# 		OSFLAG = LINUX
# 	endif
# 	ifeq ($(UNAME_S),Darwin)
# 		OSFLAG = OSX
# 	endif
# endif

postgres: 
	docker run --name postgres12 --network bank-network -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -p 5433:5432 -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=admin --owner=admin simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank --username=admin

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

server :
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/jasonLuFa/simplebank/db/sqlc Store 

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc -f

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock proto