package game

import "github.com/Bredgren/geo"

type hitbox struct {
	Label    string
	Bounds   geo.Rect
	Callback func(other *hitbox)
	Active   bool
	Owner    interface{}
}
