package Dice

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// Wild counts as 0
const Wild = 4
const wildValue = 0
const rounds = 4
const totalDice = 5

// Game implements GameInterface
type Game struct {
	round   int
	players []player
	action  chan Turn
	Choices chan Turn
	Done    bool
}

// Start the Game
func (g *Game) Start(waitgroup sync.WaitGroup) error {
	defer waitgroup.Done()
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
			g.takeTurn(i)
		}
		fmt.Println()
		g.players[0].turnOrder = g.players[0].turnOrder + len(g.players) // move to end of line
		g.resortPlayers()
	}

	g.end()
	g.Done = true

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

func (g *Game) takeTurn(playerNum int) {
	fmt.Println(g.players[playerNum].name, "is taking their turn.")
	rand.Seed(time.Now().UnixNano()) // seed once per turn

	diceRemaining := totalDice
	p := &g.players[playerNum]
	for diceRemaining > 0 {
		fmt.Println("===============================")
		fmt.Println("diceRemaining", diceRemaining)
		fmt.Println("===============================")
		var rolls []roll
		t := Turn{
			g.round,
			p,
			rolls,
		}
		d := diceRemaining
		for d > 0 {
			rnd := rand.Intn(5) + 1
			var r roll
			r.Roll.value = rnd
			t.Rolls = append(t.Rolls, r)
			d--
		}
		// TODO: Make copy and compare one that comes back
		g.action <- t
		g.Choices = make(chan Turn)
		chosen := 0
		for choices := range g.Choices {
			chosen = g.processChoices(choices)
			close(g.Choices)
		}
		diceRemaining -= chosen
	}

	return
}

// WaitForTurn is used outside the game to wait for a Turn to be playable
func (g *Game) WaitForTurn() Turn {
	var t Turn
	g.action = make(chan Turn)
	for turn := range g.action {
		t = turn
		close(g.action)
	}

	return t
}

func (g *Game) processChoices(response Turn) (chosen int) {
	total := 0
	for _, r := range response.Rolls {
		if r.Selected {
			if r.Roll.value == Wild {
				total += wildValue
			} else {
				total += r.Roll.value
			}
			chosen += 1
		}
	}
	response.ActivePlayer.score += total
	fmt.Printf("%s increases his score by %d\n", response.ActivePlayer.name, total)
	fmt.Printf("%s now has a total score of %d\n", response.ActivePlayer.name, response.ActivePlayer.score)
	fmt.Println()

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
