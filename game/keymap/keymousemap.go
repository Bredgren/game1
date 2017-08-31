package keymap

import "github.com/Bredgren/game1/game/keymap/button"

// KeyMouseMap is a bi-directional map connecting keyboard/mouse buttons and actions.
type KeyMouseMap struct {
	btnToAct map[button.KeyMouse]Action
	actToBtn map[Action]button.KeyMouse
}

// NewKeyMouseMap returns a new, initialized KeyMouseMap.
func NewKeyMouseMap() *KeyMouseMap {
	return &KeyMouseMap{
		btnToAct: map[button.KeyMouse]Action{},
		actToBtn: map[Action]button.KeyMouse{},
	}
}

// Set associates the given button and action with each other.
func (km *KeyMouseMap) Set(b button.KeyMouse, a Action) {
	oldA, oldAok := km.btnToAct[b]
	oldB, oldBok := km.actToBtn[a]
	if oldAok {
		delete(km.actToBtn, oldA)
	}
	if oldBok {
		delete(km.btnToAct, oldB)
	}

	km.btnToAct[b] = a
	km.actToBtn[a] = b
}

// GetButton returns the button associated with the action.
func (km *KeyMouseMap) GetButton(a Action) (b button.KeyMouse, ok bool) {
	b, ok = km.actToBtn[a]
	return
}

// GetAction returns the action associated with the button.
func (km *KeyMouseMap) GetAction(b button.KeyMouse) (a Action, ok bool) {
	a, ok = km.btnToAct[b]
	return
}

// DelButton removes the button and its associated action.
func (km *KeyMouseMap) DelButton(b button.KeyMouse) {
	if a, ok := km.btnToAct[b]; ok {
		delete(km.btnToAct, b)
		delete(km.actToBtn, a)
	}
}

// DelAction removes the action and its associated button.
func (km *KeyMouseMap) DelAction(a Action) {
	if b, ok := km.actToBtn[a]; ok {
		delete(km.actToBtn, a)
		delete(km.btnToAct, b)
	}
}

// Buttons returns a slice containing all buttons currently in the map.
func (km *KeyMouseMap) Buttons() []button.KeyMouse {
	s := make([]button.KeyMouse, 0, len(km.btnToAct))
	for b := range km.btnToAct {
		s = append(s, b)
	}
	return s
}

// Actions returns a slice containing all actions currently in the map.
func (km *KeyMouseMap) Actions() []Action {
	s := make([]Action, 0, len(km.actToBtn))
	for a := range km.actToBtn {
		s = append(s, a)
	}
	return s
}
