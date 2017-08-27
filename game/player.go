package game

import (
	"image/color"
	"time"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

const (
	playerMoveSpeed = 500
)

type player struct {
	pos    geo.Vec
	bounds geo.Rect
	img    *ebiten.Image
	left   bool
	right  bool
	move   float64
}

func newPlayer() *player {
	img, _ := ebiten.NewImage(16, 16, ebiten.FilterNearest)
	img.Fill(color.Black)
	p := &player{
		bounds: geo.RectWH(16, 16),
		img:    img,
	}
	return p
}

func (p *player) update(dt time.Duration) {
	if p.move == 0 {
		if p.left {
			p.move = -1
		}
		if p.right {
			p.move = 1
		}
	}
	p.pos.Add(geo.VecXY(p.move, 0).Times(playerMoveSpeed * dt.Seconds()))
	p.move = 0
	p.bounds.SetBottomMid(p.pos.XY())
}

func (p *player) draw(dst *ebiten.Image, cam *camera.Camera) {
	pos := cam.ScreenCoords(geo.VecXY(p.bounds.TopLeft()))
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(pos.XY())
	dst.DrawImage(p.img, &opts)
}

func (p *player) Pos() geo.Vec {
	return p.pos
}
