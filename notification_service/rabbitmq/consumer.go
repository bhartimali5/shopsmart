package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func ConsumeEvents(exchangeName, queueName, kind, bindingKey string) (<-chan amqp.Delivery, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	err = ch.ExchangeDeclare(exchangeName, kind, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	q, err := ch.QueueDeclare(queueName, false, false, true, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.QueueBind(q.Name, bindingKey, exchangeName, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("[*] Waiting for messages in %s", q.Name)
	return msgs, nil
}
