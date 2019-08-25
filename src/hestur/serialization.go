package hestur

// SerializationType is type of object's serialization.
type SerializationType int

const (
	// SerializationTypeFull means write all object's data.
	SerializationTypeFull SerializationType = iota

	// SerializationTypeShort means write only the frequently changed data.
	SerializationTypeShort
)
