package game

import (
	"fmt"
	"strings"
	"time"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/game1/game/keymap"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	frameTime = (time.Second / time.Nanosecond) / ebiten.FPS * time.Nanosecond
)

type gameState int

const (
	mainMenuState gameState = iota
	playState
)

// Game manages the overall state of the game.
type Game struct {
	state         gameState
	showDebugInfo bool
	timeScale     float64
	lastUpdate    time.Time
	camera        *camera.Camera
	background    *background

	actions   keymap.ActionMap
	keyLayers keymap.Layers

	// Fields only for debugging
	lastUpdateTime time.Duration
	lastDrawTime   time.Duration
	lastTimeSample time.Time

	player *player
}

// New creates, initializes, and returns a new Game.
func New(screenWidth, screenHeight int) *Game {
	p := newPlayer()

	cam := camera.New(screenWidth, screenHeight)
	cam.MaxDist = 100
	cam.MaxSpeed = 600
	cam.Ease = geo.EaseOutQuad
	cam.Target = p

	cam.Shaker.Amplitude = 30
	cam.Shaker.Duration = 1 * time.Second
	cam.Shaker.Frequency = 10
	cam.Shaker.Falloff = geo.EaseOutQuad

	g := &Game{
		state:         mainMenuState,
		showDebugInfo: true,
		timeScale:     1.0,
		camera:        cam,
		background:    newBackground(),

		keyLayers: keymap.Layers{},

		player: p,
	}

	g.actions = keymap.ActionMap{
		ActionHandlerMap: keymap.ActionHandlerMap{
			"move left":  g.handelPlayerMoveLeft,
			"move right": g.handelPlayerMoveRight,
			// "jump":       nil,
			// "uppercut":   nil,
			// "slam":       nil,
			// "punch":      nil,
			// "launch":     nil,
			// "pause":      nil,
		},
		AxisActionHandlerMap: keymap.AxisActionHandlerMap{
			"move": g.handlePlayerMove,
			// "punch horizontal": nil,
			// "punch vertical":   nil,
		},
	}

	keyMap := keymap.NewMap()
	setDefaultKeyMap(keyMap)

	g.keyLayers = append(g.keyLayers, keyMap)

	return g
}

// Update the Game by simulating the state by one frame.
func (g *Game) Update() {
	updateStart := time.Now()
	dt := g.dt(updateStart)

	g.keyLayers.Update(g.actions)

	// if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	// 	g.player.pos = g.camera.WorldCoords(geo.VecXYi(ebiten.CursorPosition()))
	// }
	// if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
	// 	g.camera.StartShake()
	// }

	g.player.update(dt)

	g.camera.Update(dt)

	if g.showDebugInfo {
		updateTime := time.Since(updateStart)
		if time.Since(g.lastTimeSample) > time.Second || updateTime > g.lastUpdateTime {
			g.lastUpdateTime = updateTime
		}
	}
}

// Draw the game to the given image. The size of the image shouud be the same as the size
// given to New.
func (g *Game) Draw(dst *ebiten.Image) {
	drawStart := time.Now()

	g.background.Draw(dst, g.camera)

	g.player.draw(dst, g.camera)

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

func (g *Game) handelPlayerMoveLeft(down bool) bool {
	g.player.left = down
	return false
}

func (g *Game) handelPlayerMoveRight(down bool) bool {
	g.player.right = down
	return false
}

func (g *Game) handlePlayerMove(val float64) bool {
	if val != 0 {
		g.player.move = val

	}
	return false
}
