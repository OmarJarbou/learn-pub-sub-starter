package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	fmt.Println("Connecting with RabbitMQ server...")
	conn_string := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(conn_string)
	if err != nil {
		log.Fatal("Error while connecting to RabbitMQ server:" + err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("Connected Successfully!")

	chann, err := conn.Channel()
	if err != nil {
		log.Fatal("Error while creating a new channel from the RabbitMQ AMQP connection: " + err.Error())
		return
	}

	pause_msg := routing.PlayingState{
		IsPaused: true,
	}
	err = pubsub.PublishJSON(chann, routing.ExchangePerilDirect, routing.PauseKey, pause_msg)
	if err != nil {
		log.Fatal("Error while publishing pause message: " + err.Error())
		return
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // if user have sent "Ctrl + C", then send the signal to signalChan
	<-signalChan                            // Blocks until recieve "Ctrl + C"
	log.Println("Server gracefully stopped")
}
