package hestur

import (
	"math"
	"time"
)

// Character is a movable, AI-driven map
type Character struct {
	MapObject

	Type [4]byte
}

func (c Character) serialize(buf *serverBuffer, serType SerializationType) {
	c.MapObject.serialize(buf, serType)

	if serType == SerializationTypeFull {
		buf.Write(uint8(c.Type[0]))
		buf.Write(uint8(c.Type[1]))
		buf.Write(uint8(c.Type[2]))
		buf.Write(uint8(c.Type[3]))
	}
}

// NewDummyCharacter creates test character.
func NewDummyCharacter() *Character {
	c := Character{
		MapObject: NewMapObject(),
		Type:      [4]byte{'d', 'u', 'm', 'm'},
	}

	c.MapObject.UpdateMethod = func(dt float32) {
		c.updateDummy(dt)
	}

	c.StartUpdateLoop(1.0 / 30.0)

	return &c
}

func (c *Character) updateDummy(dt float32) {
	x := time.Now().Sub(c.MapObject.CreationTime).Seconds()
	c.Position.X = Coordinate(float64(c.Position.Y) + 15.0 +
		math.Sin(float64(c.Position.Y)*x*4/180.0*math.Pi)*15)
}
