package comp

// Mask is a collection of components.
type Mask int64

const (
	// None is for no compnents.
	None Mask = 0
	// Position component is a 2D location in world space.
	Position Mask = 1 << iota
	// Velocity component gives an entity motion.
	Velocity
	// Rotation component allows and entity to rotate.
	Rotation
	// Gravity accelerates an entity vertically.
	Gravity
	// BoundingBox specifies the bounds of the entity. It should be specified as if the
	// entity is at the origin.
	BoundingBox
	// Camera entity to use for drawing. Implies Position is world position
	Camera
	// Follow ties this entity's position to anothers
	Follow
	// Shake jitters an entity's position
	Shake
)

// Contains retursn true if m contains other.
func (m Mask) Contains(other Mask) bool {
	return m&other == other
}
