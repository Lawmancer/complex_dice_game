package Dice

type player struct {
	Id        int
	name      string
	turnOrder int
	score     int
}

// Turn has all the information a
// player needs to make their choice
type Turn struct {
	Round        int
	ActivePlayer *player
	Rolls        []roll
}

type roll struct {
	Roll     readOnlyInt
	Selected bool
}

type readOnlyInt struct {
	value int
}

// Value gives access the read only value of a roll
func (s readOnlyInt) Value() int {
	return s.value
}
