package hestur

// Game is a root object
type Game struct {
	Width, Height int
	Cells         []Cell

	Characters []Character
	Props      []Prop
}

// NewGame creates new test game.
func NewGame() Game {
	game := Game{}
	game.Width = 100
	game.Height = 100

	cellsCount := game.Width * game.Height
	game.Cells = make([]Cell, cellsCount)

	waterBorder := 3

	for y := 0; y < game.Height; y++ {
		for x := 0; x < game.Width; x++ {
			index := x + y*game.Width
			isWater := (x < waterBorder || x >= (game.Width-waterBorder) ||
				y < waterBorder || y >= (game.Height-waterBorder))
			if isWater {
				game.Cells[index].Type = CellTypeDeepWater
			} else {
				game.Cells[index].Type = CellTypeEarth
			}
		}
	}

	game.Characters = make([]Character, 0)
	game.Props = make([]Prop, 0)

	return game
}
