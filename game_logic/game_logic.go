package game_logic

import "math/rand/v2"

// a single instance of a game
type Game struct {
	Players []string
	Dice    []int
	// Card map is [player]Card
	Cards map[string]Card
}

type Dice struct {
	White1 int
	White2 int
	Red    int
	Yellow int
	Green  int
	Blue   int
}

type Row [11]bool

type Card struct {
	Player string
	Red    Row
	Yellow Row
	Green  Row
	Blue   Row
	Skips  int
}

// Attempts to mark row at index i. Returns t/f on succes/fail
func (r Row) TryMark(i int) Row {
	if i >= len(r)-1 {
		return r.tryRowLock()
	}
	for _, v := range r[i:] {
		if v {
			return r
		}
	}
	r[i] = true
	return r
}

// Attempts to lock row. Returns t/f on success/fail
func (r Row) tryRowLock() Row {
	var marked int
	for _, v := range r {
		if v {
			marked++
		}
	}
	if marked >= 5 {
		r[len(r)-1] = true
		return r
	}
	return r
}

// Randomizes all dice
func (d *Dice) Roll() {
	d.White1 = rand.IntN(6) + 1
	d.White2 = rand.IntN(6) + 1
	d.Red = rand.IntN(6) + 1
	d.Yellow = rand.IntN(6) + 1
	d.Green = rand.IntN(6) + 1
	d.Blue = rand.IntN(6) + 1
}
