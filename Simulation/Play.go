package Simulation

import (
	"fmt"
	"github.com/DiceGame/Game"
)

// wait for rolls channel, call choose

// PlayGame simulates a game with 4 players and five
func PlayGame() {
	players := []string{"John Lennon", "Paul McCartney", "George Harrison", "Ringo Starr"}
	game := Game.Get(players)
	err := game.Start()
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
}
