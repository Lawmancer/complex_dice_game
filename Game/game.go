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

// Play a game
type Play struct {
	round   int
	players []player
	action  chan RollResult
	Choices chan RollResult
	done    bool
}

// Start playing the game
func (p *Play) Start() {
	if p.players == nil || len(p.players) == 0 {
		fmt.Println("No players, halting.")
		return
	}

	fmt.Printf("Starting a new game with %d players…\n", len(p.players))
	p.randomFirstPlayer()
	p.listPlayers()
	for p.round < rounds {
		p.round++
		fmt.Println("-----------------------")
		fmt.Println("Starting Round", p.round)
		for i := range p.players {
			p.takeTurn(i)
		}
		fmt.Println()
		p.players[0].turnOrder = p.players[0].turnOrder + len(p.players) // move to end of line
		p.resortPlayers()
	}

	p.done = true
	if p.action != nil {
		close(p.action)
	}

	winner := p.getWinner()
	for _, gamer := range p.players {
		fmt.Printf("%s has a final score of %d\n", gamer.name, gamer.score)
	}
	fmt.Printf("\n%s wins with a final score of %d\n", winner.name, winner.score)

	return
}

// Register a new player before the game starts
func (p *Play) Register(name string) (id int, err error) {
	if p.round > 0 {
		return 0, errors.New("the game has already started")
	}

	id = len(p.players) + 1
	newPlayer := player{
		id,
		name,
		0,
		0,
	}
	p.players = append(p.players, newPlayer)

	return id, nil
}

// WaitForTurn is used outside the game to wait for a Turn to be playable
func (p *Play) WaitForTurn() (t RollResult, done bool) {
	p.action = make(chan RollResult)
	for turn := range p.action {
		t = turn
		close(p.action)
	}
	if p.done {
		done = true
	}

	return
}

func (p *Play) isDone() bool {
	return p.done
}

func (p *Play) takeTurn(playerNum int) {
	fmt.Printf("\n%s is taking their turn.\n", p.players[playerNum].name)
	rand.Seed(time.Now().UnixNano()) // seed once per turn

	diceRemaining := totalDice
	gamer := &p.players[playerNum]
	for diceRemaining > 0 {
		var rolls []die
		t := RollResult{
			p.round,
			gamer,
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
		p.action <- t
		p.Choices = make(chan RollResult)
		chosen := 0
		for choices := range p.Choices {
			chosen = p.processChoices(choices)
			close(p.Choices)
		}
		diceRemaining -= chosen
	}

	return
}

func (p *Play) processChoices(response RollResult) (chosen int) {
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

func (p *Play) getWinner() player {
	sort.Slice(p.players, func(i, j int) bool {
		return p.players[i].score < p.players[j].score
	})

	return p.players[0]
}

func (p *Play) listPlayers() {
	for i, gamer := range p.players {
		fmt.Printf("Player %d: %s\n", i+1, gamer.name)
	}
	fmt.Println()
}

func (p *Play) randomFirstPlayer() {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(len(p.players))
	p.players[r].turnOrder = 1
	next := 2
	for i := 0; i < len(p.players); i++ {
		if p.players[i].turnOrder == 0 {
			p.players[i].turnOrder = next
			next++
		}
	}
	p.resortPlayers()
}

func (p *Play) resortPlayers() {
	sort.Slice(p.players, func(i, j int) bool {
		return p.players[i].turnOrder < p.players[j].turnOrder
	})
}
