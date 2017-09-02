package game

import (
	"image/color"

	"golang.org/x/image/font"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/game1/game/keymap"
	"github.com/Bredgren/game1/game/keymap/button"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

type keyLabel struct {
	action  keymap.Action
	bounds  geo.Rect
	face    font.Face
	img     map[bool]*ebiten.Image
	btnDown bool
}

func newKeyLabel(action keymap.Action, bounds geo.Rect, face font.Face) *keyLabel {
	// bounds, _ := font.BoundString(face, name)
	// width := (bounds.Max.X - bounds.Min.X).Ceil() + 4
	// height := (bounds.Max.Y - bounds.Min.Y).Ceil()
	// offset := (face.Metrics().Height - face.Metrics().Descent).Floor() - 1

	width := int(bounds.W)
	height := int(bounds.H)

	img1, _ := ebiten.NewImage(width, height, ebiten.FilterNearest)
	img1.Fill(color.RGBA{0, 0, 0, 50})
	// text.Draw(img1, name, face, 2, offset, color.Black)

	img2, _ := ebiten.NewImage(width, height, ebiten.FilterNearest)
	img2.Fill(color.RGBA{0, 0, 0, 100})
	// text.Draw(img2, name, face, 2, offset, color.White)

	k := &keyLabel{
		action: action,
		// bounds: geo.RectWH(geo.I2F2(width, height)),
		bounds: bounds,
		face:   face,
		img: map[bool]*ebiten.Image{
			false: img1,
			true:  img2,
		},
	}
	return k
}

func (k *keyLabel) draw(dst *ebiten.Image, cam *camera.Camera) {
	mouseOver := k.bounds.CollidePoint(geo.I2F2(ebiten.CursorPosition()))
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(k.bounds.TopLeft())
	dst.DrawImage(k.img[mouseOver], &opts)

	c := color.Black
	if k.btnDown {
		c = color.White
	}
	x, y := k.bounds.TopLeft()
	x += 4
	y += 14
	text.Draw(dst, string(k.action), k.face, int(x), int(y), c)
}

func (k *keyLabel) handleBtn(down bool) bool {
	k.btnDown = down
	return false
}

// func (k *keyLabel) handleAxis(val float64) bool {
// 	k.axisMove = val != 0
// 	return false
// }

type keyButton struct {
	btn     button.KeyMouse
	onClick func()
}

func (k *keyButton) draw(dst *ebiten.Image, cam *camera.Camera) {
	// opts := ebiten.DrawImageOptions{}
	// opts.GeoM.Translate(k.bounds.TopLeft())
	// dst.DrawImage(k.img[k.btnDown || k.axisMove], &opts)
}

type axisButton struct {
	axis int
}

type axisState struct {
}
