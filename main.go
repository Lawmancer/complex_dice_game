package main

import (
	"fmt"
	"github.com/DiceGame/Simulation"
	"sync"
)

var waitGroup sync.WaitGroup // 1

func main() {
	fmt.Print("-- DiceGame Running --\n\n")
	waitGroup.Add(1)
	Simulation.PlayGame(waitGroup)
	waitGroup.Wait()
	fmt.Print("\n-- DiceGame Ended --\n\n")
}
