package game_logic

import (
	"fmt"
	"math/rand/v2"
)

// a single instance of a game
type Game struct {
	Key     string
	Players []string
	Dice
	Turn
	// Card map is [player]Card
	Cards map[string]Card
}

type Dice struct {
	White1 Die
	White2 Die
	Red    Die
	Yellow Die
	Green  Die
	Blue   Die
}

type Die int

type Row [11]bool

type Card struct {
	Player string
	Red    Row
	Yellow Row
	Green  Row
	Blue   Row
	Skips  int
}

type Turn int

func (t *Turn) NextTurn(players int) {
	tmp := *t
	if int(tmp) > players {
		tmp = 0
	}
	t = &tmp
}

// Attempts to mark row at index i. Returns t/f on succes/fail
func (r Row) TryMark(i int, turn bool, dice [3]Die) (Row, bool) {
	cmp := i + 2
	if turn && (cmp == int(dice[0]+dice[2]) || cmp == int(dice[1]+dice[2])) {
		for _, v := range r[i:] {
			if v {
				return r, false
			}
			r[i] = true
			return r, true
		}
	}
	if cmp == int(dice[0]+dice[1]) {
		for _, v := range r[i:] {
			if v {
				fmt.Println("2")
				return r, false
			}
		}
		r[i] = true
		return r, true
	}
	return r, false
}

// Attempts to lock row. Returns t/f on success/fail
func (r Row) tryRowLock() (Row, bool) {
	var marked int
	for _, v := range r {
		if v {
			marked++
		}
	}
	if marked >= 5 {
		r[len(r)-1] = true
		return r, true
	}
	return r, false
}

// Randomizes all dice
func (d *Dice) Roll() {
	d.White1 = Die(rand.IntN(6) + 1)
	d.White2 = Die(rand.IntN(6) + 1)
	d.Red = Die(rand.IntN(6) + 1)
	d.Yellow = Die(rand.IntN(6) + 1)
	d.Green = Die(rand.IntN(6) + 1)
	d.Blue = Die(rand.IntN(6) + 1)
}
