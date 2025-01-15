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
	Gone  map[string]bool
	Gone2 map[string]bool
	Skips map[string]int
	// Keys should be red, yellow, green, blue
	RowsLocked map[string]bool
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
	TurnOver bool
	Player   string
	Red      Row
	Yellow   Row
	Green    Row
	Blue     Row
	Skips    int
}

type Turn int

func (t *Turn) nextTurn(players int) {
	*t++
	if int(*t) > players-1 {
		*t = 0
	}
}

// Attempts to mark row at index i. Returns t/f on succes/fail
func (r Row) TryMark(i int, turn bool, dice [3]Die, large bool, turngone bool) (Row, bool, bool) {
	fmt.Printf("\ntg: %t\n", turngone)
	if !large {
		cmp := i + 2
		if cmp == int(dice[0]+dice[1]) && !turngone {
			if cmp == 12 {
				tmp1, tmp2 := r.tryRowLock()
				return tmp1, tmp2, false
			}
			for _, v := range r[i:] {
				if v {
					return r, false, false
				}
			}
			r[i] = true
			return r, true, false
		}
		if turn && (cmp == int(dice[0]+dice[2]) || cmp == int(dice[1]+dice[2])) {
			if cmp == 12 {
				tmp1, tmp2 := r.tryRowLock()
				return tmp1, tmp2, tmp2
			}
			for _, v := range r[i:] {
				if v {
					return r, false, false
				}
			}
			r[i] = true
			return r, true, true
		}
		return r, false, false
	}

	cmp := 12 - i
	if cmp == int(dice[0]+dice[1]) && !turngone {
		if cmp == 2 {
			tmp1, tmp2 := r.tryRowLock()
			return tmp1, tmp2, false
		}
		for _, v := range r[i:] {
			if v {
				return r, false, false
			}
		}
		r[i] = true
		return r, true, false
	}
	if turn && (cmp == int(dice[0]+dice[2]) || cmp == int(dice[1]+dice[2])) {
		if cmp == 2 {
			tmp1, tmp2 := r.tryRowLock()
			return tmp1, tmp2, tmp2
		}
		for _, v := range r[i:] {
			if v {
				return r, false, false
			}
			r[i] = true
			return r, true, true
		}
	}
	return r, false, false
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

func (g *Game) ResetTurns() {
	b := make(map[string]bool)
	c := make(map[string]bool)
	for _, v := range g.Players {
		b[v] = false
	}
	for _, v := range g.Players {
		c[v] = false
	}
	g.Gone = b
	g.Gone2 = c

	for k, v := range g.Cards {
		v.TurnOver = false
		g.Cards[k] = v
	}
	fmt.Printf("\ncards: %+v\n", g.Cards)
}

// T if everyone has gone, F otherwise
func (g *Game) TurnCheck() bool {
	for _, v := range g.Gone {
		if !v {
			return false
		}
	}
	for _, v := range g.Gone2 {
		if v {
			g.Turn.nextTurn(len(g.Players))
			return true
		}
	}
	fmt.Println("returning false")
	return false
}
