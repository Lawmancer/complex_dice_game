package Simulation

import (
	"fmt"
	"github.com/DiceGame/Dice"
)

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

	go game.Start()

	for game.Done == false {

		turn, err := game.WaitForTurn()
		if err != nil {
			break
		}

		name := players[turn.ActivePlayer.Id]
		fmt.Printf("%s sees rolls: ", name)
		for _, r := range turn.Rolls {
			fmt.Printf("%d ", r.Roll.Value())
		}
		fmt.Printf("\n")
		makeChoice(turn)
		game.Choices <- turn
	}

	return
}

func makeChoice(turn Dice.Turn) {
	selections := 0
	rolls := turn.Rolls
	for i, r := range rolls {
		if r.Roll.Value() == Dice.Wild {
			turn.Rolls[i].Selected = true
			fmt.Printf("  … keeps: %d\n", turn.Rolls[i].Roll.Value())
			selections++
		} else if r.Roll.Value() == 1 || r.Roll.Value() == 2 {
			turn.Rolls[i].Selected = true
			fmt.Printf("  … keeps: %d\n", turn.Rolls[i].Roll.Value())
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
		fmt.Printf("  … keeps: %d\n", turn.Rolls[low].Roll.Value())
	}
}
