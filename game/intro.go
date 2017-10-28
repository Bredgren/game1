package game

import (
	"time"

	"github.com/Bredgren/game1/game/comp"
	"github.com/Bredgren/game1/game/gamestate"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

type introState struct {
	game *Game
}

func newIntroState(game *Game) *introState {
	game.entityState.Mask[game.player] = comp.Position | comp.Camera
	game.entityState.Position[game.player] = geo.Vec0

	game.entityState.Mask[game.camera] |= comp.Shake
	cameraHeight := game.entityState.BoundingBox[game.camera].H
	game.entityState.Position[game.camera] = geo.VecXY(0, -cameraHeight*0.4)
	game.entityState.Shake[game.camera].Shaker.Amplitude = 10
	game.entityState.Shake[game.camera].Shaker.Frequency = 200
	// p.awaken()
	return &introState{
		game: game,
	}
}

func (i *introState) Begin(previousState gamestate.State) {
	// Since introState is the first state, begin is not called
}

func (i *introState) End() {

}

func (i *introState) NextState() gamestate.State {
	// if i.p.awoke() {
	// return mainMenu
	// }
	return gamestate.Intro
}

func (i *introState) Update(dt time.Duration) {
	// i.p.update(dt)
	i.game.shakeUpdate(dt)
}

func (i *introState) Draw(dst *ebiten.Image) {
	cameraBox := i.game.entityState.BoundingBox[i.game.camera]
	cameraPos := i.game.entityState.Position[i.game.camera]
	i.game.background.Draw(dst, cameraBox.Moved(cameraPos.XY()))
	i.game.render(dst)
}
