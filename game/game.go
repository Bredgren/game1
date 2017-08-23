package game

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Game manages the overall state of the game.
type Game struct {
	showDebugInfo bool
	timeScale     float64
	lastUpdate    time.Time
	camera        *camera.Camera
	background    *background

	testImg      *ebiten.Image
	opts         *ebiten.DrawImageOptions
	cameraTarget camera.Target
	thing        *thing
}

type thing struct {
	pos geo.Vec
}

func (t *thing) Pos() geo.Vec {
	return t.pos
}

// New creates, initializes, and returns a new Game.
func New(screenWidth, screenHeight int) *Game {
	img, _ := ebiten.NewImage(10, 10, ebiten.FilterNearest)
	img.Fill(color.White)
	cam := camera.New(screenWidth, screenHeight)
	cam.MaxDist = 100
	cam.MaxSpeed = 600
	cam.Ease = geo.EaseOutQuad
	t := &thing{}
	cam.Target = t

	cam.Shaker.Amplitude = 30
	cam.Shaker.Duration = 1 * time.Second
	cam.Shaker.Frequency = 10
	cam.Shaker.Falloff = geo.EaseOutQuad

	return &Game{
		showDebugInfo: true,
		timeScale:     1.0,
		camera:        cam,
		background:    NewBackground(),
		testImg:       img,
		opts:          &ebiten.DrawImageOptions{},
		thing:         t,
	}
}

func (g *Game) Update() {
	dt := g.dt()
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.thing.pos.Add(geo.VecXY(-50, 0).Times(dt.Seconds()))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.thing.pos.Add(geo.VecXY(50, 0).Times(dt.Seconds()))
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.thing.pos.Add(geo.VecXY(0, -50).Times(dt.Seconds()))
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.thing.pos.Add(geo.VecXY(0, 50).Times(dt.Seconds()))
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.thing.pos = g.camera.WorldCoords(geo.VecXYi(ebiten.CursorPosition()))
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.camera.StartShake()
	}
	g.camera.Update(dt)
}

func (g *Game) Draw(dst *ebiten.Image) {
	g.background.Draw(dst, g.camera)

	testPos1 := g.camera.ScreenCoords(geo.VecXY(0, -40))
	g.opts.GeoM.Reset()
	g.opts.GeoM.Translate(testPos1.Floored().XY())
	dst.DrawImage(g.testImg, g.opts)

	testPos2 := g.camera.ScreenCoords(geo.VecXY(40, 0))
	g.opts.GeoM.Reset()
	g.opts.GeoM.Translate(testPos2.Floored().XY())
	dst.DrawImage(g.testImg, g.opts)

	testPos3 := g.camera.ScreenCoords(g.thing.Pos())
	g.opts.GeoM.Reset()
	g.opts.GeoM.Scale(0.5, 0.5)
	g.opts.GeoM.Translate(testPos3.Floored().XY())
	dst.DrawImage(g.testImg, g.opts)

	if g.showDebugInfo {
		g.drawDebugInfo(dst)
	}
}

func (g *Game) dt() time.Duration {
	now := time.Now()
	ns := now.Sub(g.lastUpdate).Nanoseconds()
	scaled := float64(ns) * g.timeScale
	dt := time.Duration(scaled) * time.Nanosecond
	g.lastUpdate = now
	return dt
}

func (g *Game) drawDebugInfo(dst *ebiten.Image) {
	info := []string{
		fmt.Sprintf("FPS %0.2f", ebiten.CurrentFPS()),
		fmt.Sprintf("Time Scale: %0.2f", g.timeScale),
	}
	ebitenutil.DebugPrint(dst, strings.Join(info, "\n"))
}
