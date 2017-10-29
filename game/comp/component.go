package comp

// Mask is a collection of components.
type Mask int64

const (
	// None is for no compnents.
	None Mask = 0
	// Position component is a 2D location in world space.
	Position Mask = 1 << iota
	// HitPoints destroys the entity when below 0.
	HitPoints
	// Velocity component gives an entity motion.
	Velocity
	// Rotation component allows and entity to rotate.
	Rotation
	// Gravity accelerates an entity vertically.
	Gravity
	// CollidesWithGround prevents an entity from moving bellow the ground.
	CollidesWithGround
	// BoundingBox specifies the bounds of the entity. It should be specified as if the
	// entity is at the origin.
	BoundingBox
	// Camera entity to use for drawing. Implies Position is world position
	Camera
	// Follow ties this entity's position to anothers
	Follow
	// Shake jitters an entity's position
	Shake
	// Sprite is the current image to use for drawing.
	Sprite
	// Animation updates the Sprite component to create an animation.
	Animation
	// Hitbox is a collection of rects that difine where the entity takes damage.
	Hitbox
	// Hurtbox is a collection of rect that define where the entity gives damage.
	Hurtbox
)

// Contains retursn true if m contains other.
func (m Mask) Contains(other Mask) bool {
	return m&other == other
}
