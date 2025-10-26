package pubsub

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	DURABLE SimpleQueueType = iota
	TRANSIENT
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	amqp_queue := amqp.Queue{}

	chann, err := conn.Channel()
	if err != nil {
		return nil, amqp_queue, errors.New("Error while creating a new channel from the RabbitMQ AMQP connection: " + err.Error())
	}

	queue_args := amqp.Table{
		"x-dead-letter-exchange": "peril_dlx",
	}

	queue, err := chann.QueueDeclare(queueName, (queueType == DURABLE), (queueType == TRANSIENT), (queueType == TRANSIENT), false, queue_args)
	if err != nil {
		return chann, amqp_queue, errors.New("Error while creating a new queue using the RabbitMQ AMQP channel: " + err.Error())
	}
	amqp_queue = queue

	err = chann.QueueBind(queue.Name, key, exchange, false, nil)
	if err != nil {
		return chann, amqp_queue, errors.New("Error while binding a queue to an exchange: " + err.Error())
	}

	return chann, amqp_queue, nil
}
