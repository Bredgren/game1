package keymap

import (
	"github.com/Bredgren/game1/game/keymap/button"
	"github.com/hajimehoshi/ebiten"
)

// Action is a name/label for an action.
type Action string

// KeyMap groups all map types.
type KeyMap struct {
	KeyMouse    *KeyMouseMap
	GamepadBtn  *GamepadBtnMap
	GamepadAxis *GamepadAxisMap
	btnHandlers ButtonHandlerMap
	gaHandlers  AxisHandlerMap
}

// New creates and returns a new, empty KeyMap. The btnHandlers map shared between keyboard/mouse
// and gamepad buttons. This means that if an action handler for a keyboard/mouse buttton
// stops propagation then a gamepad button that maps to the same action in a later layer
// will not be handled.
func New(btnHandlers ButtonHandlerMap, gamepadAxisHandlers AxisHandlerMap) *KeyMap {
	return &KeyMap{
		KeyMouse:    NewKeyMouseMap(),
		GamepadBtn:  NewGamepadBtnMap(),
		GamepadAxis: NewGamepadAxisMap(),
		btnHandlers: btnHandlers,
		gaHandlers:  gamepadAxisHandlers,
	}
}

// ButtonHandler is a function that handles a button state. The parameter down is true
// if the button is pressed. It should return true if no later handlers should be called
// for the same button.
type ButtonHandler func(down bool) (stopPropagation bool)

// AxisHandler is a function that handles gamepad axis state. The function will directly
// be given the result of ebiten.GamepadAxis. It should return true if no later handlers
// should be called for the same axis.
type AxisHandler func(val float64) (stopPropagation bool)

// ButtonHandlerMap maps button Actions to their handlers.
type ButtonHandlerMap map[Action]ButtonHandler

// AxisHandlerMap maps axis Actions to their handlers.
type AxisHandlerMap map[Action]AxisHandler

// Layers is a slice of KeyMaps. It enables buttons to be overloaded with multiple actions
// with the option of skipping later handlers for buttons handled at earlier layers.
type Layers []*KeyMap

// Update checks input state and calls handlers for any actions triggered. It handles
// each layer in order. If the handler for a button stops propagation then later following
// layers will not handle any actions the same button triggers.
func (l Layers) Update() {
	const gamepadID = 0
	stoppedKeys := map[button.KeyMouse]bool{}
	stoppedBtns := map[ebiten.GamepadButton]bool{}
	stoppedAxes := map[int]bool{}
	for _, keymap := range l {
		actions := map[Action]bool{}

		for _, btn := range keymap.KeyMouse.Buttons() {
			if stoppedKeys[btn] {
				continue
			}
			var down bool
			if k, ok := btn.Key(); ok {
				down = ebiten.IsKeyPressed(k)
			} else if mb, ok := btn.Mouse(); ok {
				down = ebiten.IsMouseButtonPressed(mb)
			}
			if a, ok := keymap.KeyMouse.GetAction(btn); ok {
				actions[a] = down || actions[a]
			}
		}

		for _, btn := range keymap.GamepadBtn.Buttons() {
			if stoppedBtns[btn] {
				continue
			}
			down := ebiten.IsGamepadButtonPressed(gamepadID, btn)
			if a, ok := keymap.GamepadBtn.GetAction(btn); ok {
				actions[a] = down || actions[a]
			}
		}

		for action, down := range actions {
			stop := keymap.btnHandlers[action](down)
			if b, ok := keymap.KeyMouse.GetButton(action); ok {
				stoppedKeys[b] = stop
			}
			if b, ok := keymap.GamepadBtn.GetButton(action); ok {
				stoppedBtns[b] = stop
			}
		}

		numAxis := ebiten.GamepadAxisNum(gamepadID)
		for _, axis := range keymap.GamepadAxis.Axes() {
			if stoppedAxes[axis] || axis >= numAxis {
				continue
			}
			val := ebiten.GamepadAxis(gamepadID, axis)
			if act, ok := keymap.GamepadAxis.GetAction(axis); ok {
				stop := keymap.gaHandlers[act](val)
				stoppedAxes[axis] = stop
			}
		}
	}
}
