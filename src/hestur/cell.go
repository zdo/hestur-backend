package hestur

type CellType int

const (
	Earth CellType = iota
	Grass
)

type Cell struct {
	Type CellType
}