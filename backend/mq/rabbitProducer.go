package mq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var TicketExchange = "tickets"
var DirectRouting = "direct"

type RabbitProducer struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
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
		Conn: conn,
		Ch:   ch,
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
