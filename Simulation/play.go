package Simulation

import (
	"fmt"
	"github.com/DiceGame/Dice"
	"sync"
)

// PlayGame simulates a game with 4 players and five
func PlayGame(waitGroup sync.WaitGroup) {
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

	go game.Start(waitGroup) // TODO: unhandled error

	for game.Done == false {
		selections := 0
		turn := game.WaitForTurn()

		// TODO: make new method
		fmt.Println("-----------------------")
		fmt.Println("turn.ActivePlayer.Id: ", turn.ActivePlayer.Id)
		fmt.Println("-----------------------")

		name := players[turn.ActivePlayer.Id]
		rolls := turn.Rolls
		fmt.Printf("Player %s sees rolls: %v\n", name, rolls)
		for i, r := range rolls {
			if r.Roll.Value() == Dice.Wild {
				turn.Rolls[i].Selected = true
				fmt.Printf("Player %s keeps: %d\n", name, turn.Rolls[i].Roll.Value())
				selections++
			} else if r.Roll.Value() == 1 || r.Roll.Value() == 2 {
				turn.Rolls[i].Selected = true
				fmt.Printf("Player %s keeps: %d\n", name, turn.Rolls[i].Roll.Value())
				selections++
			}
		}
		if selections == 0 {
			low := 0
			for i, r := range rolls {
				if i == 0 {
					// first is already lowers
					continue
				} else if r.Roll.Value() < turn.Rolls[i].Roll.Value() {
					low = i
				}
			}
			turn.Rolls[low].Selected = true
			fmt.Printf("Player %s keeps: %d\n", name, turn.Rolls[low].Roll.Value())
		}

		game.Choices <- turn
		if game.Done {
			continue
		}
	}

	fmt.Println("And I'm done")
	return
}
