postgres:
	docker run --name postgres-techschool -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it postgres-techschool createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres-techschool dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate
	# docker run --rm -v ${pwd}:/src -w /src sqlc/sqlc generate # on windows

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -destination ./db/mock/store.go -package mockdb github.com/bacnx/simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc server mock
