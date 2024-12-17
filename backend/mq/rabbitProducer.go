package mq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MQProducer interface {
	Publish(message []byte)
}

type RabbitProducer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitProducer(address string) (MQProducer, error) {
	conn, err := amqp.Dial(address)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	producer := &RabbitProducer{
		conn: conn,
		ch:   ch,
	}

	return producer, nil
}

func (producer *RabbitProducer) CreateQueue() (*amqp.Queue, error) {
	q, err := producer.ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	return &q, err
}

func (producer *RabbitProducer) Publish(message []byte) {
	q, err := producer.ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = producer.ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         message,
		},
	)
}
