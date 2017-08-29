package game

import (
	"image/color"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/game1/game/keymap/button"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

type keyLabel struct {
	pos    geo.Vec
	img    map[bool]*ebiten.Image
	active bool
}

func newKeyLabel(pos geo.Vec, name string, face font.Face) *keyLabel {
	img1, _ := ebiten.NewImage(30, 10, ebiten.FilterNearest)
	img1.Fill(color.White)
	text.Draw(img1, name, face, 0, 10, color.Black)
	img2, _ := ebiten.NewImage(30, 10, ebiten.FilterNearest)
	img2.Fill(color.Black)
	text.Draw(img2, name, face, 0, 10, color.White)

	k := &keyLabel{
		pos: pos,
		img: map[bool]*ebiten.Image{
			false: img1,
			true:  img2,
		},
		active: false,
	}
	return k
}

func (k *keyLabel) draw(dst *ebiten.Image, cam *camera.Camera) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(k.pos.XY())
	dst.DrawImage(k.img[k.active], &opts)
}

type keyButton struct {
	btn button.Button
}

type axisButton struct {
	axis int
}

type axisState struct {
}
