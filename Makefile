DB_URL="postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"

network:
	docker network create bank-network

postgres:
	docker run --name postgres16 -p 5432:5432 --network bank-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres16 dropdb simple_bank

migrateup:
	migrate -path db/migration -database $(DB_URL) -verbose up

migrateup1:
	migrate -path db/migration -database $(DB_URL) -verbose up 1

migratedown:
	migrate -path db/migration -database $(DB_URL) -verbose down

migratedown1:
	migrate -path db/migration -database $(DB_URL) -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

dbdocs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql doc/db.dbml -o doc/schema.sql --postgres

sqlc:
	sqlc generate
	# docker run --rm -v ${pwd}:/src -w /src sqlc/sqlc generate # on windows

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -destination ./db/mock/store.go -package mockdb github.com/bacnx/simplebank/db/sqlc Store

proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration dbdocs db_schema sqlc server mock proto
