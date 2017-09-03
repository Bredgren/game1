package game

import (
	"fmt"
	"image/color"
	"time"

	"golang.org/x/image/font/basicfont"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/game1/game/keymap"
	"github.com/Bredgren/game1/game/keymap/button"
	"github.com/Bredgren/game1/game/ui"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

const (
	buttonWidth      = 350
	buttonHeight     = 16
	axisButtonWidth  = 100
	axisButtonHeight = 14
)

type mainMenuState struct {
	p            *player
	screenHeight int
	cam          *camera.Camera
	bg           *background
	keymap       keymap.Layers
	remapAction  keymap.Action
	remap        bool
	remapAxis    bool
	remapText    *ui.Text

	menu           ui.Drawer
	btns           map[keymap.Action]*ui.Button
	actionText     map[keymap.Action]*ui.Text
	keyText        map[keymap.Action]*ui.Text
	gamepadText    map[keymap.Action]*ui.Text
	canClickButton bool

	axisMenu    ui.Drawer
	axisBtns    map[int]*ui.Button
	axisValText map[int]*ui.Text
}

func newMainMenu(p *player, screenHeight int, cam *camera.Camera, bg *background,
	km keymap.Layers) *mainMenuState {
	m := &mainMenuState{
		p:            p,
		screenHeight: screenHeight,
		cam:          cam,
		bg:           bg,
		keymap:       km,

		btns:           map[keymap.Action]*ui.Button{},
		actionText:     map[keymap.Action]*ui.Text{},
		keyText:        map[keymap.Action]*ui.Text{},
		gamepadText:    map[keymap.Action]*ui.Text{},
		canClickButton: true,

		axisBtns:    map[int]*ui.Button{},
		axisValText: map[int]*ui.Text{},
	}

	m.setupMenu()
	m.setupKeymap()

	return m
}

func (m *mainMenuState) setupMenu() {
	idleImg, _ := ebiten.NewImage(buttonWidth, buttonHeight, ebiten.FilterNearest)
	idleImg.Fill(color.NRGBA{200, 200, 200, 50})
	hoverImg, _ := ebiten.NewImage(buttonWidth, buttonHeight, ebiten.FilterNearest)
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
		left, right, move, jump, punch, punchH, punchV, uppercut, slam, launch,
	}
	for _, action := range actions {
		action := action // For use in callbacks

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
		m.actionText[action] = &ui.Text{
			Text: string(action),
			Anchor: ui.Anchor{
				Src:    geo.VecXY(0, 0.5),
				Dst:    geo.VecXY(0, 0.5),
				Offset: geo.VecXY(5, 0),
			},
			Color: color.Black,
			Face:  basicfont.Face7x13,
			Wt:    1.6,
		}

		var onClick func()

		_, isAxis := m.keymap[playerLayer].GamepadAxis.GetAxis(action)

		if isAxis {
			onClick = func() {
				m.remapAxis = true
				m.remapAction = action
				m.remapText.Text = fmt.Sprintf("Select new axis for '%s'", action)
			}
		} else {
			onClick = func() {
				m.remapAxis = false // to close the axis window if it's open
				m.remap = true
				m.remapAction = action
				m.remapText.Text = fmt.Sprintf("Press new key/mouse/gamepad button for '%s'", action)
			}
		}
		m.btns[action] = &ui.Button{
			IdleImg:     idleImg,
			HoverImg:    hoverImg,
			IdleAnchor:  ui.AnchorCenter,
			HoverAnchor: ui.AnchorCenter,
			Element: &ui.HorizontalContainer{
				Wt: 1,
				Elements: []ui.WeightedDrawer{
					m.actionText[action],
					m.keyText[action],
					m.gamepadText[action],
				},
			},
			Wt:      1,
			OnClick: onClick,
		}
		elements = append(elements, m.btns[action])
	}

	actions = []keymap.Action{
		pause, fullscreen,
	}
	for _, action := range actions {
		action := action // For use in callbacks
		b1, _ := m.keymap[generalLayer].KeyMouse.GetButton(action)
		b2, _ := m.keymap[generalLayer].GamepadBtn.GetButton(action)
		elements = append(elements, &ui.HorizontalContainer{
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
					Wt:    1.6,
				},
				&ui.Text{
					Text:   b1.String(),
					Anchor: ui.AnchorLeft,
					Color:  color.Black,
					Face:   basicfont.Face7x13,
					Wt:     1,
				},
				&ui.Text{
					Text:   fmt.Sprintf("Gamepad %d", b2),
					Anchor: ui.AnchorLeft,
					Color:  color.Black,
					Face:   basicfont.Face7x13,
					Wt:     1,
				},
			},
		})
	}

	idleImg, _ = ebiten.NewImage(buttonWidth/3, buttonHeight, ebiten.FilterNearest)
	idleImg.Fill(color.NRGBA{200, 200, 200, 50})
	hoverImg, _ = ebiten.NewImage(buttonWidth/3, buttonHeight, ebiten.FilterNearest)
	hoverImg.Fill(color.NRGBA{100, 100, 100, 50})
	b := &ui.Button{
		IdleImg:     idleImg,
		HoverImg:    hoverImg,
		IdleAnchor:  ui.AnchorCenter,
		HoverAnchor: ui.AnchorCenter,
		Element: &ui.HorizontalContainer{
			Wt: 1,
			Elements: []ui.WeightedDrawer{
				&ui.Text{
					Text:   "Restore Default",
					Anchor: ui.AnchorCenter,
					Color:  color.Black,
					Face:   basicfont.Face7x13,
					Wt:     1,
				},
			},
		},
		Wt: 1,
		OnClick: func() {
			setDefaultKeyMap(m.keymap[playerLayer])
			setDefaultKeyMap(m.keymap[uiLayer])
			m.updateText()
		},
	}

	m.btns[keymap.Action("reset")] = b
	elements = append(elements, b)

	m.menu = &ui.VerticalContainer{
		Wt:       1,
		Elements: elements,
	}

	m.updateText()
}

