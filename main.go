package main

import (
	"fmt"
	"github.com/DiceGame/Simulation"
)

func main() {
	fmt.Print("-- DiceGame Running --\n\n")

	Simulation.PlayGame()

	fmt.Print("\n-- DiceGame Ended --\n\n")
}
