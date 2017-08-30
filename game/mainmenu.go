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
}

func newMainMenu(p *player, screenHeight int, cam *camera.Camera) *mainMenuState {
	return &mainMenuState{
		p:            p,
		screenHeight: screenHeight,
		cam:          cam,
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

}

func (m *mainMenuState) draw(dst *ebiten.Image, cam *camera.Camera) {

}
