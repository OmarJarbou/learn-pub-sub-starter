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

	chann, err := conn.Channel()
	if err != nil {
		log.Fatal("Error while creating a new channel from the RabbitMQ AMQP connection: " + err.Error())
		return
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal("Error while welcoming client: " + err.Error())
		return
	}

	// _, _, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queue_name, routing.PauseKey, pubsub.TRANSIENT)
	// if err != nil {
	// 	log.Fatal("Error while declaring a new queue and bind it to an exchange: " + err.Error())
	// 	return
	// }

	game_state := gamelogic.NewGameState(username)
	pauseErrorsChan := make(chan error, 1)
	moveErrorsChan := make(chan error, 1)
	warErrorsChan := make(chan error, 1)
	go func() {
		pauseErrorsChan <- pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+username, routing.PauseKey, pubsub.TRANSIENT, handlerPause(game_state))
	}()
	go func() {
		moveErrorsChan <- pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+username, routing.ArmyMovesPrefix+".*", pubsub.TRANSIENT, handlerMove(game_state, chann))
	}()
	go func() {
		warErrorsChan <- pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix, routing.WarRecognitionsPrefix+".*", pubsub.DURABLE, handlerWar(game_state))
	}()

	for {
		// Check for subscription errors in a non-blocking way
		select {
		case err := <-pauseErrorsChan:
			if err != nil {
				log.Fatal("Error while subscribing to the pause queue: " + err.Error())
				return
			}
		case err := <-moveErrorsChan:
			if err != nil {
				log.Fatal("Error while subscribing to the move queue: " + err.Error())
				return
			}
		case err := <-warErrorsChan:
			if err != nil {
				log.Fatal("Error while subscribing to the war queue: " + err.Error())
				return
			}
		default:
			// No error, continue with normal flow
		}

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
			army_move, err := game_state.CommandMove(words)
			if err != nil {
				fmt.Println(err.Error())
			}
			err = pubsub.PublishJSON(chann, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+username, army_move)
			if err != nil {
				log.Fatal("Error while publishing the move message: " + err.Error())
				return
			}
			fmt.Println("The move was published successfully!")
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
