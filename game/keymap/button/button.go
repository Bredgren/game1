package button

import "github.com/hajimehoshi/ebiten"

// Button combines ebiten's Key, GamepadButton, and MouseButton.
type Button int

// FromKey converts ebiten.Key to a Button.
func FromKey(key ebiten.Key) Button {
	return Button(key)
}

// FromGamepadButton converts ebiten.GamepadButton to a Button.
func FromGamepadButton(gb ebiten.GamepadButton) Button {
	return Button(int(ebiten.KeyMax) + int(gb))
}

// FromMouseButton converts ebiten.MouseButton to a Button.
func FromMouseButton(mb ebiten.MouseButton) Button {
	return Button(int(ebiten.KeyMax) + int(ebiten.GamepadButtonMax) + int(mb))
}

// IsKey returns true if the Button is an ebiten.Key.
func (b Button) IsKey() bool {
	return int(b) <= int(ebiten.KeyMax)
}

// IsGamepadButton returns true if the Button is an ebiten.GamepadButton.
func (b Button) IsGamepadButton() bool {
	return int(ebiten.KeyMax) < int(b) && int(b) <= int(ebiten.KeyMax)+int(ebiten.GamepadButtonMax)
}

// IsMouseButton returns true if the Button is an ebiten.GamepadButton.
func (b Button) IsMouseButton() bool {
	return int(ebiten.KeyMax)+int(ebiten.GamepadButtonMax) < int(b)
}

// Key converts Button to ebiten.Key. The return value ok is false if it is actually a
// GamepadButton or MouseButton.
func (b Button) Key() (k ebiten.Key, ok bool) {
	return ebiten.Key(b), b.IsKey()
}

// GamepadButton converts Button to ebiten.GamepadButton. The return value ok is false
// if it is actually a Key or MouseButton.
func (b Button) GamepadButton() (gb ebiten.GamepadButton, ok bool) {
	return ebiten.GamepadButton(b), b.IsGamepadButton()
}

// MouseButton converts Button to ebiten.MouseButton. The return value ok is false if it
// is actually a Key or GamepadButton.
func (b Button) MouseButton() (mb ebiten.MouseButton, ok bool) {
	return ebiten.MouseButton(b), b.IsMouseButton()
}
