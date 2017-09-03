package game

import (
	"time"

	"github.com/Bredgren/game1/game/camera"
	"github.com/hajimehoshi/ebiten"
)

type playState struct {
	p            *player
	cam          *camera.Camera
	target       *dynamicCameraTarget
	bg           *background
	screenHeight int
}

func newPlayState(p *player, screenHeight int, cam *camera.Camera, bg *background) *playState {
	return &playState{
		p:            p,
		cam:          cam,
		target:       newDynamicCameraTarget(p, screenHeight),
		bg:           bg,
		screenHeight: screenHeight,
	}
}

func (p *playState) begin(previousState gameStateName) {
	p.cam.Target = p.target
}

func (p *playState) end() {
}

func (p *playState) nextState() gameStateName {
	return play
}

func (p *playState) update(dt time.Duration) {
	p.p.update(dt)
	p.target.update(dt)
}

func (p *playState) draw(dst *ebiten.Image, cam *camera.Camera) {
	p.bg.Draw(dst, cam)
	p.p.draw(dst, cam)
}
