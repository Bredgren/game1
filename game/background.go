package game

import (
	"image"
	"image/color"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/game1/game/util"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

type background struct {
	color1        color.Color
	color2        color.Color
	maxHeight     float64
	clouds        []*ebiten.Image
	cloudScaleMin geo.Vec
	cloudScaleMax geo.Vec

	// cloudTest *ebiten.Image
	// cloudW    int
	// cloudH    int
	// cloudX    float64
	// cloudY    float64
}

func NewBackground() *background {
	return &background{
		color1:    color.NRGBA{255, 140, 68, 255},
		color2:    color.NRGBA{0, 0, 10, 255},
		maxHeight: 500,
		clouds: []*ebiten.Image{
			// These were found manually with the cloudFinder method below
			makeCloud(192, 90, 2.24, 0.77),
			makeCloud(80, 80, -6.38, -0.55),
			makeCloud(73, 64, -8.57, -10.46),
			makeCloud(84, 97, -10.27, -1.56),
			makeCloud(147, 124, -14.71, 2.8),
			makeCloud(140, 153, 2.85, 14.22),
			makeCloud(130, 157, 13.3, 21.91),
			makeCloud(105, 65, 13.66, 23.44),
			makeCloud(94, 184, 27.74, 28.81),
			makeCloud(104, 84, 34.29, 32.91),
		},
		cloudScaleMin: geo.VecXY(0.5, 1),
		cloudScaleMax: geo.VecXY(5, 3),

		// cloudTest: makeCloud(100, 100, 0, 0),
		// cloudW:    100,
		// cloudH:    100,
		// cloudX:    0,
		// cloudY:    0,
	}
}

func (b *background) Draw(dst *ebiten.Image, cam *camera.Camera) {
	height := geo.Clamp(-cam.Center().Y, 0, b.maxHeight) / b.maxHeight
	dst.Fill(util.LerpColor(b.color1, b.color2, height))

	// b.cloudFinder(dst, cam)

	pos := cam.ScreenCoords(geo.VecXY(-100, -100))
	opts := ebiten.DrawImageOptions{}
	// opts.GeoM.Rotate(math.Pi / 3)
	// opts.GeoM.Scale(3, 1)
	opts.GeoM.Translate(pos.XY())
	for _, cloud := range b.clouds {
		dst.DrawImage(cloud, &opts)
		w, _ := cloud.Size()
		opts.GeoM.Translate(float64(w+10), 0)
	}
}

func makeCloud(width, height int, xOff, yOff float64) *ebiten.Image {
	pix := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < width*height; i++ {
		x, y := float64(i%width), float64(i/width)
		// By filtering out values < 0.5 then rescaling between 0 and 1 we get more isolated
		// clouds. This way we can more easily contain a cloud in a rectangle and not cutoff
		// others at the edges of the rectangle.
		val := geo.Map(filter(geo.PerlinOctave(x*0.01+xOff, y*0.01+yOff, 0, 3, 0.5), 0.5), 0.5, 1, 0, 1)
		pix.Pix[4*i+3] = uint8(0xff * val)
	}

	img, _ := ebiten.NewImage(width, height, ebiten.FilterLinear)
	img.ReplacePixels(pix.Pix)
	return img
}

func filter(n, min float64) float64 {
	if n > min {
		return n
	}
	return 0
}

// func (b *background) cloudFinder(dst *ebiten.Image, cam *camera.Camera) {
// 	speed := 1.0
// 	if ebiten.IsKeyPressed(ebiten.KeyShift) {
// 		speed = 5
// 	}
//
// 	change := false
// 	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
// 		b.cloudX -= 0.01 * speed
// 		change = true
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyRight) {
// 		b.cloudX += 0.01 * speed
// 		change = true
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyUp) {
// 		b.cloudY -= 0.01 * speed
// 		change = true
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyDown) {
// 		b.cloudY += 0.01 * speed
// 		change = true
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyJ) {
// 		b.cloudW -= 1 * int(speed)
// 		change = true
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyL) {
// 		b.cloudW += 1 * int(speed)
// 		change = true
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyI) {
// 		b.cloudH -= 1 * int(speed)
// 		change = true
// 	}
// 	if ebiten.IsKeyPressed(ebiten.KeyK) {
// 		b.cloudH += 1 * int(speed)
// 		change = true
// 	}
//
// 	if change {
// 		log.Println(b.cloudW, b.cloudH, b.cloudX, b.cloudY)
// 		b.cloudTest = makeCloud(b.cloudW, b.cloudH, b.cloudX, b.cloudY)
// 	}
//
// 	pos := cam.ScreenCoords(geo.VecXY(-50, 0))
// 	opts := ebiten.DrawImageOptions{}
// 	opts.GeoM.Scale(1, 1)
// 	opts.GeoM.Translate(pos.XY())
// 	dst.DrawImage(b.cloudTest, &opts)
// }
