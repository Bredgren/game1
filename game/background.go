package game

import (
	"image"
	"image/color"
	"math"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/game1/game/util"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

type background struct {
	skycolor1      color.Color
	skyclor1       color.Color
	maxHeight      float64
	clouds         []*ebiten.Image
	cloudScaleMin  geo.Vec
	cloudScaleMax  geo.Vec
	cloudMinHight  float64
	cloudThickness float64
	padding        float64

	groundColor color.Color
	groundImg   *ebiten.Image

	// cloudTest *ebiten.Image
	// cloudW    int
	// cloudH    int
	// cloudX    float64
	// cloudY    float64
}

func NewBackground() *background {
	b := &background{
		skycolor1: color.NRGBA{255, 140, 68, 255},
		skyclor1:  color.NRGBA{0, 0, 10, 255},
		maxHeight: 700, // When the background becomes dark
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
		cloudScaleMin:  geo.VecXY(0.5, 0.5),
		cloudScaleMax:  geo.VecXY(10, 2),
		cloudMinHight:  150, // Lowest a cloud can be
		cloudThickness: 700, // Vertical size of the area a cloud can be

		groundColor: color.NRGBA{60, 60, 60, 255},

		// cloudTest: makeCloud(100, 100, 0, 0),
		// cloudW:    100,
		// cloudH:    100,
		// cloudX:    0,
		// cloudY:    0,
	}

	// To make sure clouds don't suddenly appear on screen we select a padding equal to
	// the maximum size of a cloud and we will draw clouds off screen up to that distance.
	for _, cloud := range b.clouds {
		w, h := geo.I2F2(cloud.Size())
		diagonal := math.Hypot(w, h)
		max := diagonal * math.Max(b.cloudScaleMax.X, b.cloudScaleMax.Y)
		b.padding = math.Max(max, b.padding)
	}

	b.groundImg, _ = ebiten.NewImage(1, 1, ebiten.FilterNearest)
	b.groundImg.Fill(b.groundColor)

	return b
}

func (b *background) Draw(dst *ebiten.Image, cam *camera.Camera) {
	height := geo.Clamp(-cam.Center().Y, 0, b.maxHeight) / b.maxHeight
	dst.Fill(util.LerpColor(b.skycolor1, b.skyclor1, height))

	// b.cloudFinder(dst, cam)

	b.drawClouds(dst, cam)
	b.drawGround(dst, cam)
}

func (b *background) drawClouds(dst *ebiten.Image, cam *camera.Camera) {
	topLeft := cam.WorldCoords(geo.VecXY(-b.padding, -b.padding))
	screenSize := geo.VecXYi(dst.Size())
	bottomRight := cam.WorldCoords(geo.VecXY(b.padding, b.padding).Plus(screenSize))
	area := geo.RectCornersVec(topLeft, bottomRight)

	// cutoff is used to create some cloudless gaps
	cutoff := 0.6
	// gap is the spacing between clouds, lowering creates more/thicker clouds
	gap := 50

	opts := ebiten.DrawImageOptions{}
	// Round to nearest mutliple of gap becuase we need the x values to be consistent.
	left := float64(((int(area.Left()) + gap/2) / gap) * gap)
	right := float64(((int(area.Right()) + gap/2) / gap) * gap)
	for x := left; x < right; x += float64(gap) {
		noise := geo.Perlin(x*0.5, 0.12345, 0.678901)
		if noise < cutoff {
			continue
		}
		y := geo.Map(noise, cutoff, 1, -b.cloudMinHight, -b.cloudMinHight-b.cloudThickness)

		noise2 := geo.Perlin(x, 0.678901, 0.12345)

		opts.GeoM.Reset()

		opts.GeoM.Rotate(noise2 * 2 * math.Pi)

		xScale := geo.Map(noise2, 0, 1, b.cloudScaleMin.X, b.cloudScaleMax.X)
		yScale := geo.Map(noise2, 0, 1, b.cloudScaleMin.Y, b.cloudScaleMax.Y)
		opts.GeoM.Scale(xScale, yScale)

		pos := cam.ScreenCoords(geo.VecXY(x, y))
		opts.GeoM.Translate(pos.XY())

		cloudIndex := int(math.Floor(noise2 * float64(len(b.clouds)+1)))
		dst.DrawImage(b.clouds[cloudIndex], &opts)
	}
}

func (b *background) drawGround(dst *ebiten.Image, cam *camera.Camera) {
	dstSize := geo.VecXYi(dst.Size())
	cameraBottomRight := cam.WorldCoords(dstSize)
	if cameraBottomRight.Y < 0 {
		return // Ground is not visible
	}

	cameraTopLeft := cam.WorldCoords(geo.Vec0)
	groundTopLeft := geo.VecXY(cameraTopLeft.X, 0)
	topLeft := cam.ScreenCoords(groundTopLeft)

	groundArea := geo.RectCornersVec(topLeft, dstSize)

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Scale(groundArea.Size())
	opts.GeoM.Translate(groundArea.TopLeft())

	dst.DrawImage(b.groundImg, &opts)
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
