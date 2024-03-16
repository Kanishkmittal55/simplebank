
# Set Default directory and name
DIR ?= db/migration

# The name of your Docker Compose file
COMPOSE_FILE=docker-compose.yaml

# Variables
DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

up:
	docker-compose up

down:
	docker-compose -f $(COMPOSE_FILE) down --rmi all --volumes
	@echo "Containers, networks, volumes, and images removed successfully."

startdb: postgres createdb migrateup

postgres:
	docker run --name postgres15 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres15 dropdb simple_bank

runRootDb:
	docker exec -it postgres15 psql -U root simple_bank

logs:
	docker logs postgres15

migrateup:
	 migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	 migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	 migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	 migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover -short -count=1 ./...

server:
	go run main.go

mock:
	 mockgen -package mockdb -destination db/mock/store.go github.com/kanishkmittal55/simplebank/db/sqlc Store

enter:
	docker exec -it 24cd1e5f7aaa /bin/bash

enter_into_psql:
	docker exec -it 24cd1e5f7aaa psql -U ledger_migration -d ledger

simplebankcont:
	 docker run --name simplebank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@08346868a80d:5432/simple_bank?sslmode=disable" simplebank:latest

# For adding an endpoint with one or many new tables added in the database
# 1. We add the migration up and down files...
# Example Usage - make migratefiles DIR=your_directory NAME=your_migration_name , DIR is defaulted to db/migration
migratefiles:
	 migrate create -ext sql -dir $(DIR) -seq $(NAME)

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/db/db.dbml

db_docs:
	dbdocs build docs/db/db.dbml

proto:
	rm -f pb/*.go
	rm -f docs/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
    --openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto && \
	statik -src=./docs/swagger -dest=./docs -f


evans:
	evans --host localhost --port 8080 --reflection --package pb --service simplebank statik

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: postgres createdb dropdb migrateup migratedown sqlc server mock migrateup1 migratedown1 simplebankcont down up db_docs proto evans startdb redis db_schema new_migration
