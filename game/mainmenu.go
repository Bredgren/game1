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
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type mainMenuState struct {
	p              *player
	screenHeight   int
	cam            *camera.Camera
	bg             *background
	keymap         keymap.Layers
	remapAction    keymap.Action
	remap          bool
	keyLabels      map[keymap.Action]*keyLabel
	menu           ui.Drawer
	btn            *ui.Button
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
		// remapAction: jump,

		keyLabels:      map[keymap.Action]*keyLabel{},
		canClickButton: true,
	}

	idleImg, _ := ebiten.NewImage(40, 20, ebiten.FilterNearest)
	idleImg.Fill(color.NRGBA{200, 200, 200, 50})
	hoverImg, _ := ebiten.NewImage(40, 20, ebiten.FilterNearest)
	hoverImg.Fill(color.NRGBA{200, 200, 200, 200})

	m.btn = &ui.Button{
		IdleImg:  idleImg,
		HoverImg: hoverImg,
		IdleAnchor: ui.Anchor{
			Src: geo.VecXY(0.5, 0.5),
			Dst: geo.VecXY(0.5, 0.5),
		},
		HoverAnchor: ui.Anchor{
			Src: geo.VecXY(0.5, 0.5),
			Dst: geo.VecXY(0.5, 0.5),
		},
		Element: &ui.Text{
			Anchor: ui.Anchor{
				Src: geo.VecXY(0.5, 0.5),
				Dst: geo.VecXY(0.5, 0.5),
			},
			Text:  "btn",
			Color: color.Black,
			Face:  basicfont.Face7x13,
			Wt:    1,
		},
		Wt: 1,
	}

	m.menu = &ui.VerticalContainer{
		Wt: 1,
		Elements: []ui.WeightedDrawer{
			&ui.HorizontalContainer{
				Wt: 1,
				Elements: []ui.WeightedDrawer{
					&ui.Text{
						Anchor: ui.Anchor{
							Src:    geo.VecXY(0.5, 0.5),
							Dst:    geo.VecXY(0.5, 0.5),
							Offset: geo.VecXY(0, 0),
						},
						Text:  "text1",
						Color: color.Black,
						Face:  basicfont.Face7x13,
						Wt:    1,
					},
					m.btn,
				},
			},
			&ui.HorizontalContainer{
				Wt: 1,
				Elements: []ui.WeightedDrawer{
					&ui.Text{
						Anchor: ui.Anchor{
							Src:    geo.VecXY(1, 0.5),
							Dst:    geo.VecXY(1, 0.5),
							Offset: geo.VecXY(0, -10),
						},
						Text:  "text2",
						Color: color.Black,
						Face:  basicfont.Face7x13,
						Wt:    1,
					},
					&ui.Text{
						Anchor: ui.Anchor{
							Src:    geo.VecXY(0, 0.5),
							Dst:    geo.VecXY(0, 0.5),
							Offset: geo.VecXY(0, 10),
						},
						Text:  "text3",
						Color: color.Black,
						Face:  basicfont.Face7x13,
						Wt:    1,
					},
				},
			},
		},
	}

	m.setupKeyLabels()
	m.setupKeymap()

	return m
}

func (m *mainMenuState) setupKeyLabels() {
	keyOptionsPos := geo.VecXY(100, 100)
	keyOptionVGap := 2.0
	keyLabels := []*keyLabel{
		newKeyLabel(left, geo.RectCornersVec(keyOptionsPos, keyOptionsPos.Plus(geo.VecXY(50, 20))), basicfont.Face7x13),
		// newKeyLabel(right, basicfont.Face7x13),
		// newKeyLabel(move, basicfont.Face7x13),
		// newKeyLabel(jump, basicfont.Face7x13),
	}

	for _, kl := range keyLabels {
		kl.bounds.SetTopLeft(keyOptionsPos.XY())
		keyOptionsPos.Y += kl.bounds.H + keyOptionVGap
		m.keyLabels[kl.action] = kl
	}
}

func (m *mainMenuState) setupKeymap() {
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

	// UI handlers
	uiHandlers := keymap.ButtonHandlerMap{
		left: m.keyLabels[left].handleBtn,
		// right: m.keyLabels[right].handleBtn,
		// jump:  m.keyLabels[jump].handleBtn,
		click: m.handleMouseDown,
	}
	uiAxisHandlers := keymap.AxisHandlerMap{
	// move: m.keyLabels[move].handleAxis,
	}
	m.keymap[uiLayer] = keymap.New(uiHandlers, uiAxisHandlers)
	setDefaultKeyMap(m.keymap[uiLayer])
	m.keymap[uiLayer].KeyMouse.Set(button.FromMouse(ebiten.MouseButtonLeft), click)
}

func (m *mainMenuState) keyRemapHandler(btn button.KeyMouse) keymap.ButtonHandler {
	return func(down bool) bool {
		remap := m.remap
		if down && remap {
			log.Println("remap key to", btn)
			m.keymap[playerLayer].KeyMouse.Set(btn, m.remapAction)
			m.remap = false
		}
		return remap
	}
}

func (m *mainMenuState) btnRemapHandler(btn ebiten.GamepadButton) keymap.ButtonHandler {
	return func(down bool) bool {
		remap := m.remap
		if down && remap {
			log.Println("remap gamepad btn to", btn)
			m.keymap[playerLayer].GamepadBtn.Set(btn, m.remapAction)
			m.remap = false
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
	m.btn.Update()
}

func (m *mainMenuState) draw(dst *ebiten.Image, cam *camera.Camera) {
	m.bg.Draw(dst, cam)
	m.p.draw(dst, cam)

	for _, kl := range m.keyLabels {
		kl.draw(dst, cam)
	}

	ebitenutil.DrawRect(dst, 100, 150, 100, 100, color.NRGBA{100, 100, 100, 50})
	m.menu.Draw(dst, geo.RectXYWH(100, 150, 100, 100))
}

func (m *mainMenuState) handleMouseDown(down bool) bool {
	if m.canClickButton && down && m.btn.Hover {
		log.Println("click")
		m.canClickButton = false
		return true
	}
	m.canClickButton = !down
	return false
}
