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
	fmt.Println("Starting Peril client...")

	fmt.Println("Connecting with RabbitMQ server...")
	conn_string := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(conn_string)
	if err != nil {
		log.Fatal("Error while connecting to RabbitMQ server: " + err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("Connected Successfully!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal("Error while welcoming client: " + err.Error())
		return
	}

	queue_name := routing.PauseKey + "." + username
	pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queue_name, routing.PauseKey, pubsub.TRANSIENT)

	game_state := gamelogic.NewGameState(username)
	for {
		words := gamelogic.GetInput()
		if len(words) < 1 {
			fmt.Println("you must enter a command from possible commands list.")
			continue
		}
		if words[0] == "spawn" {
			err = game_state.CommandSpawn(words)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else if words[0] == "move" {
			_, err = game_state.CommandMove(words)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else if words[0] == "status" {
			game_state.CommandStatus()
		} else if words[0] == "help" {
			gamelogic.PrintClientHelp()
		} else if words[0] == "spam" {
			fmt.Println("Spamming not allowed yet!")
		} else if words[0] == "quit" {
			gamelogic.PrintQuit()
			break
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
