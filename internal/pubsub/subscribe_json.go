package pubsub

import (
	"encoding/json"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

const (
	ACK Acktype = iota
	NACK_REQUEUE
	NACK_DISCARD
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) Acktype,
) error {
	chann, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return errors.New("Error while declaring a new queue and bind it to an exchange: " + err.Error())
	}

	deliveries_chann, err := chann.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return errors.New("Error while creating the consuming channel: " + err.Error())
	}

	// Channel to capture errors from the goroutine
	errChan := make(chan error, 1)

	go func() {
		var body T
		for delivery := range deliveries_chann {
			err := json.Unmarshal(delivery.Body, &body)
			if err != nil {
				errChan <- errors.New("Error while decoding message's body: " + err.Error())
				return
			}
			ack := handler(body)
			switch ack {
			case ACK:
				delivery.Ack(false)
				fmt.Println("ACK: Message processed successfully, remove it from the queue")
			case NACK_REQUEUE:
				delivery.Nack(false, true)
				fmt.Println("NACK_REQUEUE: Requeue the failed message to proccess it again")
			case NACK_DISCARD:
				delivery.Nack(false, false)
				fmt.Println("NACK_DISCARD: Discard the failed message and send it to dead letter queue")
			}
		}
		errChan <- nil // Signal completion without error
	}()

	// Return the first error from the goroutine
	return <-errChan
}
