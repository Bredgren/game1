package camera

import (
	"time"

	"github.com/Bredgren/geo"
)

// Target is any object that has a function for returning a position that the Camera
// can use as a target position.
type Target interface {
	Pos() geo.Vec
}

// Camera manages a 2-D camera. If its distance from Target is greater than MaxDist when
// Update is called then its position is moved directly toward the target so that the
// distance is equal to MaxDist. If the Camera is less than MaxDist then its velocity
// is directly toward the Target with a magnitude that is a percentage of MaxSpeed determined
// by the Ease function. The Ease function is given the ratio <distance to Target>/MaxDist
// and the return value is multiplied by MaxSpeed.
type Camera struct {
	pos      geo.Vec
	offset   geo.Vec
	halfSize geo.Vec
	Target   Target
	MaxDist  float64
	MaxSpeed float64
	Ease     geo.EaseFn
	// Shaker is optional. Set its fields and call the Camera's StartShake functions. If
	// Shaker.Falloff is nil then the Shaker's ShakeConst is used.
	Shaker     geo.Shaker
	shakerTime time.Time
}

// New creates, initializes, and returns a new Camera. The parameters width and height
// are the dimensions of the image tha Camera will be used to draw to. The default Ease
// fuction is linear, MaxSpeed=0, and MaxDist=0. This results in perfectly sticking
// to the Target, though MaxDist=0 on its own is sufficient for that behavior.
func New(width, height int) *Camera {
	return &Camera{
		halfSize:   geo.VecXYi(width/2, height/2),
		Ease:       geo.EaseLinear,
		shakerTime: time.Now(),
	}
}

// Update updates the Camera's state simulating dt time passed.
func (c *Camera) Update(dt time.Duration) {
	c.shakerTime = c.shakerTime.Add(dt)

	target := c.Target.Pos()
	distToTarget2 := target.Dist2(c.pos)
	max2 := c.MaxDist * c.MaxDist
	if distToTarget2 > max2 {
		c.pos = target.Plus(c.pos.Minus(target).WithLen(c.MaxDist))
	} else {
		ratio := distToTarget2 / max2
		speed := c.Ease(ratio) * c.MaxSpeed
		vel := target.Minus(c.pos).WithLen(speed)
		c.pos.Add(vel.Times(dt.Seconds()))
	}

	if c.Shaker.Falloff != nil {
		c.offset = c.Shaker.Shake(c.shakerTime)
	} else {
		c.offset = c.Shaker.ShakeConst(c.shakerTime)
	}
}

// ScreenCoords takes a position in world coordinates and returns its position on the screen.
func (c *Camera) ScreenCoords(pos geo.Vec) geo.Vec {
	return pos.Minus(c.topLeft())
}

// WorldCoords takes a position in screen coordinates and returns its position in the world.
func (c *Camera) WorldCoords(pos geo.Vec) geo.Vec {
	return pos.Plus(c.topLeft())
}

func (c *Camera) topLeft() geo.Vec {
	return c.Center().Minus(c.halfSize)
}

// Center returns the camera's center position in world coordinates.
func (c *Camera) Center() geo.Vec {
	cameraCenter := c.pos.Plus(c.offset)
	cameraCenter.Floor()
	return cameraCenter
}

// StartShake restarts the Shaker time.
func (c *Camera) StartShake() {
	c.Shaker.StartTime = c.shakerTime
}
