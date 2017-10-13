package sprite

import (
	"time"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

// FrameDesc describes a single frame of a sprite
type FrameDesc struct {
	Img *ebiten.Image
	// Points maps a name to a list of positions relative to the sprite's top left corner
	Points map[string][]geo.Vec
	// Rects maps a name to a list of rectangles relative to the sprite's top left corner
	Rects    map[string][]geo.Rect
	Duration time.Duration
}

// Desc holds information that is shared between all instances of the same sprite.
type Desc struct {
	Name   string
	Frames []FrameDesc
}

// Sprite is an instance of a sprite described by Desc. This allows one to create many
// copies of the same sprite, each animating independently, but all of the common data
// is shared.
type Sprite struct {
	*Desc
	Loop           bool
	frame          int
	untilNextFrame time.Duration
}

// Img returns the image for the current frame.
func (s *Sprite) Img() *ebiten.Image {
	return s.Frames[s.frame].Img
}

// Points returns the points associated with the given name for the current frame.
func (s *Sprite) Points(name string) []geo.Vec {
	// TODO: fill in points and rects on missing frames in Psd
	frame := s.frame
	points, ok := s.Frames[frame].Points[name]
	for !ok && frame > 0 {
		frame--
		points, ok = s.Frames[frame].Points[name]
	}
	return points
}

// Rects returns the rectangles associated with the given name for the current frame.
func (s *Sprite) Rects(name string) []geo.Rect {
	return s.Frames[s.frame].Rects[name]
}

// Start the sprite from the first frame.
func (s *Sprite) Start() {
	s.frame = 0
	s.untilNextFrame = s.Frames[s.frame].Duration
}

// Ended returns true when the sprite is not looping and is done with the last frame.
func (s *Sprite) Ended() bool {
	return !s.Loop && s.frame >= len(s.Frames)-1 && s.untilNextFrame <= 0
}

// Update subtracts dt from the time remaining on the current frame then advances the
// frame when it reaches 0. If Loop is true then the animation will start back from
// the beginning after reaching the end.
func (s *Sprite) Update(dt time.Duration) {
	s.untilNextFrame -= dt
	if s.Ended() {
		return
	}

	if s.untilNextFrame <= 0 {
		s.frame++
		if s.Loop && s.frame >= len(s.Frames) {
			s.frame = 0
		}
		s.untilNextFrame = s.Frames[s.frame].Duration
	}
}

// Draw the current frame to dst with the given options.
func (s *Sprite) Draw(dst *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if len(s.Frames) == 0 {
		return
	}

	dst.DrawImage(s.Frames[s.frame].Img, opts)
}
