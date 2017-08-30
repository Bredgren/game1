package game

import (
	"time"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

type introState struct {
	p    *player
	bg   *background
	wait time.Duration
}

func newIntroState(p *player, screenHeight int, cam *camera.Camera, bg *background) *introState {
	p.SetPos(geo.Vec0)
	cam.Target = fixedCameraTarget{geo.VecXY(0, -float64(screenHeight)*0.4)}
	return &introState{
		p:    p,
		bg:   bg,
		wait: 3 * time.Second,
	}
}

func (i *introState) begin(previousState gameStateName) {
	// Since introState is the first state, begin is not called
}

func (i *introState) end() {

}

func (i *introState) nextState() gameStateName {
	if i.wait <= 0 {
		return mainMenu
	}
	return intro
}

func (i *introState) update(dt time.Duration) {
	i.wait -= dt
}

func (i *introState) draw(dst *ebiten.Image, cam *camera.Camera) {
	i.bg.Draw(dst, cam)
	i.p.draw(dst, cam)
}
