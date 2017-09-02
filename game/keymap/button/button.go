package button

import "github.com/hajimehoshi/ebiten"

// KeyMouse combines ebiten's Key and MouseButton.
type KeyMouse int

// FromKey converts Key to a KeyMouse.
func FromKey(key ebiten.Key) KeyMouse {
	return KeyMouse(key)
}

// FromMouse converts MouseButton to a KeyMouse.
func FromMouse(mb ebiten.MouseButton) KeyMouse {
	return KeyMouse(int(ebiten.KeyMax) + int(mb))
}

// IsKey returns true if the KeyMouse is an Key.
func (km KeyMouse) IsKey() bool {
	return int(km) < int(ebiten.KeyMax)
}

// IsMouse returns true if the KeyMouse is a MouseButton.
func (km KeyMouse) IsMouse() bool {
	return int(km) >= int(ebiten.KeyMax)
}

// Key converts KeyMouse to Key. The return value ok is false if it is actually a MouseButton.
func (km KeyMouse) Key() (k ebiten.Key, ok bool) {
	return ebiten.Key(km), km.IsKey()
}

// Mouse converts KeyMouse to MouseButton. The return value ok is false if it is actually a Key.
func (km KeyMouse) Mouse() (mb ebiten.MouseButton, ok bool) {
	return ebiten.MouseButton(int(km) - int(ebiten.KeyMax)), km.IsMouse()
}
