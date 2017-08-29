package game

import (
	"image/color"

	"golang.org/x/image/font"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/game1/game/keymap/button"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

type keyLabel struct {
	name   string
	bounds geo.Rect
	img    map[bool]*ebiten.Image
	active bool
}

func newKeyLabel(name string, face font.Face) *keyLabel {
	bounds, _ := font.BoundString(face, name)
	width := (bounds.Max.X - bounds.Min.X).Ceil() + 4
	height := (bounds.Max.Y - bounds.Min.Y).Ceil()
	offset := (face.Metrics().Height - face.Metrics().Descent).Floor() - 1

	img1, _ := ebiten.NewImage(width, height, ebiten.FilterNearest)
	img1.Fill(color.RGBA{0, 0, 0, 50})
	text.Draw(img1, name, face, 2, offset, color.Black)

	img2, _ := ebiten.NewImage(width, height, ebiten.FilterNearest)
	img2.Fill(color.RGBA{0, 0, 0, 150})
	text.Draw(img2, name, face, 2, offset, color.White)

	k := &keyLabel{
		name:   name,
		bounds: geo.RectWH(geo.I2F2(width, height)),
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
	opts.GeoM.Translate(k.bounds.TopLeft())
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
