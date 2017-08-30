package game

import (
	"time"

	"github.com/Bredgren/game1/game/camera"
	"github.com/hajimehoshi/ebiten"
)

type gameStateName int

const (
	intro gameStateName = iota
	mainMenu
	play
)

type gameState interface {
	begin(previousState gameStateName)
	end()
	nextState() gameStateName
	update(dt time.Duration)
	draw(dst *ebiten.Image, cam *camera.Camera)
}
