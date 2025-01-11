package game_logic

type Player string

// a single instance of a game
type Game struct {
	Players []Player
	Dice    []int
	// Card map is [player]Card
	Cards map[Player]Card
}

type Card struct {
	Player string
	Red    map[int]bool
	Yellow map[int]bool
	Green  map[int]bool
	Blue   map[int]bool
	Skips  int
}

func (c *Card) RowLock(index int) (e error) {
	return nil
}
