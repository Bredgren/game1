package ui

import (
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

// VerticalContainer is a container that orders each of its elements vertically. Their
// widths will be the full width of the bounding rectangle given to Draw and their heights
// will be determined by their Weights relative to the other elements.
type VerticalContainer struct {
	Elements []WeightedDrawer
	Wt       float64
}

// Draw draws all of the container's elements to dst within bounds.
func (v *VerticalContainer) Draw(dst *ebiten.Image, bounds geo.Rect) {
	totalWeight := 0.0
	for _, e := range v.Elements {
		totalWeight += e.Weight()
	}
	heights := make([]float64, len(v.Elements))
	for i, e := range v.Elements {
		heights[i] = e.Weight() / totalWeight * bounds.H
	}
	subBounds := bounds
	for i, h := range heights {
		subBounds.H = h
		v.Elements[i].Draw(dst, subBounds)
		subBounds.Y += h
	}
}

// Weight returns the relative weight of the container so that containers may be nested.
func (v *VerticalContainer) Weight() float64 {
	return v.Wt
}

// HorizontalContainer is a container that orders each of its elements horizontally. Their
// heights will be the full height of the bounding rectangle given to Draw and their widths
// will be determined by their Weights relative to the other elements.
type HorizontalContainer struct {
	Elements []WeightedDrawer
	Wt       float64
}

// Draw draws all of the container's elements to dst within bounds.
func (h *HorizontalContainer) Draw(dst *ebiten.Image, bounds geo.Rect) {
	totalWeight := 0.0
	for _, e := range h.Elements {
		totalWeight += e.Weight()
	}
	widths := make([]float64, len(h.Elements))
	for i, e := range h.Elements {
		widths[i] = e.Weight() / totalWeight * bounds.W
	}
	subBounds := bounds
	for i, w := range widths {
		subBounds.W = w
		h.Elements[i].Draw(dst, subBounds)
		subBounds.X += w
	}
}

// Weight returns the relative weight of the container so that containers may be nested.
func (h *HorizontalContainer) Weight() float64 {
	return h.Wt
}
