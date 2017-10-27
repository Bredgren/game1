package game

import (
	"github.com/Bredgren/geo"
)

type fixedCameraTarget struct {
	pos geo.Vec
}

func (ct fixedCameraTarget) Pos() geo.Vec {
	return ct.pos
}

// type dynamicCameraTarget struct {
// 	p      *player
// 	offset geo.Vec
// 	pos    geo.Vec
// }
//
// func newDynamicCameraTarget(p *player, screenHeight int) *dynamicCameraTarget {
// 	return &dynamicCameraTarget{
// 		p:      p,
// 		offset: geo.VecXY(0, -float64(screenHeight)*0.4),
// 		pos:    p.Pos(),
// 	}
// }
//
// func (ct *dynamicCameraTarget) update(dt time.Duration) {
// 	offset := ct.offset
// 	offset.Y = -math.Max(0, ct.p.Pos().Y-offset.Y)
// 	ct.pos = ct.p.Pos().Plus(offset)
// }
//
// func (ct *dynamicCameraTarget) Pos() geo.Vec {
// 	return ct.pos
// }
