package mq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	WaitExchange   = "wait_exchange"
	WaitQueue      = "wait_queue"
	WaitRoutingKey = "wait_routing_key"

	RetryCountHeader = "x-retry-count"
	DirectRouting    = "direct"
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

func (consumer *RabbitConsumer) DeclareQueue(name string, args amqp.Table) error {
	q, err := consumer.Ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		args,  // arguments
	)

	consumer.queues[name] = &q
	return err
}

func (consumer *RabbitConsumer) QueueBind(name, exchange, routingKey string) error {
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

func (producer *RabbitConsumer) Publish(exchange, routingKey string, message []byte, headers amqp.Table) (err error) {

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
			Headers:      headers,
			Body:         message,
		},
	)

	return err
}

func (consumer *RabbitConsumer) CreateWaitQueue(delay int32, waitQueue string, waitKey string, exchange string, routingKey string) error {

	err := consumer.DeclareQueue(waitQueue, amqp.Table{
		"x-dead-letter-exchange":    exchange,
		"x-dead-letter-routing-key": routingKey,
		"x-message-ttl":             delay * 1000,     // turns into DLM after this time peroid
		"x-expires":                 delay * 1000 * 2, // delete queue if not access
	})

	if err != nil {
		return err
	}

	err = consumer.QueueBind(waitQueue, WaitExchange, waitKey)

	if err != nil {
		return err
	}

	return nil
}
