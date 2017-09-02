package ui

import (
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

// Button is a container that holds one sub-element and images for its states.
type Button struct {
	IdleImg     *ebiten.Image
	HoverImg    *ebiten.Image
	Element     WeightedDrawer
	IdleAnchor  Anchor
	HoverAnchor Anchor
	Wt          float64
	Hover       bool
	OnClick     func()
	lastRect    geo.Rect
}

// Update sets b.Hover according to the current mouse position and the position that
// the last call to Draw put the button at.
func (b *Button) Update() {
	mousePos := geo.VecXY(geo.I2F2(ebiten.CursorPosition()))
	b.Hover = b.lastRect.CollidePoint(mousePos.XY())
}

// Draw draws the Button to the image within bounds.
func (b *Button) Draw(dst *ebiten.Image, bounds geo.Rect) {
	opts := ebiten.DrawImageOptions{}
	var imgBounds geo.Rect
	var topLeft geo.Vec
	if b.Hover {
		imgBounds = geo.RectWH(geo.I2F2(b.HoverImg.Size()))
		topLeft = b.HoverAnchor.TopLeft(imgBounds, bounds)
		opts.GeoM.Translate(topLeft.XY())
		dst.DrawImage(b.HoverImg, &opts)
	} else {
		imgBounds = geo.RectWH(geo.I2F2(b.IdleImg.Size()))
		topLeft = b.IdleAnchor.TopLeft(imgBounds, bounds)
		opts.GeoM.Translate(topLeft.XY())
		dst.DrawImage(b.IdleImg, &opts)
	}
	imgBounds.SetTopLeft(topLeft.XY())
	b.lastRect = imgBounds

	if b.Element != nil {
		b.Element.Draw(dst, imgBounds)
	}
}

// Weight returns the relative weight for allocating space within a container.
func (b *Button) Weight() float64 {
	return b.Wt
}
