package game

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"golang.org/x/image/font/basicfont"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/game1/game/keymap"
	"github.com/Bredgren/game1/game/keymap/button"
	"github.com/Bredgren/game1/game/ui"
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
	remapText    *ui.Text

	menu           ui.Drawer
	btns           map[keymap.Action]*ui.Button
	keyText        map[keymap.Action]*ui.Text
	gamepadText    map[keymap.Action]*ui.Text
	canClickButton bool
}

func newMainMenu(p *player, screenHeight int, cam *camera.Camera, bg *background,
	km keymap.Layers) *mainMenuState {
	m := &mainMenuState{
		p:            p,
		screenHeight: screenHeight,
		cam:          cam,
		bg:           bg,
		keymap:       km,

		// remap:       true,
		// remapAction: left,

		btns:           map[keymap.Action]*ui.Button{},
		keyText:        map[keymap.Action]*ui.Text{},
		gamepadText:    map[keymap.Action]*ui.Text{},
		canClickButton: true,
	}

	m.setupMenu()
	m.setupKeymap()

	return m
}

func (m *mainMenuState) setupMenu() {

	idleImg, _ := ebiten.NewImage(100, 20, ebiten.FilterNearest)
	idleImg.Fill(color.NRGBA{200, 200, 200, 50})
	hoverImg, _ := ebiten.NewImage(100, 20, ebiten.FilterNearest)
	hoverImg.Fill(color.NRGBA{100, 100, 100, 50})

	var elements []ui.WeightedDrawer

	m.remapText = &ui.Text{
		Anchor: ui.AnchorCenter,
		Color:  color.Black,
		Face:   basicfont.Face7x13,
		Wt:     0.5,
	}

	elements = append(elements, m.remapText)

	actions := []keymap.Action{
		left, right,
	}
	m.keyText = map[keymap.Action]*ui.Text{}
	m.gamepadText = map[keymap.Action]*ui.Text{}
	for _, action := range actions {
		action := action
		m.keyText[action] = &ui.Text{
			Anchor: ui.AnchorLeft,
			Color:  color.Black,
			Face:   basicfont.Face7x13,
			Wt:     1,
		}
		m.gamepadText[action] = &ui.Text{
			Anchor: ui.AnchorLeft,
			Color:  color.Black,
			Face:   basicfont.Face7x13,
			Wt:     1,
		}
		m.btns[action] = &ui.Button{
			IdleImg:     idleImg,
			HoverImg:    hoverImg,
			IdleAnchor:  ui.AnchorCenter,
			HoverAnchor: ui.AnchorCenter,
			Element: &ui.HorizontalContainer{
				Wt: 1,
				Elements: []ui.WeightedDrawer{
					&ui.Text{
						Text: string(action),
						Anchor: ui.Anchor{
							Src:    geo.VecXY(0, 0.5),
							Dst:    geo.VecXY(0, 0.5),
							Offset: geo.VecXY(5, 0),
						},
						Color: color.Black,
						Face:  basicfont.Face7x13,
						Wt:    2,
					},
					m.gamepadText[action],
					m.keyText[action],
				},
			},
			Wt: 1,
			OnClick: func() {
				log.Println("remap", action)
				m.remap = true
				m.remapAction = action
				m.remapText.Text = fmt.Sprintf("Remap action '%s'", action)
			},
		}
		elements = append(elements, m.btns[action])
	}

	m.menu = &ui.VerticalContainer{
		Wt:       1,
		Elements: elements,
	}

	m.updateText()
}

func (m *mainMenuState) updateText() {
	actions := []keymap.Action{
		left, right,
	}
	for _, action := range actions {
		if btn, ok := m.keymap[playerLayer].KeyMouse.GetButton(action); ok {
			m.keyText[action].Text = fmt.Sprintf("%d", btn)
		} else {
			m.keyText[action].Text = "-"
		}
		if btn, ok := m.keymap[playerLayer].GamepadBtn.GetButton(action); ok {
			m.gamepadText[action].Text = fmt.Sprintf("%d", btn)
		} else {
			m.gamepadText[action].Text = "-"
		}
	}
}

