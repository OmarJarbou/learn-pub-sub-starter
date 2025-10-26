package gamelogic

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

type GameState struct {
	Player Player
	Paused bool
	mu     *sync.RWMutex
}

func (gs *GameState) CommandSpam(words []string, ch *amqp.Channel) error {
	if len(words) < 2 {
		return errors.New("usage: spam <loops>")
	}

	loops, err := strconv.Atoi(words[1])
	if err != nil {
		return errors.New("spam command only accepts integers: " + err.Error())
	}

	for i := loops; i >= 0; i-- {
		game_log := routing.GameLog{
			CurrentTime: time.Now(),
			Message:     GetMaliciousLog(),
			Username:    gs.GetUsername(),
		}
		err := pubsub.PublishGob(ch, routing.ExchangePerilTopic, routing.GameLogSlug+"."+gs.GetUsername(), game_log)
		if err != nil {
			fmt.Println("Error while publishing the war result message: " + err.Error())
		}
	}
	fmt.Println(words[1] + " logs has been published successfully!")

	return nil
}

func NewGameState(username string) *GameState {
	return &GameState{
		Player: Player{
			Username: username,
			Units:    map[int]Unit{},
		},
		Paused: false,
		mu:     &sync.RWMutex{},
	}
}

func (gs *GameState) resumeGame() {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Paused = false
}

func (gs *GameState) pauseGame() {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Paused = true
}

func (gs *GameState) isPaused() bool {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.Paused
}

func (gs *GameState) addUnit(u Unit) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Player.Units[u.ID] = u
}

func (gs *GameState) removeUnitsInLocation(loc Location) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	for k, v := range gs.Player.Units {
		if v.Location == loc {
			delete(gs.Player.Units, k)
		}
	}
}

func (gs *GameState) UpdateUnit(u Unit) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Player.Units[u.ID] = u
}

func (gs *GameState) GetUsername() string {
	return gs.Player.Username
}

func (gs *GameState) getUnitsSnap() []Unit {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	Units := []Unit{}
	for _, v := range gs.Player.Units {
		Units = append(Units, v)
	}
	return Units
}

func (gs *GameState) GetUnit(id int) (Unit, bool) {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	u, ok := gs.Player.Units[id]
	return u, ok
}

func (gs *GameState) GetPlayerSnap() Player {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	Units := map[int]Unit{}
	for k, v := range gs.Player.Units {
		Units[k] = v
	}
	return Player{
		Username: gs.Player.Username,
		Units:    Units,
	}
}
