.PHONY: build cleardb clearuser simulate spam create_user create_event

build:
	go build -o client main.go client.go

cleardb:
	docker exec -it tickets-postgres-1 psql -U root -d ticket -c "TRUNCATE TABLE event_zones, events, tickets, orders RESTART IDENTITY;"

clearcache:
	docker exec -it tickets-redis-1 redis-cli flushall

clearuser:
	docker exec -it tickets-postgres-1 psql -U root -d ticket -c "TRUNCATE TABLE users CASCADE"

simulate:
	./client -attempt 3000 -event_id 1

spam:
	./client -attempt 5 -event_id 1 -same_user true

spam_ip:
	./client -attempt 8 -port 8081 -event_id 1

create_user:
	./client -action user

create_event:
	./client -zone A,3,10\|B,3,10\|C,3,10 -action event

