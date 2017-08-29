package game

import (
	"math"
	"time"

	"github.com/Bredgren/geo"
)

type playerCameraTarget struct {
	g      *Game
	p      *player
	offset geo.Vec
	pos    geo.Vec
}

func newPlayerCameraTarget(g *Game, p *player, screenHeight int) *playerCameraTarget {
	return &playerCameraTarget{
		g:      g,
		p:      p,
		offset: geo.VecXY(0, -float64(screenHeight)*0.4),
		pos:    p.Pos(),
	}
}

func (ct *playerCameraTarget) update(dt time.Duration) {
	switch ct.g.state {
	case mainMenuState:
		ct.pos.Y = ct.offset.Y
	case playState:
		offset := ct.offset
		offset.Y = -math.Max(0, ct.p.Pos().Y-offset.Y)
		ct.pos = ct.p.Pos().Plus(offset)
	}
}

func (ct *playerCameraTarget) Pos() geo.Vec {
	return ct.pos
}
