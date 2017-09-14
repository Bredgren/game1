package sprite

import (
	"encoding/json"
	"fmt"
	"image"
	"time"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

// SheetDesc maps from sprite name to a description of that sprite. Normally one may want
// to parse this from a json file.
type SheetDesc map[string]Desc

// Desc describes a single sprite by its indvidual frames and duration.
type Desc struct {
	Frames   []FrameDesc `json:"frames"`
	Duration Duration    `json:"duration"`
}

// FrameDesc describes the position and size of a single frame. When drawing a sprite the
// anchor for each frame will be put at the position the sprite is drawn at. The Weight
// field can be used to give a larger slice of the sprite's duration to individual frames.
// If the weight of all frames is the same then the sprite duration will be divided evenly
// between all frames. A frame whose weight is 2x the weight of another will be shown for
// twice as long.
type FrameDesc struct {
	X      int     `json:"x"`
	Y      int     `json:"y"`
	W      int     `json:"w"`
	H      int     `json:"h"`
	Weight float64 `json:"weight"`
	Anchor struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"anchor"`
}

// AddSheet extracts sprites from the sheet according to the given layout. Sprites may
// be retrieved afterward via the Get function. Sprite names (the keys of SheetDesc) must
// be unique among all sheets added with this function. A non-nil error is returned if
// this is not the case.
func AddSheet(img *ebiten.Image, layout SheetDesc) error {
	for spriteName, spriteDesc := range layout {
		if _, exists := sprites[spriteName]; exists {
			return fmt.Errorf("sprite with name '%s' already exists", spriteName)
		}

		s := Sprite{
			src:      img,
			frames:   make([]frame, len(spriteDesc.Frames)),
			duration: spriteDesc.Duration.Duration,
		}
		totalWeight := 0.0
		for i, frameDesc := range spriteDesc.Frames {
			xMax := frameDesc.X + frameDesc.W
			yMax := frameDesc.Y + frameDesc.H
			rect := image.Rect(frameDesc.X, frameDesc.Y, xMax, yMax)
			s.frames[i] = frame{
				opts: ebiten.DrawImageOptions{
					SourceRect: &rect,
				},
				anchor: geo.VecXYi(frameDesc.Anchor.X, frameDesc.Anchor.Y),
			}

			totalWeight += frameDesc.Weight
		}

		if totalWeight == 0 {
			totalWeight = 1
		}

		cumTime := 0.0
		for i, frameDesc := range spriteDesc.Frames {
			normalizedWeight := frameDesc.Weight / totalWeight
			cumTime += float64(s.duration.Nanoseconds()) * normalizedWeight
			s.frames[i].endTime = time.Duration(cumTime) * time.Nanosecond
		}

		sprites[spriteName] = s
	}
	return nil
}

// Duration is just time.Duration that can be used in json.
type Duration struct {
	time.Duration
}

// UnmarshalJSON parses a duration from a string, or if it's an int then assumes it's
// nanoseconds.
func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' {
		sd := string(b[1 : len(b)-1])
		d.Duration, err = time.ParseDuration(sd)
		return
	}

	var id int64
	id, err = json.Number(string(b)).Int64()
	d.Duration = time.Duration(id)

	return
}

// MarshalJSON formats for json.
func (d Duration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}
