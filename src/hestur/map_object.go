package hestur

// Coordinate is used for indicating position on the axis.
type Coordinate float32

// Position is used for indicating 2D-position on the map.
type Position struct {
	X, Y Coordinate
}

// MapObject is a basic struct for any object located on the map.
type MapObject struct {
	Name     string
	Position Position
}
