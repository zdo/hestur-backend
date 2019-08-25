package hestur

import (
	"time"
)

// MapObjectID is ID of map object.
type MapObjectID uint32

// Coordinate is used for indicating position on the axis.
type Coordinate float32

// Position is used for indicating 2D-position on the map.
type Position struct {
	X, Y Coordinate
}

// MapObject is a basic struct for any object located on the map.
type MapObject struct {
	Game     *Game
	ID       MapObjectID
	Name     string
	Position Position

	isDied       bool
	CreationTime time.Time
	UpdateMethod func(float32)
}

// NewMapObject creates stub map object.
func NewMapObject() MapObject {
	mo := MapObject{
		Game:         nil,
		ID:           0,
		Name:         "Anonymous",
		Position:     Position{X: 0, Y: 0},
		CreationTime: time.Now(),
	}
	return mo
}

func (mo MapObject) serialize(buf *serverBuffer, serType SerializationType) {
	buf.Write(uint32(mo.ID))

	if serType == SerializationTypeFull {
		buf.Write(uint8(0xff))
	} else {
		buf.Write(uint8(0x00))
	}

	if mo.isDied {
		buf.Write(uint8(1))
	} else {
		buf.Write(uint8(0))
	}

	if serType == SerializationTypeFull {
		buf.WriteString(mo.Name)
	}

	buf.Write(float32(mo.Position.X))
	buf.Write(float32(mo.Position.Y))
}

// StartUpdateLoop launches update loop for that map object.
func (mo *MapObject) StartUpdateLoop(wantedFramesPerSec float32) {
	go mo.updateLoop(wantedFramesPerSec)
}

func (mo *MapObject) updateLoop(wantedFramesPerSec float32) {
	prevTime := time.Now()

	for !mo.isDied {
		time.Sleep(time.Duration(wantedFramesPerSec*1000) * time.Millisecond)
		now := time.Now()
		dt := float32(now.Sub(prevTime).Seconds())
		mo.UpdateMethod(dt)
		prevTime = now
	}
}
