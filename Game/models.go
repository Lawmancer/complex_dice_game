package Game

type player struct {
	Id        int
	name      string
	turnOrder int
	score     int
}

// RollResult has all the information a
// player needs to make their choice
type RollResult struct {
	Round        int
	ActivePlayer *player
	Dice         []die
}

type die struct {
	Pips     readOnlyInt
	Selected bool
}

type readOnlyInt struct {
	value int
}

// Value gives access the read only value of a roll
func (s readOnlyInt) Value() int {
	return s.value
}
