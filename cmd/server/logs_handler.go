package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerLogs() func(routing.GameLog) pubsub.Acktype {
	return func(gamelog routing.GameLog) pubsub.Acktype {
		err := gamelogic.WriteLog(gamelog)
		if err != nil {
			fmt.Println("Error while writing game log: " + err.Error())
			return pubsub.NACK_DISCARD
		}
		return pubsub.ACK
	}
}
