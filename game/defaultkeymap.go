package game

import (
	"github.com/Bredgren/game1/game/keymap"
	"github.com/Bredgren/game1/game/keymap/button"
	"github.com/hajimehoshi/ebiten"
)

// Xbox Elite controller buttons:
// GamepadButton0  A
// GamepadButton1  B
// GamepadButton2  X
// GamepadButton3  Y
// GamepadButton4  LB
// GamepadButton5  RB
// GamepadButton6  Select
// GamepadButton7  Start
// GamepadButton8  L Stick
// GamepadButton9  R Stick
// GamepadButton10 Up D-Pad
// GamepadButton11 Right D-Pad
// GamepadButton12 Down D-Pad
// GamepadButton13 Left D-Pad
//
// Axis0 Left X
// Axis1 Left Y
// Axis2 Right X
// Axis3 Right Y
// Axis4 Left Trigger (-1 default)
// Axis5 Right Trigger (-1 default)

func setDefaultKeyMap(km *keymap.KeyMap) {
	for _, btn := range km.KeyMouse.Buttons() {
		km.KeyMouse.DelButton(btn)
	}
	for _, btn := range km.GamepadBtn.Buttons() {
		km.GamepadBtn.DelButton(btn)
	}
	for _, axis := range km.GamepadAxis.Axes() {
		km.GamepadAxis.DelAxis(axis)
	}

	km.KeyMouse.Set(button.FromKey(ebiten.KeyA), left)
	km.KeyMouse.Set(button.FromKey(ebiten.KeyD), right)
	km.KeyMouse.Set(button.FromKey(ebiten.KeySpace), jump)
	km.KeyMouse.Set(button.FromKey(ebiten.KeyW), uppercut)
	km.KeyMouse.Set(button.FromKey(ebiten.KeyS), slam)
	km.KeyMouse.Set(button.FromMouse(ebiten.MouseButtonLeft), punch)
	km.KeyMouse.Set(button.FromMouse(ebiten.MouseButtonRight), launch)

	km.GamepadBtn.Set(ebiten.GamepadButton13, left)
	km.GamepadBtn.Set(ebiten.GamepadButton11, right)
	km.GamepadBtn.Set(ebiten.GamepadButton0, jump)
	km.GamepadBtn.Set(ebiten.GamepadButton3, uppercut)
	km.GamepadBtn.Set(ebiten.GamepadButton5, launch)
	km.GamepadBtn.Set(ebiten.GamepadButton2, slam)
	km.GamepadBtn.Set(ebiten.GamepadButton7, pause)

	km.GamepadAxis.Set(0, move)
	km.GamepadAxis.Set(2, punchH)
	km.GamepadAxis.Set(3, punchV)
}

var defaultKeyMap *keymap.KeyMap

func init() {
	defaultKeyMap = keymap.New(nil, nil)
	setDefaultKeyMap(defaultKeyMap)
}
