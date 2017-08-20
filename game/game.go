package game

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Game manages the overall state of the game.
type Game struct {
	showDebugInfo bool
	timeScale     float64
}

// NewGame creates, initializes, and returns a new Game.
func NewGame() *Game {
	return &Game{
		showDebugInfo: true,
		timeScale:     1.0,
	}
}

func (g *Game) Update() {
}

func (g *Game) Draw(dst *ebiten.Image) {
	if g.showDebugInfo {
		g.drawDebugInfo(dst)
	}
}

func (g *Game) drawDebugInfo(dst *ebiten.Image) {
	info := []string{
		fmt.Sprintf("FPS %0.2f", ebiten.CurrentFPS()),
		fmt.Sprintf("Time Scale: %0.2f", g.timeScale),
	}
	ebitenutil.DebugPrint(dst, strings.Join(info, "\n"))
}
