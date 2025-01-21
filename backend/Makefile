DB_CONTAINER_NAME=ticket_db
DB_CONNECTION_STRING=postgresql://root:secret@localhost:5432/ticket?sslmode=disable
DB_NAME=ticket
DB_USERNAME=root
DB_PASSWORD=secret

REDIS_CONTAINER_NAME=ticket_redis
MQ_CONTAINER_NAME=ticket_mq

MIGRATION_PATH=db/migration

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server redis mq

postgres:
	docker run --name $(DB_CONTAINER_NAME) -p 5432:5432 -e POSTGRES_USER=$(DB_USERNAME) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:16-alpine

redis:
	docker run -d --name $(REDIS_CONTAINER_NAME) -p 6379:6379 redis:7.4.2-alpine

mq:
	docker run -d --name $(MQ_CONTAINER_NAME) -p 5672:5672 rabbitmq:4.0.4-alpine 

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
