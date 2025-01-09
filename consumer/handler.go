package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"ticket/consumer/model"
	rbmq "ticket/consumer/mq"

	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

var (
	TicketExchange   = "ticket_exchange"
	TicketQueue      = "ticket_queue"
	TicketRoutingKey = "ticket_routing_key"
)

var maxRetryConut int32 = 10

type Handler struct {
	db *gorm.DB
	mq *rbmq.RabbitConsumer
}

func NewHandler(db *gorm.DB, mq *rbmq.RabbitConsumer) *Handler {
	return &Handler{
		db: db,
		mq: mq,
	}
}

func deserialize(b []byte) (model.Message, error) {
	var msg model.Message
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}

func (handler *Handler) startListenEvent() {
	err := handler.mq.DeclareExchange(TicketExchange, rbmq.DirectRouting)
	if err != nil {
		log.Fatalf("Failed to declare ticket exchange: %v \n", err)
	}

	err = handler.mq.DeclareExchange(rbmq.WaitExchange, rbmq.DirectRouting)
	if err != nil {
		log.Fatalf("Failed to declare wait exchange: %v \n", err)
	}

	err = handler.mq.DeclareQueue(TicketQueue, amqp091.Table{})
	if err != nil {
		log.Fatalf("Failed to declare queue: %v \n", err)
	}

	err = handler.mq.QueueBind(TicketQueue, TicketExchange, TicketRoutingKey)
	if err != nil {
		log.Fatalf("Failed to Bind queue: %v \n", err)
	}

	msgs, err := handler.mq.Consume(TicketQueue)
	if err != nil {
		log.Fatalf("Failed to consume: %v \n", err)
	}

	go handler.eventCreateHandler(msgs)

	var forever chan struct{}
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (handler *Handler) eventCreateHandler(msgs <-chan amqp091.Delivery) {
	for d := range msgs {
		msg, err := deserialize(d.Body)
		if err != nil {
			log.Printf("Failed to deserialize: %v \n", err)
		}

		log.Printf("Received Event ID: %v, start create", msg)

		err = model.CreateEvent(handler.db, msg.EventID)

		if err != nil {
			// pass message to DLE by sending to a wait queue that nobody consume, which trigger TTL
			err := handler.sendBackToEventQueue(d)

			if err != nil {
				// adding a notification would be a good idea, since event creation isn't complete
				log.Printf("[Possible Message Loss] Failed sending EventID %d back to queue due to: %v", msg.EventID, err)

				// mark event status as failed
				err := model.MarkEventFailure(handler.db, msg.EventID)

				if err != nil {
					log.Printf("Failed marking EventID %d as failure: %v", msg.EventID, err)
				}
			}
		}

		// If message failed -> pass to wait queue -> no need to NACK
		d.Ack(false)
	}
}

func (handler *Handler) sendBackToEventQueue(d amqp091.Delivery) error {
	delay := 3
	waitQueue := fmt.Sprintf("%s@%d", rbmq.WaitQueue, delay)
	waitKey := fmt.Sprintf("%s@%d", rbmq.WaitRoutingKey, delay)

	err := handler.mq.CreateWaitQueue(delay, waitQueue, waitKey, TicketExchange, TicketRoutingKey)
	if err != nil {
		return err
	}

	var retryCount int32 = 0
	count, ok := d.Headers[rbmq.RetryCountHeader]
	if ok {
		retryCount = count.(int32)
	}
	log.Printf("Receive Retry Count: %v", retryCount)

	if retryCount > maxRetryConut {
		return fmt.Errorf("Max event create retry exceeded")
	}

	err = handler.mq.Publish(rbmq.WaitExchange, waitKey, d.Body, amqp091.Table{
		rbmq.RetryCountHeader: retryCount + 1,
	})

	if err != nil {
		return err
	}

	return nil
}
