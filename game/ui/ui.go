package ui

import (
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

// Weighted provides the function Weight which returns the relative space an element would
// like to take up.
type Weighted interface {
	Weight() float64
}

// Drawer provides the function Draw which draws an element to an image within a specified
// bounding rectangle.
type Drawer interface {
	Draw(dst *ebiten.Image, bounds geo.Rect)
}

// WeightedDrawer combines the Weighted and Drawer interfaces.
type WeightedDrawer interface {
	Weighted
	Drawer
}
