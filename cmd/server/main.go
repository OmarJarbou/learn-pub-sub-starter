package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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
		log.Fatal("Error while connecting to RabbitMQ server: " + err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("Connected Successfully!")

	chann, err := conn.Channel()
	if err != nil {
		log.Fatal("Error while creating a new channel from the RabbitMQ AMQP connection: " + err.Error())
		return
	}

	gamelogic.PrintServerHelp()
	for {
		words := gamelogic.GetInput()
		if len(words) != 1 {
			fmt.Println("you must enter a one word command from possible commands list.")
			continue
		}
		if words[0] == "pause" {
			fmt.Println("Publishing pause message...")
			pause_msg := routing.PlayingState{
				IsPaused: true,
			}
			err = pubsub.PublishJSON(chann, routing.ExchangePerilDirect, routing.PauseKey, pause_msg)
			if err != nil {
				log.Fatal("Error while publishing pause message: " + err.Error())
				return
			}
		} else if words[0] == "resume" {
			fmt.Println("Publishing resume message...")
			pause_msg := routing.PlayingState{
				IsPaused: false,
			}
			err = pubsub.PublishJSON(chann, routing.ExchangePerilDirect, routing.PauseKey, pause_msg)
			if err != nil {
				log.Fatal("Error while publishing resume message: " + err.Error())
				return
			}
		} else if words[0] == "exit" {
			fmt.Println("Exiting...")
			break
		} else if words[0] == "help" {
			gamelogic.PrintServerHelp()
		} else {
			fmt.Println("Unknown Command")
		}
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // if user have sent "Ctrl + C", then send the signal to signalChan
	<-signalChan                            // Blocks until recieve "Ctrl + C"
	log.Println("Server gracefully stopped")
}
