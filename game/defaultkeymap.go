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

var defaultPlayKeyMap = keymap.Map{
	KeyMap: keymap.KeyMap{
		button.FromKey(ebiten.KeyA):                     "move left",
		button.FromKey(ebiten.KeyD):                     "move right",
		button.FromKey(ebiten.KeySpace):                 "jump",
		button.FromKey(ebiten.KeyW):                     "uppercut",
		button.FromKey(ebiten.KeyS):                     "slam",
		button.FromMouseButton(ebiten.MouseButtonLeft):  "punch",
		button.FromMouseButton(ebiten.MouseButtonRight): "launch",
		button.FromKey(ebiten.KeyEscape):                "pause",

		button.FromGamepadButton(ebiten.GamepadButton0): "jump",
		button.FromGamepadButton(ebiten.GamepadButton2): "slam",
		button.FromGamepadButton(ebiten.GamepadButton3): "uppercut",
		button.FromGamepadButton(ebiten.GamepadButton5): "launch",
		button.FromGamepadButton(ebiten.GamepadButton7): "pause",
	},
	AxisMap: keymap.AxisMap{
		0: "move",
		2: "punch horizontal",
		3: "punch vertical",
	},
}

func setDefaultKeyMap(km keymap.Map) {
	for k := range km.KeyMap {
		delete(km.KeyMap, k)
	}
	for k := range km.AxisMap {
		delete(km.AxisMap, k)
	}

	for k, v := range defaultPlayKeyMap.KeyMap {
		km.KeyMap[k] = v
	}
	for k, v := range defaultPlayKeyMap.AxisMap {
		km.AxisMap[k] = v
	}
}
