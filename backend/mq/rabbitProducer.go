package mq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitProducer struct {
	Conn   *amqp.Connection
	Ch     *amqp.Channel
	queues map[string]*amqp.Queue
}

func NewRabbitProducer(address string) (*RabbitProducer, error) {
	conn, err := amqp.Dial(address)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	producer := &RabbitProducer{
		Conn:   conn,
		Ch:     ch,
		queues: make(map[string]*amqp.Queue),
	}

	return producer, nil
}

func (producer *RabbitProducer) DeclareQueue(name string) error {
	q, err := producer.Ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	producer.queues[name] = &q
	return err
}

func (producer *RabbitProducer) Publish(queue string, message []byte) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q, ok := producer.queues[queue]
	if !ok {
		return fmt.Errorf("Queue Name: %v not found", queue)
	}

	err = producer.Ch.PublishWithContext(ctx,
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

	return err
}

func (producer *RabbitProducer) Close() {
	producer.Conn.Close()
	producer.Ch.Close()
}