func (m *mainMenuState) setupKeymap() {
	//// Setup remap layer
	// Button handlers
	remapHandlers := keymap.ButtonHandlerMap{}
	for key := ebiten.Key0; key < ebiten.KeyMax; key++ {
		action := keymap.Action(fmt.Sprintf("key%d", key))
		remapHandlers[action] = m.keyRemapHandler(button.FromKey(key))
	}
	remapHandlers[keymap.Action("mouse0")] = m.keyRemapHandler(button.FromMouse(ebiten.MouseButtonLeft))
	remapHandlers[keymap.Action("mouse1")] = m.keyRemapHandler(button.FromMouse(ebiten.MouseButtonMiddle))
	remapHandlers[keymap.Action("mouse2")] = m.keyRemapHandler(button.FromMouse(ebiten.MouseButtonRight))

	// Gamepad handlers
	for btn := ebiten.GamepadButton0; btn < ebiten.GamepadButtonMax; btn++ {
		action := keymap.Action(fmt.Sprintf("btn%d", btn))
		remapHandlers[action] = m.btnRemapHandler(btn)
	}

	// Axis handlers
	axisHandlers := keymap.AxisHandlerMap{}
	// // We don't know how many axes there will be at this point so just do alot :P
	// for axis := 0; axis < 100; axis++ {
	// 	action := keymap.Action(fmt.Sprintf("axis%d", axis))
	// 	axisHandlers[action] = m.axisRemapHandler(axis)
	// }

	m.keymap[remapLayer] = keymap.New(remapHandlers, axisHandlers)

	// Button actions
	for key := ebiten.Key0; key < ebiten.KeyMax; key++ {
		action := keymap.Action(fmt.Sprintf("key%d", key))
		m.keymap[remapLayer].KeyMouse.Set(button.FromKey(key), action)
	}
	m.keymap[remapLayer].KeyMouse.Set(button.FromMouse(ebiten.MouseButtonLeft), "mouse0")
	m.keymap[remapLayer].KeyMouse.Set(button.FromMouse(ebiten.MouseButtonMiddle), "mouse1")
	m.keymap[remapLayer].KeyMouse.Set(button.FromMouse(ebiten.MouseButtonRight), "mouse2")

	// Gamepad actions
	for btn := ebiten.GamepadButton0; btn < ebiten.GamepadButtonMax; btn++ {
		action := keymap.Action(fmt.Sprintf("btn%d", btn))
		m.keymap[remapLayer].GamepadBtn.Set(btn, action)
	}

	// Axis actions
	// for axis := 0; axis < 100; axis++ {
	// 	action := keymap.Action(fmt.Sprintf("axis%d", axis))
	// 	m.keymap[remapLayer].GamepadAxis.Set(axis, action)
	// }

	//// Setup UI handlers
	leftClickHandlers := keymap.ButtonHandlerMap{
		leftClick: m.leftMouseHandler,
	}
	m.keymap[leftClickLayer] = keymap.New(leftClickHandlers, nil)
	m.keymap[leftClickLayer].KeyMouse.Set(button.FromMouse(ebiten.MouseButtonLeft), leftClick)

	// UI handlers
	uiHandlers := keymap.ButtonHandlerMap{
	// left: m.keyLabels[left].handleBtn,
	// right: m.keyLabels[right].handleBtn,
	// jump:  m.keyLabels[jump].handleBtn,
	}
	uiAxisHandlers := keymap.AxisHandlerMap{
	// move: m.keyLabels[move].handleAxis,
	}
	m.keymap[uiLayer] = keymap.New(uiHandlers, uiAxisHandlers)
	setDefaultKeyMap(m.keymap[uiLayer])
}

func (m *mainMenuState) keyRemapHandler(btn button.KeyMouse) keymap.ButtonHandler {
	return func(down bool) bool {
		if !m.canClickButton && btn.IsMouse() {
			// This prevents us from always immediately remapping to left mouse
			return false
		}

		if down && m.remap {
			log.Println("remap key to", btn)
			m.keymap[playerLayer].KeyMouse.Set(btn, m.remapAction)
			m.remap = false
			m.remapText.Text = ""
			m.updateText()

			if btn.IsMouse() {
				// This prevents us from clicking a button if remapping to left mouse while hover
				// over a button
				m.canClickButton = false
			}

			return true
		}

		// No reason to stop propagation here because either the button is up or is not
		// remappable
		return false
	}
}

func (m *mainMenuState) btnRemapHandler(btn ebiten.GamepadButton) keymap.ButtonHandler {
	return func(down bool) bool {
		remap := m.remap
		if down && remap {
			log.Println("remap gamepad btn to", btn)
			m.keymap[playerLayer].GamepadBtn.Set(btn, m.remapAction)
			m.remap = false
			m.updateText()
		}
		return remap
	}
}

// func (m *mainMenuState) axisRemapHandler(axis int) keymap.AxisHandler {
// 	return func(val float64) bool {
// 		remap := m.remap
// 		if val != 0 && remap {
// 			log.Println("remap axis to", axis)
// 			m.keymap[playerLayer].GamepadAxis.Set(axis, m.remapAction)
// 			m.remap = false
// 		}
// 		return remap
// 	}
// }

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

	for _, b := range m.btns {
		b.Update()
	}
}

func (m *mainMenuState) draw(dst *ebiten.Image, cam *camera.Camera) {
	m.bg.Draw(dst, cam)
	m.p.draw(dst, cam)

	// ebitenutil.DrawRect(dst, 100, 150, 100, 100, color.NRGBA{100, 100, 100, 50})
	m.menu.Draw(dst, geo.RectXYWH(100, 150, 100, 100))
}

func (m *mainMenuState) leftMouseHandler(down bool) bool {
	if m.canClickButton && down {
		for _, b := range m.btns {
			if b.Hover {
				b.OnClick()
				m.canClickButton = false
				return true
			}
		}
	}
	m.canClickButton = !down
	return false
}
