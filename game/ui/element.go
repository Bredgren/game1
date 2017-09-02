package ui

import (
	"image/color"

	"golang.org/x/image/font"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

// Text is an element that contains text.
type Text struct {
	Anchor Anchor
	Text   string
	Color  color.Color
	Face   font.Face
	Wt     float64
}

// Draw draws the Text to the image within bounds.
func (t *Text) Draw(dst *ebiten.Image, bounds geo.Rect) {
	b, _ := font.BoundString(t.Face, t.Text)
	width := float64((b.Max.X - b.Min.X).Ceil())
	height := float64((b.Max.Y - b.Min.Y).Ceil())
	offset := (t.Face.Metrics().Height - t.Face.Metrics().Descent).Floor() // - 1
	textBounds := geo.RectWH(width, height)

	textBounds.SetTopLeft(t.Anchor.TopLeft(textBounds, bounds).XY())
	text.Draw(dst, t.Text, t.Face, int(textBounds.X), int(textBounds.Y)+offset, t.Color)
}

// Weight returns the relative weight for allocating space within a container.
func (t *Text) Weight() float64 {
	return t.Wt
}
