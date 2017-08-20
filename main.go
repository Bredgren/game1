package main

import (
	"log"

	"github.com/Bredgren/game1/game"
	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth = 300
	screenHight = 200
)

var theGame *game.Game

func update(screen *ebiten.Image) error {
	theGame.Update()

	if ebiten.IsRunningSlowly() {
		return nil
	}

	theGame.Draw(screen)

	return nil
}

func init() {
	theGame = game.NewGame()
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHight, 3, "Game Title"); err != nil {
		log.Fatal(err)
	}
}
