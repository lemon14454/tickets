package mq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

var TicketExchange = "tickets"
var TicketQueue = ""
var DirectRouting = "direct"

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

	consumer := &RabbitConsumer{
		Conn:   conn,
		Ch:     ch,
		queues: make(map[string]*amqp.Queue),
	}

	return consumer, nil
}

func (consumer *RabbitConsumer) DeclareExchange(name, exchangeType string) error {
	err := consumer.Ch.ExchangeDeclare(
		name,         // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	return err
}

func (consumer *RabbitConsumer) DeclareQueue(name string) error {
	q, err := consumer.Ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	err = consumer.Ch.Qos(
		1,
		0,
		false,
	)

	consumer.queues[name] = &q
	return err
}

func (consumer *RabbitConsumer) QueueBind(name, routingKey, exchange string) error {
	q, ok := consumer.queues[name]
	if !ok {
		return fmt.Errorf("Queue %s not exists", name)
	}

	err := consumer.Ch.QueueBind(
		q.Name,     // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil,
	)

	return err
}

func (consumer *RabbitConsumer) Consume(queue string) (<-chan amqp.Delivery, error) {

	q, ok := consumer.queues[queue]
	if !ok {
		return make(chan amqp.Delivery), fmt.Errorf("Queue Name: %v not found", queue)
	}

	msgs, err := consumer.Ch.Consume(
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
