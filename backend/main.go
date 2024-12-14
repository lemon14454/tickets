package main

import (
	"context"
	"log"
	"ticket/backend/api"
	db "ticket/backend/db/sqlc"
	"ticket/backend/util"

	"github.com/jackc/pgx/v5"
)

func main() {
	config := util.LoadConfig(".")

	conn, err := pgx.Connect(context.Background(), config.DB_SOURCE)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v \n", err)
	}
	defer conn.Close(context.Background())

	store := db.NewStore(conn)
	runAPIServer(config, store)
}

func runAPIServer(config *util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("Failed to create server: %v \n", err)
	}

	err = server.Start(config.HTTP_SERVER_ADDRESS)
	if err != nil {
		log.Fatalf("Failed to start server: %v \n", err)
	}
}
