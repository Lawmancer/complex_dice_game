package Dice

type player struct {
	id        int
	name      string
	turnOrder int
	score     int
}

type turn struct {
	round        int
	activePlayer player
	rolls        []int
}
