package main

import (
	"log"
	rbmq "ticket/consumer/mq"
	"ticket/consumer/util"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config := util.LoadConfig(".")

	db, err := gorm.Open(postgres.Open(config.DB_SOURCE), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v \n", err)
	}

	mq, err := rbmq.NewRabbitConsumer(config.BROKER_ADDRESS)
	if err != nil {
		log.Fatalf("Failed to create server: %v \n", err)
	}

	handler := NewHandler(db, mq)
	handler.startListenEvent()
}
