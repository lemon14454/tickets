DB_CONTAINER_NAME=ticket_db
DB_CONNECTION_STRING=postgresql://root:secret@localhost:5432/ticket?sslmode=disable
DB_NAME=ticket
DB_USERNAME=root
DB_PASSWORD=secret

.PHONY: test build server

test:
	go test -v -cover ./...

build:
	go build -o server main.go handler.go

server:
	./server
