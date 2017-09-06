package sprite

import "github.com/hajimehoshi/ebiten"

type SheetDesc map[string]Desc
type Desc []FrameDesc

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

func NewSheet(img *ebiten.Image, layout SheetDesc) {
}
