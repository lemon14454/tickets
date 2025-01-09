package mq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var DirectRouting = "direct"

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

func (producer *RabbitProducer) DeclareExchange(name, exchangeType string) error {
	err := producer.Ch.ExchangeDeclare(
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

func (producer *RabbitProducer) DeclareQueue(name string, args amqp.Table) error {
	q, err := producer.Ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		args,  // arguments
	)

	producer.queues[name] = &q
	return err
}

func (producer *RabbitProducer) QueueBind(name, exchange, routingKey string) error {
	q, ok := producer.queues[name]
	if !ok {
		return fmt.Errorf("Queue %s not exists", name)
	}

	err := producer.Ch.QueueBind(
		q.Name,     // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil,
	)

	return err
}

func (producer *RabbitProducer) Publish(exchange, routingKey string, message []byte) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = producer.Ch.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
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
