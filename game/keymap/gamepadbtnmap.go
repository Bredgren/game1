package keymap

import (
	"github.com/hajimehoshi/ebiten"
)

// GamepadBtnMap is a bi-directional map connecting gamepad buttons and actions.
type GamepadBtnMap struct {
	btnToAct map[ebiten.GamepadButton]Action
	actToBtn map[Action]ebiten.GamepadButton
}

// NewGamepadBtnMap returns a new, initialized GamepadBtnMap.
func NewGamepadBtnMap() *GamepadBtnMap {
	return &GamepadBtnMap{
		btnToAct: map[ebiten.GamepadButton]Action{},
		actToBtn: map[Action]ebiten.GamepadButton{},
	}
}

// Set associates the given button and action with each other.
func (gm *GamepadBtnMap) Set(b ebiten.GamepadButton, a Action) {
	oldA, oldAok := gm.btnToAct[b]
	oldB, oldBok := gm.actToBtn[a]
	if oldAok {
		delete(gm.actToBtn, oldA)
	}
	if oldBok {
		delete(gm.btnToAct, oldB)
	}

	gm.btnToAct[b] = a
	gm.actToBtn[a] = b
}

// GetButton returns the button associated with the action.
func (gm *GamepadBtnMap) GetButton(a Action) (b ebiten.GamepadButton, ok bool) {
	b, ok = gm.actToBtn[a]
	return
}

// GetAction returns the action associated with the button.
func (gm *GamepadBtnMap) GetAction(b ebiten.GamepadButton) (a Action, ok bool) {
	a, ok = gm.btnToAct[b]
	return
}

// DelButton removes the button and its associated action.
func (gm *GamepadBtnMap) DelButton(b ebiten.GamepadButton) {
	if a, ok := gm.btnToAct[b]; ok {
		delete(gm.btnToAct, b)
		delete(gm.actToBtn, a)
	}
}

// DelAction removes the action and its associated button.
func (gm *GamepadBtnMap) DelAction(a Action) {
	if b, ok := gm.actToBtn[a]; ok {
		delete(gm.actToBtn, a)
		delete(gm.btnToAct, b)
	}
}

// Buttons returns a slice containing all buttons currently in the map.
func (gm *GamepadBtnMap) Buttons() []ebiten.GamepadButton {
	s := make([]ebiten.GamepadButton, 0, len(gm.btnToAct))
	for b := range gm.btnToAct {
		s = append(s, b)
	}
	return s
}

// Actions returns a slice containing all actions currently in the map.
func (gm *GamepadBtnMap) Actions() []Action {
	s := make([]Action, 0, len(gm.actToBtn))
	for a := range gm.actToBtn {
		s = append(s, a)
	}
	return s
}
