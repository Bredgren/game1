package game

import (
	"fmt"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Game manages the overall state of the game.
type Game struct {
	showDebugInfo bool
	timeScale     float64
	lastUpdate    time.Time
}

// NewGame creates, initializes, and returns a new Game.
func NewGame() *Game {
	return &Game{
		showDebugInfo: true,
		timeScale:     1.0,
	}
}

func (g *Game) Update() {
	dt := g.dt()
	_ = dt
}

func (g *Game) Draw(dst *ebiten.Image) {
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
