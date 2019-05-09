package Dice

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const wild = 4
const rounds = 4

// Game implements GameInterface
type Game struct {
	round   int
	players []player
}

// Start the Game
func (g *Game) Start() error {
	if g.players == nil || len(g.players) == 0 {
		return errors.New("no players")
	}

	fmt.Printf("Starting a new Game with %d players…\n", len(g.players))
	g.randomFirstPlayer()
	g.listPlayers()
	for g.round < rounds {
		g.round++
		fmt.Println("Starting Round", g.round+1)
		for i := range g.players {
			g.turn(i)
		}
		fmt.Println()
		g.players[0].turnOrder = g.players[0].turnOrder + len(g.players) // move to end of line
		g.resortPlayers()
	}

	return nil
}

// Register a new player before the Game starts
func (g *Game) Register(name string) (id int, err error) {
	if g.round > 0 {
		return 0, errors.New("the game has already started")
	}

	id = len(g.players) + 1
	newPlayer := player{
		id,
		name,
		0,
		0,
	}
	g.players = append(g.players, newPlayer)

	return id, nil
}

func (g *Game) turn(playerNum int) (rolls []int) {
	fmt.Println(g.players[playerNum].name, " is taking their turn.")
	// TODO: probably some kind of while loop? Recursion? until no choices left.
	// todo: take turn (roll dice, wait for response, update player totals, continue until no rolls left)
	return
}

func (g *Game) end() {
	// Decide Winner
	// TODO: After all four rounds have been completed the player with the lowest combined score wins.
}

func (g *Game) listPlayers() {
	for i, p := range g.players {
		fmt.Printf("Player %d: %s\n", i+1, p.name)
	}
	fmt.Println()
}

func (g *Game) randomFirstPlayer() {
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

func (g *Game) resortPlayers() {
	sort.Slice(g.players, func(i, j int) bool {
		return g.players[i].turnOrder < g.players[j].turnOrder
	})
}
