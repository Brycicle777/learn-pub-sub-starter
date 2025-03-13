package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type simpleQueueType int

const (
	Durable simpleQueueType = iota
	Transient
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("Error marshalling value: %v", err)
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	})
	if err != nil {
		return fmt.Errorf("Error publishing to channel: %v", err)
	}
	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType simpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	durable := simpleQueueType == Durable
	autoDelete := simpleQueueType == Transient
	exclusive := simpleQueueType == Transient

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating connection channel: %v", err)
		return nil, amqp.Queue{}, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		exclusive,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error creating queue: %v", err)
		return nil, amqp.Queue{}, err
	}

	err = ch.QueueBind(
		queueName,
		key,
		exchange,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error binding queue: %v", err)
		return nil, amqp.Queue{}, err
	}

	return ch, q, nil
}
