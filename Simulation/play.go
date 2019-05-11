package Simulation

import (
	"fmt"
	"github.com/DiceGame/Game"
)

// PlayGame simulates a game with 4 players and five
func PlayGame() {
	playerNames := []string{"John Lennon", "Paul McCartney", "George Harrison", "Ringo Starr"}

	players := make(map[int]string)
	game := Game.Play{}
	for _, name := range playerNames {
		id, err := game.Register(name)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		players[id] = name
	}

	go game.Start()

	for game.Done == false {
		turn, done := game.WaitForTurn()
		if done {
			break
		}

		name := players[turn.ActivePlayer.Id]
		fmt.Printf("%s sees rolls: ", name)
		for _, r := range turn.Dice {
			fmt.Printf("%d ", r.Pips.Value())
		}
		fmt.Printf("\n")
		makeChoice(turn)
		game.Choices <- turn
	}

	return
}

func makeChoice(turn Game.RollResult) {
	selections := 0
	dice := turn.Dice
	for i, r := range dice {
		if r.Pips.Value() == Game.Wild || r.Pips.Value() <= 2 {
			turn.Dice[i].Selected = true
			fmt.Printf("  … keeps: %d\n", turn.Dice[i].Pips.Value())
			selections++
		}
	}
	if selections == 0 {
		low := 0
		for i, r := range dice {
			if i == 0 {
				// first is already lowers
				continue
			} else if r.Pips.Value() < turn.Dice[i].Pips.Value() {
				low = i
			}
		}
		turn.Dice[low].Selected = true
		fmt.Printf("  … keeps: %d\n", turn.Dice[low].Pips.Value())
	}
}
