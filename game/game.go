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

const (
	frameTime = (time.Second / time.Nanosecond) / 60 * time.Nanosecond
)

// Game manages the overall state of the game.
type Game struct {
	showDebugInfo bool
	timeScale     float64
	lastUpdate    time.Time
	camera        *camera.Camera
	background    *background

	// Fields only for debugging
	lastUpdateTime time.Duration
	lastDrawTime   time.Duration
	lastTimeSample time.Time

	// temporary stuff
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
	updateStart := time.Now()
	dt := g.dt(updateStart)

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

	if g.showDebugInfo {
		updateTime := time.Since(updateStart)
		if time.Since(g.lastTimeSample) > time.Second || updateTime > g.lastUpdateTime {
			g.lastUpdateTime = updateTime
		}
	}
}

func (g *Game) Draw(dst *ebiten.Image) {
	drawStart := time.Now()

	g.background.Draw(dst, g.camera)

	testPos := g.camera.ScreenCoords(g.thing.Pos())
	g.opts.GeoM.Reset()
	g.opts.GeoM.Scale(0.5, 0.5)
	g.opts.GeoM.Translate(testPos.Floored().XY())
	dst.DrawImage(g.testImg, g.opts)

	if g.showDebugInfo {
		drawTime := time.Since(drawStart)
		if time.Since(g.lastTimeSample) > time.Second || drawTime > g.lastDrawTime {
			g.lastDrawTime = drawTime
		}
		if time.Since(g.lastTimeSample) > time.Second {
			g.lastTimeSample = drawStart
		}

		g.drawDebugInfo(dst)
	}
}

func (g *Game) dt(now time.Time) time.Duration {
	ns := now.Sub(g.lastUpdate).Nanoseconds()
	scaled := float64(ns) * g.timeScale
	dt := time.Duration(scaled) * time.Nanosecond
	g.lastUpdate = now
	// Cap dt at twice the frame time to prevent large jumps
	maxDt := 2 * frameTime
	if dt > maxDt {
		dt = maxDt
	}
	return dt
}

func (g *Game) drawDebugInfo(dst *ebiten.Image) {
	info := []string{
		fmt.Sprintf("Update+Draw: %0.2f+%0.2f = %0.2f/%0.2f %0.2f%%",
			g.lastUpdateTime.Seconds()*1000, g.lastDrawTime.Seconds()*1000,
			(g.lastUpdateTime+g.lastDrawTime).Seconds()*1000, frameTime.Seconds()*1000,
			(g.lastUpdateTime+g.lastDrawTime).Seconds()/frameTime.Seconds()*100),
		fmt.Sprintf("FPS %0.2f", ebiten.CurrentFPS()),
		fmt.Sprintf("Time Scale: %0.2f", g.timeScale),
	}
	ebitenutil.DebugPrint(dst, strings.Join(info, "\n"))
}
