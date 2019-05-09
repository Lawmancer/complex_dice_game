package Simulation

import (
	"fmt"
	"github.com/DiceGame/Dice"
)

// wait for rolls channel, call choose

// PlayGame simulates a game with 4 players and five
func PlayGame() {
	playerNames := []string{"John Lennon", "Paul McCartney", "George Harrison", "Ringo Starr"}

	players := make(map[int]string)
	game := Dice.Game{}
	for _, name := range playerNames {
		id, err := game.Register(name)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		players[id] = name
	}

	err := game.Start()
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
