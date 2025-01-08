package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"ticket/consumer/model"
	rbmq "ticket/consumer/mq"
	"ticket/consumer/util"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func deserialize(b []byte) (model.Message, error) {
	var msg model.Message
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}

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

	err = mq.DeclareQueue(rbmq.TicketQueue)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v \n", err)
	}

	err = mq.QueueBind(rbmq.TicketQueue, "", rbmq.TicketExchange)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v \n", err)
	}

	msgs, err := mq.Consume(rbmq.TicketQueue)
	if err != nil {
		log.Fatalf("Failed to consume: %v \n", err)
	}

	go func() {
		for d := range msgs {
			msg, err := deserialize(d.Body)
			if err != nil {
				log.Printf("Failed to deserialize: %v \n", err)
			}

			log.Printf("Received Event ID: %v", msg)

			var evt model.Event
			result := db.First(&evt, msg.EventID)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				log.Printf("Unable to find Event: %v \n", err)
				// normally it won't happend
				return
			}

			var eventZone []model.EventZone
			result = db.Where("event_id = ?", msg.EventID).Find(&eventZone)
			if result.Error != nil {
				log.Printf("Unable to find EventZone: %v \n", err)
				return
			}

			tickets := make([]model.Ticket, 0)
			for _, zone := range eventZone {
				for row := range zone.Rows {
					for seat := range zone.Seats {
						tickets = append(tickets, model.Ticket{
							EventID: evt.ID,
							ZoneID:  zone.ID,
							Row:     row + 1,
							Seat:    seat + 1,
						})
					}
				}
			}

			err = db.Transaction(func(tx *gorm.DB) error {
				if err := db.Create(&tickets).Error; err != nil {
					return err
				}

				if err := db.Model(&evt).Update("status", model.EventStatusCreated).Update("updated_at", time.Now()).Error; err != nil {
					return err
				}

				return nil
			})

			d.Ack(false)
		}
	}()

	var forever chan struct{}
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
