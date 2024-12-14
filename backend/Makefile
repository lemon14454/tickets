DB_CONTAINER_NAME=ticket_db
DB_CONNECTION_STRING=postgresql://root:secret@localhost:5432/ticket?sslmode=disable
DB_NAME=ticket
DB_USERNAME=root
DB_PASSWORD=secret

NETWORK_NAME=ticket-network
MIGRATION_PATH=db/migration

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server

network:
	docket network create -d bridge $(NETWORK_NAME)

postgres:
	docker run --name $(DB_CONTAINER_NAME) --network $(NETWORK_NAME) -p 5432:5432 -e POSTGRES_USER=$(DB_USERNAME) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:16-alpine

createdb:
	docker exec -it $(DB_CONTAINER_NAME) createdb --username=$(DB_USERNAME) --owner=$(DB_USERNAME) $(DB_NAME)

dropdb:
	docker exec -it $(DB_CONTAINER_NAME) dropdb --username=$(DB_USERNAME) --owner=$(DB_USERNAME) $(DB_NAME)

migrateup:
	migrate -path $(MIGRATION_PATH) -database $(DB_CONNECTION_STRING) -verbose up

migratedown:
	migrate -path $(MIGRATION_PATH) -database $(DB_CONNECTION_STRING) -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go
