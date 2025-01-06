package mq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitConsumer struct {
	Conn   *amqp.Connection
	Ch     *amqp.Channel
	queues map[string]*amqp.Queue
}

func NewRabbitConsumer(address string) (*RabbitConsumer, error) {
	conn, err := amqp.Dial(address)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	producer := &RabbitConsumer{
		Conn:   conn,
		Ch:     ch,
		queues: make(map[string]*amqp.Queue),
	}

	return producer, nil
}

func (producer *RabbitConsumer) DeclareQueue(name string) error {
	q, err := producer.Ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	err = producer.Ch.Qos(
		1,
		0,
		false,
	)

	producer.queues[name] = &q
	return err
}

func (producer *RabbitConsumer) Consume(queue string) (<-chan amqp.Delivery, error) {

	q, ok := producer.queues[queue]
	if !ok {
		return make(chan amqp.Delivery), fmt.Errorf("Queue Name: %v not found", queue)
	}

	msgs, err := producer.Ch.Consume(
		q.Name,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	return msgs, err
}
