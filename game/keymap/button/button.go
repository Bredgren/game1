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
	return KeyMouse(int(ebiten.KeyMax) + 1 + int(mb))
}

// IsKey returns true if the KeyMouse is an Key.
func (km KeyMouse) IsKey() bool {
	return int(km) <= int(ebiten.KeyMax)
}

// IsMouse returns true if the KeyMouse is a MouseButton.
func (km KeyMouse) IsMouse() bool {
	return int(km) > int(ebiten.KeyMax)
}

// Key converts KeyMouse to Key. The return value ok is false if it is actually a MouseButton.
func (km KeyMouse) Key() (k ebiten.Key, ok bool) {
	return ebiten.Key(km), km.IsKey()
}

// Mouse converts KeyMouse to MouseButton. The return value ok is false if it is actually a Key.
func (km KeyMouse) Mouse() (mb ebiten.MouseButton, ok bool) {
	return ebiten.MouseButton(int(km) - int(ebiten.KeyMax) - 1), km.IsMouse()
}

func (km KeyMouse) String() string {
	if mb, ok := km.Mouse(); ok {
		switch mb {
		case ebiten.MouseButtonLeft:
			return "LMB"
		case ebiten.MouseButtonMiddle:
			return "MMB"
		case ebiten.MouseButtonRight:
			return "RMB"
		}
	}

	k, ok := km.Key()
	if !ok {
		return "N/A"
	}

	switch k {
	case ebiten.Key0:
		return "0"
	case ebiten.Key1:
		return "1"
	case ebiten.Key2:
		return "2"
	case ebiten.Key3:
		return "3"
	case ebiten.Key4:
		return "4"
	case ebiten.Key5:
		return "5"
	case ebiten.Key6:
		return "6"
	case ebiten.Key7:
		return "7"
	case ebiten.Key8:
		return "8"
	case ebiten.Key9:
		return "9"
	case ebiten.KeyA:
		return "A"
	case ebiten.KeyB:
		return "B"
	case ebiten.KeyC:
		return "C"
	case ebiten.KeyD:
		return "D"
	case ebiten.KeyE:
		return "E"
	case ebiten.KeyF:
		return "F"
	case ebiten.KeyG:
		return "G"
	case ebiten.KeyH:
		return "H"
	case ebiten.KeyI:
		return "I"
	case ebiten.KeyJ:
		return "J"
	case ebiten.KeyK:
		return "K"
	case ebiten.KeyL:
		return "L"
	case ebiten.KeyM:
		return "M"
	case ebiten.KeyN:
		return "N"
	case ebiten.KeyO:
		return "O"
	case ebiten.KeyP:
		return "P"
	case ebiten.KeyQ:
		return "Q"
	case ebiten.KeyR:
		return "R"
	case ebiten.KeyS:
		return "S"
	case ebiten.KeyT:
		return "T"
	case ebiten.KeyU:
		return "U"
	case ebiten.KeyV:
		return "V"
	case ebiten.KeyW:
		return "W"
	case ebiten.KeyX:
		return "X"
	case ebiten.KeyY:
		return "Y"
	case ebiten.KeyZ:
		return "Z"
	case ebiten.KeyAlt:
		return "Alt"
	case ebiten.KeyApostrophe:
		return "'"
	case ebiten.KeyBackslash:
		return "\\"
	case ebiten.KeyBackspace:
		return "Backspace"
	case ebiten.KeyCapsLock:
		return "Caps lock"
	case ebiten.KeyComma:
		return ","
	case ebiten.KeyControl:
		return "Ctrl"
	case ebiten.KeyDelete:
		return "Del"
	case ebiten.KeyDown:
		return "Down"
	case ebiten.KeyEnd:
		return "End"
	case ebiten.KeyEnter:
		return "Enter"
	case ebiten.KeyEqual:
		return "="
	case ebiten.KeyEscape:
		return "Esc"
	case ebiten.KeyF1:
		return "F1"
	case ebiten.KeyF2:
		return "F2"
	case ebiten.KeyF3:
		return "F3"
	case ebiten.KeyF4:
		return "F4"
	case ebiten.KeyF5:
		return "F5"
	case ebiten.KeyF6:
		return "F6"
	case ebiten.KeyF7:
		return "F7"
	case ebiten.KeyF8:
		return "F8"
	case ebiten.KeyF9:
		return "F9"
	case ebiten.KeyF10:
		return "F10"
	case ebiten.KeyF11:
		return "F11"
	case ebiten.KeyF12:
		return "F12"
	case ebiten.KeyGraveAccent:
		return "`"
	case ebiten.KeyHome:
		return "Home"
	case ebiten.KeyInsert:
		return "Insert"
	case ebiten.KeyLeft:
		return "Left"
	case ebiten.KeyLeftBracket:
		return "["
	case ebiten.KeyMinus:
		return "-"
	case ebiten.KeyPageDown:
		return "PgDn"
	case ebiten.KeyPageUp:
		return "PgUp"
	case ebiten.KeyPeriod:
		return "."
	case ebiten.KeyRight:
		return "Right"
	case ebiten.KeyRightBracket:
		return "]"
	case ebiten.KeySemicolon:
		return ";"
	case ebiten.KeyShift:
		return "Shift"
	case ebiten.KeySlash:
		return "/"
	case ebiten.KeySpace:
		return "Space"
	case ebiten.KeyTab:
		return "Tab"
	case ebiten.KeyUp:
		return "Up"
	}

	return "-"
}
