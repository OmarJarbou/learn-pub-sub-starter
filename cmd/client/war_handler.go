package main

import (
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(rw gamelogic.RecognitionOfWar) pubsub.Acktype {
		outcome, winner, loser := gs.HandleWar(rw)
		var log_message string
		war_happened := false
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NACK_REQUEUE // so other player/consumer can try to consume it
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NACK_DISCARD
		case gamelogic.WarOutcomeOpponentWon:
			log_message = winner + " won a war against " + loser
			war_happened = true
		case gamelogic.WarOutcomeYouWon:
			log_message = winner + " won a war against " + loser
			war_happened = true
		case gamelogic.WarOutcomeDraw:
			log_message = "A war between " + winner + " and " + loser + " resulted in a draw"
			war_happened = true
		}
		if war_happened {
			game_log := routing.GameLog{
				CurrentTime: time.Now(),
				Message:     log_message,
				Username:    gs.GetUsername(),
			}
			err := pubsub.PublishGob(ch, routing.ExchangePerilTopic, routing.GameLogSlug+"."+gs.GetUsername(), game_log)
			if err != nil {
				fmt.Println("Error while publishing the war result message: " + err.Error())
				return pubsub.NACK_REQUEUE
			}
			fmt.Println("The war result was published successfully!")
			return pubsub.ACK
		}
		return pubsub.NACK_DISCARD
	}
}
