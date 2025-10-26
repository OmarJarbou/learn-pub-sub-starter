package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)

func handlerWar(gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(rw gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Print("> ")
		outcome, _, _ := gs.HandleWar(rw)
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NACK_REQUEUE // so other player/consumer can try to consume it
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NACK_DISCARD
		case gamelogic.WarOutcomeOpponentWon:
		case gamelogic.WarOutcomeYouWon:
		case gamelogic.WarOutcomeDraw:
			return pubsub.ACK
		}
		return pubsub.NACK_DISCARD
	}
}
