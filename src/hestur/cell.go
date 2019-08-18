package hestur

// CellType is a type of cell.
type CellType int

const (
	// CellTypeDeepWater can't be passed through.
	CellTypeDeepWater CellType = iota

	// CellTypeEarth is a default land type.
	CellTypeEarth
)

// Cell contains all data related to specific map's cell.
type Cell struct {
	Type CellType
}

// CanPassThrough indicates can cell be passed through.
func (cell Cell) CanPassThrough() bool {
	switch cell.Type {
	case CellTypeDeepWater:
		return false
	case CellTypeEarth:
		return true
	default:
		return false
	}
}
