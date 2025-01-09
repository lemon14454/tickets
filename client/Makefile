.PHONY: build client

build:
	go build -o client main.go client.go

client:
	client -ip 127.0.0.1 -port 8080 -user testuser -pwd 12345678 -zone A,10,10\|B,10,10\|C,10,10 -attempt 10000
