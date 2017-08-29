package keymap

import (
	"github.com/Bredgren/game1/game/keymap/button"
	"github.com/hajimehoshi/ebiten"
)

// Action is a name/label for an action.
type Action string

// ButtonHandler is a function that handles a button state. The parameter down is true
// if the button is pressed. It should return true if no later handlers should be called
// for the same button.
type ButtonHandler func(down bool) (stopPropagation bool)

// AxisHandler is a function that handles gamepad axis state. The function will directly
// be given the result of ebiten.GamepadAxis. It should return true if no later handlers
// should be called for the same axis.
type AxisHandler func(val float64) (stopPropagation bool)

// ActionHandlerMap maps an Action to its handler.
type ActionHandlerMap map[Action]ButtonHandler

// AxisActionHandlerMap maps an Action to its axis handler.
type AxisActionHandlerMap map[Action]AxisHandler

// KeyMap maps a Button to an Action.
type KeyMap map[button.Button]Action

// StoppedBtnSet holds a set of Buttons that stopped propagation.
type StoppedBtnSet map[button.Button]bool

// Update calls the handler for all buttons that haven't been stopped. A button is stopped
// if it maps to true in stoppedBtns. The stoppedBtns map is updated according to the
// the return values of the handlers that are executed.
func (km KeyMap) Update(ahm ActionHandlerMap, stoppedBtns StoppedBtnSet) {
	gamepadID := 0 // Assume one gamepad for now

	// We need to OR together all "down" vaules for different buttons that map to the same
	// action, otherwise they can cancel each other out.
	actionValues := map[Action]bool{}

	for btn, action := range km {
		if stoppedBtns[btn] {
			continue
		}

		if _, ok := ahm[action]; ok {
			var down bool
			if k, ok := btn.Key(); ok {
				down = ebiten.IsKeyPressed(k)
			} else if gb, ok := btn.GamepadButton(); ok {
				down = ebiten.IsGamepadButtonPressed(gamepadID, gb)
			} else if mb, ok := btn.MouseButton(); ok {
				down = ebiten.IsMouseButtonPressed(mb)
			}
			actionValues[action] = down || actionValues[action]
		}
	}

	for action, down := range actionValues {
		res := ahm[action](down)
		for btn, a := range km {
			if a == action {
				stoppedBtns[btn] = res
			}
		}
	}
}

// AxisMap maps an axis id to an Action.
type AxisMap map[int]Action

// StoppedAxisSet holds a set of axis that stopped propagation.
type StoppedAxisSet map[int]bool

// Update calls the handler for all axes that haven't been stopped. An axis is stopped
// if it maps to true in stoppedAxis. The stoppedAxis map is updated according to the
// the return values of the handlers that are executed.
func (am AxisMap) Update(ahm AxisActionHandlerMap, stoppedAxis StoppedAxisSet) {
	gamepadID := 0 // Assume one gamepad for now
	numAxis := ebiten.GamepadAxisNum(0)

	for axis, action := range am {
		if stoppedAxis[axis] || axis > numAxis {
			continue
		}

		if actionFn, ok := ahm[action]; ok {
			stoppedAxis[axis] = actionFn(ebiten.GamepadAxis(gamepadID, axis))
		}
	}
}

// Map combines KeyMap and AxisMap.
type Map struct {
	KeyMap
	AxisMap
}

// NewMap creates a new Map type with empty KeyMap and AxisMap.
func NewMap() Map {
	return Map{
		KeyMap:  make(KeyMap),
		AxisMap: make(AxisMap),
	}
}

//ActionMap groups the handler map types.
type ActionMap struct {
	ActionHandlerMap
	AxisActionHandlerMap
}

// Layers is a slice of Maps. This can be used to combine and use multiple Maps at once.
// Maps are handled in order and if more than one KeyMap has a handler for the same Button
// then their earlier ones have the option to stop propagation to later handlers.
type Layers []Map

// Update calls the Update methods for each Map type in the list.
func (l Layers) Update(am ActionMap) {
	stoppedBtns := StoppedBtnSet{}
	stoppedAxis := StoppedAxisSet{}

	for _, m := range l {
		m.KeyMap.Update(am.ActionHandlerMap, stoppedBtns)
		m.AxisMap.Update(am.AxisActionHandlerMap, stoppedAxis)
	}
}
