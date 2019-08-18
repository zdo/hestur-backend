package hestur

// CharacterType indicates type of the character.
type CharacterType int

const (
	// CharacterTypeDummy is test dummy character.
	CharacterTypeDummy CharacterType = iota
)

// Character is a movable, AI-driven map
type Character struct {
	MapObject

	Type CharacterType
}
