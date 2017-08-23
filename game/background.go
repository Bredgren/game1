package game

import (
	"image/color"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/game1/game/util"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

type background struct {
	color1    color.Color
	color2    color.Color
	maxHeight float64
}

func NewBackground() *background {
	return &background{
		color1:    color.NRGBA{255, 140, 68, 255},
		color2:    color.NRGBA{0, 0, 10, 255},
		maxHeight: 500,
	}
}

func (b *background) Draw(dst *ebiten.Image, cam *camera.Camera) {
	height := geo.Clamp(-cam.Center().Y, 0, b.maxHeight) / b.maxHeight
	dst.Fill(util.LerpColor(b.color1, b.color2, height))
}
