package hestur

// PropType indicates type of the prop.
type PropType int

const (
	// PropTypeDummy is test dummy prop.
	PropTypeDummy PropType = iota
)

// Prop is a static map object without AI.
type Prop struct {
	MapObject

	Type PropType
}
