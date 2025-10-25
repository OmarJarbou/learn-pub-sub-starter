package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(am gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		move_outcome := gs.HandleMove(am)
		if move_outcome == gamelogic.MoveOutComeSafe || move_outcome == gamelogic.MoveOutcomeMakeWar {
			return pubsub.ACK
		} else {
			return pubsub.NACK_DISCARD
		}
	}
}
