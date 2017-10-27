package gamestate

import (
	"time"

	"github.com/hajimehoshi/ebiten"
)

type State int

const (
	Intro State = iota
	MainMenu
	Play
)

type GameState interface {
	Begin(previousState State)
	End()
	NextState() State
	Update(dt time.Duration)
	Draw(dst *ebiten.Image)
}
