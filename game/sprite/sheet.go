package sprite

import "github.com/hajimehoshi/ebiten"

// SheetDesc maps from sprite name to a description of that sprite.
type SheetDesc map[string]Desc

// Desc is a list of frames that make up a sprite.
type Desc []FrameDesc

// FrameDesc describes the position and size of a single frame. When drawing a sprite the
// anchor for each frame will be put at the position the sprite is drawn at.
type FrameDesc struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	W      int `json:"w"`
	H      int `json:"h"`
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
	return nil
}
