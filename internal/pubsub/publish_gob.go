package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(val)
	if err != nil {
		return errors.New("Error while converting/encoding message to gob format: " + err.Error())
	}

	published_msg := amqp.Publishing{
		ContentType: "application/gob",
		Body:        buffer.Bytes(),
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, published_msg)
	if err != nil {
		return errors.New("Error while publishing the message to the RabbitMQ exchange: " + err.Error())
	}

	return nil
}
