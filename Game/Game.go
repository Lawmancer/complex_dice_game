package Game

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const wild = 4
const rounds = 4

type iGame interface {
	Start() error
}

// game implements GameInterface
type game struct {
	turnNum int
	players []player
}

// Get an instance of game
func Get(playerNames []string) iGame {

	var players []player
	var rolls []int
	for _, playerName := range playerNames {
		newPlayer := player{
			playerName,
			0,
			0,
			rolls,
		}
		players = append(players, newPlayer)
	}

	return &game{
		1,
		players,
	}
}

func (g game) Error() string {
	return fmt.Sprintf("Error: ")
}

// Start the game
func (g *game) Start() error {

	if g.players == nil || len(g.players) == 0 {
		return errors.New("no players")
	}
	g.randomFirstPlayer()

	for i, p := range g.players {
		fmt.Printf("Player %d: %s\n", i+1, p.name)
	}
	fmt.Println()

	for round := 0; round < rounds; round++ {
		fmt.Println("Starting Round", round+1)
		for i := 0; i < len(g.players); i++ {
			g.turn(i)
		}
		fmt.Println()
		g.players[0].turnOrder = g.players[0].turnOrder + len(g.players) // move to end of line
		g.resortPlayers()
	}

	return nil
}

func (g *game) turn(playerNum int) (rolls []int) {
	fmt.Println(g.players[playerNum].name, " is taking their turn.")
	// TODO: probably some kind of while loop? Recursion? until no choices left.
	// todo: take turn (roll dice, wait for response, update player totals, continue until no rolls left)
	return
}

func (g *game) end() {
	// Decide Winner
	// TODO: After all four rounds have been completed the player with the lowest combined score wins.
}

func (g *game) randomFirstPlayer() {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(len(g.players))
	g.players[r].turnOrder = 1
	next := 2
	for i := 0; i < len(g.players); i++ {
		if g.players[i].turnOrder == 0 {
			g.players[i].turnOrder = next
			next++
		}
	}
	g.resortPlayers()
}

func (g *game) resortPlayers() {
	sort.Slice(g.players, func(i, j int) bool {
		return g.players[i].turnOrder < g.players[j].turnOrder
	})
}
