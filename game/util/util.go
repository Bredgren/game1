package util

import (
	"image/color"

	"github.com/Bredgren/geo"
)

func lerpColorComponent(a, b uint32, t float64) uint8 {
	return uint8(geo.Map(geo.Lerp(float64(a), float64(b), t), 0, 0xffff, 0, 0xff))
}

func LerpColor(c1, c2 color.Color, t float64) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	r := lerpColorComponent(r1, r2, t)
	g := lerpColorComponent(g1, g2, t)
	b := lerpColorComponent(b1, b2, t)
	a := lerpColorComponent(a1, a2, t)
	return color.NRGBA{r, g, b, a}
}
