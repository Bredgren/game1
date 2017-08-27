package main

import (
	"log"

	"github.com/Bredgren/game1/game"
	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 600
	screenHeight = 400
)

var canChangeFullscreen = true

var theGame *game.Game

func togglFullscreen() {
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		if canChangeFullscreen {
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
			canChangeFullscreen = false
		}
	} else {
		canChangeFullscreen = true
	}
}

func update(screen *ebiten.Image) error {
	togglFullscreen()

	theGame.Update()

	if ebiten.IsRunningSlowly() {
		return nil
	}

	theGame.Draw(screen)

	return nil
}

func init() {
	theGame = game.New(screenWidth, screenHeight)
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Game Title"); err != nil {
		log.Fatal(err)
	}
}
