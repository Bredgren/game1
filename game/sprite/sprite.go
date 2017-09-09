package sprite

import (
	"time"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

var sprites = map[string]Sprite{}

// Get returns the sprite with the give name that was added via Add or AddSheet. Returns
// an empty Sprite if one with the name hasn't been added.
func Get(name string) Sprite {
	return sprites[name]
}

type frame struct {
	opts    ebiten.DrawImageOptions
	anchor  geo.Vec
	endTime time.Duration
}

// Sprite is an image or collection of images that can play in succession.
type Sprite struct {
	src      *ebiten.Image
	frames   []frame
	duration time.Duration
	timeLeft time.Duration
	curFrame int
	loop     bool
}

// Begin starts the sprites animation if it has more than one frame. If loop is true then
// the animation will replay from the beginning each time that it reaches the end.
func (s *Sprite) Begin(loop bool) {
	s.timeLeft = s.duration
	s.curFrame = 0
	s.loop = loop
}

// Ended returns true if it has reached the end of its animation and isn't looping.
func (s *Sprite) Ended() bool {
	return s.timeLeft <= 0
}

// Update advances the sprite by the given time, changing frames when needed.
func (s *Sprite) Update(dt time.Duration) {
	s.timeLeft -= dt
	if s.timeLeft <= 0 && s.loop {
		s.Begin(true)
	}

	if s.duration-s.timeLeft > s.frames[s.curFrame].endTime {
		s.curFrame++
		if s.curFrame >= len(s.frames) {
			s.curFrame = len(s.frames) - 1
		}
	}
}

// Draw draws the spite's current frame to dst at the given position. The anchor position
// of the current frame will be placed at pos. The geom parameter can be used to apply
// extra transformations to the sprite before drawing it.
func (s *Sprite) Draw(dst *ebiten.Image, pos geo.Vec, geom ebiten.GeoM) {
	if len(s.frames) == 0 {
		return
	}

	pos.Sub(s.frames[s.curFrame].anchor)

	s.frames[s.curFrame].opts.GeoM.Reset()
	s.frames[s.curFrame].opts.GeoM.Concat(geom)
	s.frames[s.curFrame].opts.GeoM.Translate(pos.XY())
	dst.DrawImage(s.src, &s.frames[s.curFrame].opts)
}

// Size returns the width and height of the current frame as a vector.
func (s *Sprite) Size() geo.Vec {
	return geo.VecPoint(s.frames[s.curFrame].opts.SourceRect.Size())
}
