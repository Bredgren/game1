package game

import (
	"time"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

type mainMenuState struct {
	p            *player
	screenHeight int
	cam          *camera.Camera
	bg           *background
}

func newMainMenu(p *player, screenHeight int, cam *camera.Camera, bg *background) *mainMenuState {
	return &mainMenuState{
		p:            p,
		screenHeight: screenHeight,
		cam:          cam,
		bg:           bg,
	}
}

func (m *mainMenuState) begin(previousState gameStateName) {
	m.cam.Target = fixedCameraTarget{geo.VecXY(m.p.pos.X, -float64(m.screenHeight)*0.4)}
}

func (m *mainMenuState) end() {

}

func (m *mainMenuState) nextState() gameStateName {
	return mainMenu
}

func (m *mainMenuState) update(dt time.Duration) {
	m.p.update(dt)
}

func (m *mainMenuState) draw(dst *ebiten.Image, cam *camera.Camera) {
	m.bg.Draw(dst, cam)
	m.p.draw(dst, cam)
}
