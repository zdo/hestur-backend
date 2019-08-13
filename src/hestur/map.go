package hestur

type Map struct {
	Width, Height int
	Cells         []Cell
	// Characters    []Character
}

func CreateEmptyMap(width, height int) Map {
	m := Map{}
	m.Cells = make([]Cell, width*height)
	m.Width = width
	m.Height = height

	for i := 0; i < width; i++ {
		x := i
		y := i
		index := x + y*width
		m.Cells[index] = Cell{Type: Grass}
	}

	return m
}