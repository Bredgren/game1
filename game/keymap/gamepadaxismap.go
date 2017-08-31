package keymap

// GamepadAxisMap is a bi-directional map connecting gamepad axes and actions.
type GamepadAxisMap struct {
	axisToAct map[int]Action
	actToaxis map[Action]int
}

// NewGamepadAxisMap returns a new, initialized GamepadAxisMap.
func NewGamepadAxisMap() *GamepadAxisMap {
	return &GamepadAxisMap{
		axisToAct: map[int]Action{},
		actToaxis: map[Action]int{},
	}
}

// Set associates the given axis and action with each other.
func (gm *GamepadAxisMap) Set(ax int, a Action) {
	oldA, oldAok := gm.axisToAct[ax]
	oldAx, oldAxok := gm.actToaxis[a]
	if oldAok {
		delete(gm.actToaxis, oldA)
	}
	if oldAxok {
		delete(gm.axisToAct, oldAx)
	}

	gm.axisToAct[ax] = a
	gm.actToaxis[a] = ax
}

// GetAxis returns the axis associated with the action.
func (gm *GamepadAxisMap) GetAxis(a Action) (ax int, ok bool) {
	ax, ok = gm.actToaxis[a]
	return
}

// GetAction returns the action associated with the button.
func (gm *GamepadAxisMap) GetAction(b int) (a Action, ok bool) {
	a, ok = gm.axisToAct[b]
	return
}

// DelAxis removes the axis and its associated action.
func (gm *GamepadAxisMap) DelAxis(ax int) {
	if a, ok := gm.axisToAct[ax]; ok {
		delete(gm.axisToAct, ax)
		delete(gm.actToaxis, a)
	}
}

// DelAction removes the action and its associated axis.
func (gm *GamepadAxisMap) DelAction(a Action) {
	if ax, ok := gm.actToaxis[a]; ok {
		delete(gm.actToaxis, a)
		delete(gm.axisToAct, ax)
	}
}

// Axes returns a slice containing all axes currently in the map.
func (gm *GamepadAxisMap) Axes() []int {
	s := make([]int, 0, len(gm.axisToAct))
	for ax := range gm.axisToAct {
		s = append(s, ax)
	}
	return s
}

// Actions returns a slice containing all actions currently in the map.
func (gm *GamepadAxisMap) Actions() []Action {
	s := make([]Action, 0, len(gm.actToaxis))
	for a := range gm.actToaxis {
		s = append(s, a)
	}
	return s
}
