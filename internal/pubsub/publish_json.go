package pubsub

import (
	"context"
	"encoding/json"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	message, err := json.Marshal(val)
	if err != nil {
		return errors.New("Error while converting/encoding message to json format: " + err.Error())
	}

	published_msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, published_msg)
	if err != nil {
		return errors.New("Error while publishing the message to the RabbitMQ exchange: " + err.Error())
	}

	return nil
}
