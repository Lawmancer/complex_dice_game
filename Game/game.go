package Game

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// Wild counts as 0
const Wild = 4
const wildValue = 0
const rounds = 4
const totalDice = 5

// Play implements GameInterface // TODO correct comment
type Play struct {
	round   int
	players []player
	action  chan RollResult
	Choices chan RollResult
	Done    bool
}

// Start playing the game
func (g *Play) Start() {
	if g.players == nil || len(g.players) == 0 {
		fmt.Println("No players, halting.")
		return
	}

	fmt.Printf("Starting a new game with %d players…\n", len(g.players))
	g.randomFirstPlayer()
	g.listPlayers()
	for g.round < rounds {
		g.round++
		fmt.Println("-----------------------")
		fmt.Println("Starting Round", g.round)
		for i := range g.players {
			g.takeTurn(i)
		}
		fmt.Println()
		g.players[0].turnOrder = g.players[0].turnOrder + len(g.players) // move to end of line
		g.resortPlayers()
	}

	g.Done = true
	if g.action != nil {
		close(g.action)
	}

	winner := g.getWinner()
	for _, p := range g.players {
		fmt.Printf("%s has a final score of %d\n", p.name, p.score)
	}
	fmt.Printf("\n%s wins with a final score of %d\n", winner.name, winner.score)

	return
}

// Register a new player before the game starts
func (g *Play) Register(name string) (id int, err error) {
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

func (g *Play) takeTurn(playerNum int) {
	fmt.Printf("\n%s is taking their turn.\n", g.players[playerNum].name)
	rand.Seed(time.Now().UnixNano()) // seed once per turn

	diceRemaining := totalDice
	p := &g.players[playerNum]
	for diceRemaining > 0 {
		var rolls []die
		t := RollResult{
			g.round,
			p,
			rolls,
		}
		d := diceRemaining
		for d > 0 {
			rnd := rand.Intn(5) + 1
			var r die
			r.Pips.value = rnd
			t.Dice = append(t.Dice, r)
			d--
		}
		g.action <- t
		g.Choices = make(chan RollResult)
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
func (g *Play) WaitForTurn() (t RollResult, done bool) {
	g.action = make(chan RollResult)
	for turn := range g.action {
		t = turn
		close(g.action)
	}
	if g.Done {
		done = true
	}

	return
}

func (g *Play) processChoices(response RollResult) (chosen int) {
	total := 0
	for _, r := range response.Dice {
		if r.Selected {
			if r.Pips.value == Wild {
				total += wildValue
			} else {
				total += r.Pips.value
			}
			chosen += 1
		}
	}
	response.ActivePlayer.score += total
	fmt.Printf("  … score goes up %d (total: %d)\n", total, response.ActivePlayer.score)

	return
}

func (g *Play) getWinner() player {
	sort.Slice(g.players, func(i, j int) bool {
		return g.players[i].score < g.players[j].score
	})

	return g.players[0]
}

func (g *Play) listPlayers() {
	for i, p := range g.players {
		fmt.Printf("Player %d: %s\n", i+1, p.name)
	}
	fmt.Println()
}

func (g *Play) randomFirstPlayer() {
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

func (g *Play) resortPlayers() {
	sort.Slice(g.players, func(i, j int) bool {
		return g.players[i].turnOrder < g.players[j].turnOrder
	})
}
