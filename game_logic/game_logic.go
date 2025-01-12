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
	Red    [11]bool
	Yellow [11]bool
	Green  [11]bool
	Blue   [11]bool
	Skips  int
}

func (c *Card) RowLock(index int) (e error) {
	return nil
}
