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

func (t *Turn) nextTurn(players int) {
	fmt.Println("\nnextturn called\n")
	*t++
	if int(*t) > players-1 {
		*t = 0
	}
}

// Attempts to mark row at index i. Returns t/f on succes/fail
func (r Row) TryMark(i int, turn bool, dice [3]Die) (Row, bool, bool) {
	cmp := i + 2
	if turn && (cmp == int(dice[0]+dice[2]) || cmp == int(dice[1]+dice[2])) {
		for _, v := range r[i:] {
			if v {
				return r, false, false
			}
			r[i] = true
			return r, true, true
		}
	}
	if cmp == int(dice[0]+dice[1]) {
		for _, v := range r[i:] {
			if v {
				return r, false, false
			}
		}
		r[i] = true
		return r, true, false
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
}

// T if everyone has gone, F otherwise
func (g *Game) TurnCheck() bool {
	for _, v := range g.Gone {
		if !v {
			fmt.Println("returning false")
			return false
		}
	}
	for _, v := range g.Gone2 {
		if v {
			g.Turn.nextTurn(len(g.Players))
			fmt.Println("returning true")
			fmt.Printf("\ngone1: %+v\ngone2: %+v\n", g.Gone, g.Gone2)
			return true
		}
	}
	fmt.Println("returning false")
	return false
}
