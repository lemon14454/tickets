package main

import (
	"context"
	"log"
	"ticket/backend/api"
	db "ticket/backend/db/sqlc"
	"ticket/backend/util"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config := util.LoadConfig(".")

	dbpool, err := pgxpool.New(context.Background(), config.DB_SOURCE)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v \n", err)
	}
	defer dbpool.Close()

	store := db.NewStore(dbpool)

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("Failed to create server: %v \n", err)
	}
	defer server.Close()

	err = server.Start(config.HTTP_SERVER_ADDRESS)
	if err != nil {
		log.Fatalf("Failed to start server: %v \n", err)
	}
}
