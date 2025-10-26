package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(am gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		move_outcome := gs.HandleMove(am)
		if move_outcome == gamelogic.MoveOutComeSafe {
			return pubsub.ACK
		} else if move_outcome == gamelogic.MoveOutcomeMakeWar {
			war_message := gamelogic.RecognitionOfWar{
				Attacker: gs.Player,
				Defender: am.Player,
			}
			err := pubsub.PublishJSON(ch, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix+"."+gs.GetUsername(), war_message)
			if err != nil {
				fmt.Println("Error while publishing the war message: " + err.Error())
				return pubsub.NACK_REQUEUE
			}
			fmt.Println("The war was published successfully!")
			return pubsub.ACK
		} else {
			return pubsub.NACK_DISCARD
		}
	}
}