func (m *mainMenuState) setupAxisMenu() {
	idleImg, _ := ebiten.NewImage(axisButtonWidth, axisButtonHeight, ebiten.FilterNearest)
	idleImg.Fill(color.NRGBA{200, 200, 200, 50})
	hoverImg, _ := ebiten.NewImage(axisButtonWidth, axisButtonHeight, ebiten.FilterNearest)
	hoverImg.Fill(color.NRGBA{100, 100, 100, 50})

	var elements []ui.WeightedDrawer
	elements = append(elements, &ui.Text{
		Anchor: ui.AnchorCenter,
		Color:  color.Black,
		Face:   basicfont.Face7x13,
		Text:   "Select Axis",
		Wt:     1,
	})

	for axis := 0; axis < ebiten.GamepadAxisNum(0); axis++ {
		axis := axis
		m.axisValText[axis] = &ui.Text{
			Anchor: ui.AnchorLeft,
			Color:  color.Black,
			Face:   basicfont.Face7x13,
			Text:   "(0)",
			Wt:     1,
		}
		m.axisBtns[axis] = &ui.Button{
			IdleImg:     idleImg,
			HoverImg:    hoverImg,
			IdleAnchor:  ui.AnchorCenter,
			HoverAnchor: ui.AnchorCenter,
			Element: &ui.HorizontalContainer{
				Wt: 1,
				Elements: []ui.WeightedDrawer{
					&ui.Text{
						Anchor: ui.Anchor{
							Src:    geo.VecXY(0, 0.5),
							Dst:    geo.VecXY(0, 0.5),
							Offset: geo.VecXY(2, 0),
						},
						Color: color.Black,
						Face:  basicfont.Face7x13,
						Text:  fmt.Sprintf("Axis %d", axis),
						Wt:    1,
					},
					m.axisValText[axis],
				},
			},
			Wt: 1,
			OnClick: func() {
				m.keymap[playerLayer].GamepadAxis.Set(axis, m.remapAction)
				m.keymap[uiLayer].GamepadAxis.Set(axis, m.remapAction)
				m.remapAxis = false
				m.updateText()
			},
		}
		elements = append(elements, m.axisBtns[axis])
	}

	m.axisMenu = &ui.VerticalContainer{
		Wt:       1,
		Elements: elements,
	}
}

func (m *mainMenuState) updateText() {
	actions := []keymap.Action{
		left, right, move, jump, punch, punchH, punchV, uppercut, slam, launch,
	}
	for _, action := range actions {
		if btn, ok := m.keymap[playerLayer].KeyMouse.GetButton(action); ok {
			m.keyText[action].Text = btn.String()
			m.keyText[action].Color = color.Black
		} else {
			m.keyText[action].Text = "N/A"
			if _, valid := defaultKeyMap.KeyMouse.GetButton(action); valid {
				m.keyText[action].Color = color.NRGBA{200, 0, 0, 200}
			} else {
				m.keyText[action].Color = color.NRGBA{0, 0, 0, 100}
			}
		}

		if btn, ok := m.keymap[playerLayer].GamepadBtn.GetButton(action); ok {
			m.gamepadText[action].Text = fmt.Sprintf("Gamepad %d", btn)
			m.gamepadText[action].Color = color.Black
		} else if axis, ok := m.keymap[playerLayer].GamepadAxis.GetAxis(action); ok {
			m.gamepadText[action].Text = fmt.Sprintf("Axis %d", axis)
			m.gamepadText[action].Color = color.Black
		} else {
			m.gamepadText[action].Text = "N/A"
			_, validBtn := defaultKeyMap.GamepadBtn.GetButton(action)
			_, validAxis := defaultKeyMap.GamepadAxis.GetAxis(action)
			if validBtn || validAxis {
				m.gamepadText[action].Color = color.NRGBA{200, 0, 0, 200}
			} else {
				m.gamepadText[action].Color = color.NRGBA{0, 0, 0, 100}
			}
		}
	}
}

