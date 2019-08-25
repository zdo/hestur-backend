package hestur

import (
	"sync"
)

// Game is a root object
type Game struct {
	Width, Height int
	Cells         []Cell

	shouldStop     bool
	Dt             float64
	TimeSinceStart float64

	mapObjectsLock *sync.RWMutex
	Characters     []*Character
	props          []*Prop

	nextMapObjectID     MapObjectID
	nextMapObjectIDLock *sync.Mutex
}

// NewGame creates new test game.
func NewGame() Game {
	game := Game{}
	game.Width = 100
	game.Height = 100

	game.shouldStop = false

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

	game.nextMapObjectIDLock = new(sync.Mutex)
	game.mapObjectsLock = new(sync.RWMutex)
	game.Characters = make([]*Character, 0)
	game.props = make([]*Prop, 0)

	return game
}

func (game *Game) generateMapObjectID() MapObjectID {
	game.nextMapObjectIDLock.Lock()
	nextID := game.nextMapObjectID
	game.nextMapObjectID++
	game.nextMapObjectIDLock.Unlock()

	return nextID
}

// RegisterCharacter adds character to the internal processing list.
func (game *Game) RegisterCharacter(c *Character) {
	game.mapObjectsLock.Lock()
	c.MapObject.ID = game.generateMapObjectID()
	game.Characters = append(game.Characters, c)
	game.mapObjectsLock.Unlock()
}
