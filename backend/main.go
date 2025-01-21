package main

import (
	"context"
	"log"
	"ticket/backend/api"
	db "ticket/backend/db/sqlc"
	"ticket/backend/util"

	"github.com/jackc/pgx/v5/pgxpool"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config := util.LoadConfig(".")

	dbpool, err := pgxpool.New(context.Background(), config.DB_SOURCE)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v \n", err)
	}
	defer dbpool.Close()
	runDBMigration(config.MIGRATION_URL, config.DB_SOURCE)

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

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatalf("connect create new migrate instance: %v", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrate up: %v", err)
	}

	log.Print("db migrated successfully")
}