func (m *mainMenuState) setupKeymap() {
	//// Setup remap layer
	// Button handlers
	remapHandlers := keymap.ButtonHandlerMap{}
	for key := ebiten.Key0; key <= ebiten.KeyMax; key++ {
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

	m.keymap[remapLayer] = keymap.New(remapHandlers, nil)

	// Button actions
	for key := ebiten.Key0; key <= ebiten.KeyMax; key++ {
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

	//// Setup UI handlers
	leftClickHandlers := keymap.ButtonHandlerMap{
		leftClick: m.leftMouseHandler,
	}
	m.keymap[leftClickLayer] = keymap.New(leftClickHandlers, nil)
	m.keymap[leftClickLayer].KeyMouse.Set(button.FromMouse(ebiten.MouseButtonLeft), leftClick)

	colorFn := func(action keymap.Action) keymap.ButtonHandler {
		return func(down bool) bool {
			if down {
				m.actionText[action].Color = color.White
			} else {
				m.actionText[action].Color = color.Black
			}
			return false
		}
	}

	axisFn := func(action keymap.Action) keymap.AxisHandler {
		return func(val float64) bool {
			var axis int
			fmt.Sscanf(m.gamepadText[action].Text, "Axis %d", &axis)
			m.gamepadText[action].Text = fmt.Sprintf("Axis %d (%.2f)", axis, val)
			return false
		}
	}

	// UI handlers
	uiHandlers := keymap.ButtonHandlerMap{
		left:     colorFn(left),
		right:    colorFn(right),
		jump:     colorFn(jump),
		uppercut: colorFn(uppercut),
		slam:     colorFn(slam),
		punch:    colorFn(punch),
		launch:   colorFn(launch),
	}
	uiAxisHandlers := keymap.AxisHandlerMap{
		move:   axisFn(move),
		punchH: axisFn(punchH),
		punchV: axisFn(punchV),
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

		_, valid := defaultKeyMap.KeyMouse.GetButton(m.remapAction)
		if down && m.remap && valid {
			m.keymap[playerLayer].KeyMouse.Set(btn, m.remapAction)
			m.keymap[uiLayer].KeyMouse.Set(btn, m.remapAction)
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
		_, valid := defaultKeyMap.GamepadBtn.GetButton(m.remapAction)
		if down && m.remap && valid {
			m.keymap[playerLayer].GamepadBtn.Set(btn, m.remapAction)
			m.keymap[uiLayer].GamepadBtn.Set(btn, m.remapAction)
			m.remap = false
			m.updateText()
		}

		// No reason to stop propagation here because either the button is up or is not
		// remappable
		return false
	}
}

func (m *mainMenuState) begin(previousState gameStateName) {
	m.cam.Target = fixedCameraTarget{geo.VecXY(m.p.pos.X, -float64(m.screenHeight)*0.4)}
	if m.axisMenu == nil {
		// Initialize here so that we have the correct number of gamepad axes.
		m.setupAxisMenu()
	}
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
	if m.remapAxis {
		for _, b := range m.axisBtns {
			b.Update()
		}

		for axis := 0; axis < ebiten.GamepadAxisNum(0); axis++ {
			m.axisValText[axis].Text = fmt.Sprintf("(%.2f)", ebiten.GamepadAxis(0, axis))
		}
	}
}

func (m *mainMenuState) draw(dst *ebiten.Image, cam *camera.Camera) {
	m.bg.Draw(dst, cam)
	m.p.draw(dst, cam)

	x, y := 120.0, 20.0
	height := 229.0
	m.menu.Draw(dst, geo.RectXYWH(x, y, buttonWidth, height))

	if m.remapAxis {
		height = 106
		x, y = x+buttonWidth+10, y+50
		m.axisMenu.Draw(dst, geo.RectXYWH(x, y, axisButtonWidth, height))
	}
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
		if m.remapAxis {
			for _, b := range m.axisBtns {
				if b.Hover {
					b.OnClick()
					m.canClickButton = false
					return true
				}
			}
		}
	}
	m.canClickButton = !down
	return false
}
