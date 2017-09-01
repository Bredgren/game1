package game

import (
	"fmt"
	"log"
	"time"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/game1/game/keymap"
	"github.com/Bredgren/game1/game/keymap/button"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

type mainMenuState struct {
	p            *player
	screenHeight int
	cam          *camera.Camera
	bg           *background
	keymap       keymap.Layers
	remapAction  keymap.Action
	remap        bool
}

func newMainMenu(p *player, screenHeight int, cam *camera.Camera, bg *background,
	km keymap.Layers) *mainMenuState {
	m := &mainMenuState{
		p:            p,
		screenHeight: screenHeight,
		cam:          cam,
		bg:           bg,
		keymap:       km,
		// remap:        true,
		// remapAction:  jump,
	}
	m.setupKeymap()
	return m
}

func (m *mainMenuState) setupKeymap() {
	remapHandlers := keymap.ButtonHandlerMap{}
	for key := ebiten.Key0; key < ebiten.KeyMax; key++ {
		remapHandlers[keymap.Action(fmt.Sprintf("key%d", key))] = m.keyRemapHandler(key)
	}
	m.keymap[remapLayer] = keymap.New(remapHandlers, nil)

	for key := ebiten.Key0; key < ebiten.KeyMax; key++ {
		m.keymap[remapLayer].KeyMouse.Set(button.FromKey(key), keymap.Action(fmt.Sprintf("key%d", key)))
	}
	log.Println(m.keymap[remapLayer])

	m.keymap[uiLayer] = keymap.New(nil, nil)
}

func (m *mainMenuState) keyRemapHandler(key ebiten.Key) keymap.ButtonHandler {
	return func(down bool) bool {
		remap := m.remap
		if down && remap {
			log.Println("remap to", key)
			m.keymap[playerLayer].KeyMouse.Set(button.FromKey(key), m.remapAction)
			m.remap = false
		}
		return remap
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
