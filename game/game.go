package game

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/game1/game/keymap"
	"github.com/Bredgren/game1/game/keymap/button"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	frameTime = (time.Second / time.Nanosecond) / ebiten.FPS * time.Nanosecond
)

const (
	ignore     = "ignore"
	left       = "left"
	right      = "right"
	move       = "move"
	jump       = "jump"
	fullscreen = "fullscreen"
	pause      = "pause"
)

const (
	generalLayer = iota
	remapLayer
	uiLayer
	playerLayer
	numInputLayers
)

// Game manages the overall state of the game.
type Game struct {
	state         gameStateName
	states        map[gameStateName]gameState
	showDebugInfo bool
	timeScale     float64
	lastUpdate    time.Time
	camera        *camera.Camera
	background    *background
	// inputDisabled       bool
	canToggleFullscreen bool

	// keyLabels map[string]*keyLabel

	// actions   keymap.ActionMap
	keymap keymap.Layers

	player *player

	// Fields only for debugging
	lastUpdateTime time.Duration
	lastDrawTime   time.Duration
	lastTimeSample time.Time
}

// New creates, initializes, and returns a new Game.
func New(screenWidth, screenHeight int) *Game {
	p := newPlayer()

	cam := camera.New(screenWidth, screenHeight)
	// cam.MaxDist = 100
	// cam.MaxSpeed = 600
	// cam.Ease = geo.EaseOutExpo
	//
	// cam.MaxDist = 80
	// cam.MaxSpeed = 600
	// cam.Ease = geo.EaseInExpo

	cam.Shaker.Amplitude = 30
	cam.Shaker.Duration = 1 * time.Second
	cam.Shaker.Frequency = 10
	cam.Shaker.Falloff = geo.EaseOutQuad

	// keyOptionsPos := geo.VecXY(100, 100)
	// keyOptionVGap := 2.0
	// keyLabels := []*keyLabel{
	// 	newKeyLabel(left, basicfont.Face7x13),
	// 	newKeyLabel(right, basicfont.Face7x13),
	// 	newKeyLabel(move, basicfont.Face7x13),
	// 	newKeyLabel(jump, basicfont.Face7x13),
	// }
	//
	// for _, kl := range keyLabels {
	// 	kl.bounds.SetTopLeft(keyOptionsPos.XY())
	// 	keyOptionsPos.Y += kl.bounds.H + keyOptionVGap
	// }

	bg := newBackground()

	g := &Game{
		state: intro,
		states: map[gameStateName]gameState{
			intro:    newIntroState(p, screenHeight, cam, bg),
			mainMenu: newMainMenu(p, screenHeight, cam),
		},
		showDebugInfo: true,
		timeScale:     1.0,
		camera:        cam,
		background:    bg,

		// keyLabels: map[string]*keyLabel{},

		keymap: make(keymap.Layers, numInputLayers),

		player: p,
	}

	generalActions := keymap.ButtonHandlerMap{
		pause: func(down bool) bool {
			if down {
				log.Println("pause not implement yet")
			}
			return false
		},
		fullscreen: func(down bool) bool {
			if down && g.canToggleFullscreen {
				ebiten.SetFullscreen(!ebiten.IsFullscreen())
				g.canToggleFullscreen = false
			} else if !down {
				g.canToggleFullscreen = true
			}
			return false
		},
	}
	g.keymap[generalLayer] = keymap.New(generalActions, nil)

	g.keymap[generalLayer].KeyMouse.Set(button.FromKey(ebiten.KeyEscape), pause)
	g.keymap[generalLayer].KeyMouse.Set(button.FromKey(ebiten.KeyF11), fullscreen)
	g.keymap[generalLayer].GamepadBtn.Set(ebiten.GamepadButton7, pause)
	g.keymap[generalLayer].GamepadBtn.Set(ebiten.GamepadButton6, fullscreen)

	// for _, kl := range keyLabels {
	// 	g.keyLabels[kl.name] = kl
	// }

	// g.actions = keymap.ActionMap{
	// 	ActionHandlerMap: keymap.ActionHandlerMap{
	// 		ignore: func(_ bool) bool { return g.inputDisabled },
	// 		left:   g.handlePlayerMoveLeft,
	// 		right:  g.handlePlayerMoveRight,
	// 		jump:   g.handlePlayerJump,
	// 		// "uppercut":   nil,
	// 		// "slam":       nil,
	// 		// "punch":      nil,
	// 		// "launch":     nil,
	// 		fullscreen: func(down bool) bool {
	// 			if down && g.canToggleFullscreen {
	// 				ebiten.SetFullscreen(!ebiten.IsFullscreen())
	// 				g.canToggleFullscreen = false
	// 			} else if !down {
	// 				g.canToggleFullscreen = true
	// 			}
	// 			return false
	// 		},
	// 		pause: func(down bool) bool {
	// 			if down {
	// 				log.Println("pause not implement yet")
	// 			}
	// 			return false
	// 		},
	// 	},
	// 	AxisActionHandlerMap: keymap.AxisActionHandlerMap{
	// 		ignore: func(_ float64) bool { return g.inputDisabled },
	// 		move:   g.handlePlayerMove,
	// 		// "punch horizontal": nil,
	// 		// "punch vertical":   nil,
	// 	},
	// }

	// keyMap := keymap.NewMap()
	// setDefaultKeyMap(keyMap)
	//
	// // Keys that can't be remapped
	// fixedKeyMap := keymap.NewMap()
	// fixedKeyMap.KeyMap[button.FromKey(ebiten.KeyEscape)] = pause
	// fixedKeyMap.KeyMap[button.FromKey(ebiten.KeyF11)] = fullscreen
	// fixedKeyMap.KeyMap[button.FromGamepadButton(ebiten.GamepadButton7)] = pause
	// fixedKeyMap.KeyMap[button.FromGamepadButton(ebiten.GamepadButton6)] = fullscreen
	//
	// // This keymap layer is for disabling all input
	// disableKeyMap := keymap.NewMap()
	// for i := ebiten.Key0; i < ebiten.KeyMax; i++ {
	// 	disableKeyMap.KeyMap[button.FromKey(i)] = ignore
	// }
	// for i := ebiten.GamepadButton0; i < ebiten.GamepadButtonMax; i++ {
	// 	disableKeyMap.KeyMap[button.FromGamepadButton(i)] = ignore
	// }
	// disableKeyMap.KeyMap[button.FromMouseButton(ebiten.MouseButtonLeft)] = ignore
	// disableKeyMap.KeyMap[button.FromMouseButton(ebiten.MouseButtonMiddle)] = ignore
	// disableKeyMap.KeyMap[button.FromMouseButton(ebiten.MouseButtonRight)] = ignore
	// // We don't know how many axes there will be so just do alot :P
	// for i := 0; i < 100; i++ {
	// 	disableKeyMap.AxisMap[i] = ignore
	// }
	//
	// g.keyLayers = append(g.keyLayers, fixedKeyMap)
	// g.keyLayers = append(g.keyLayers, disableKeyMap)
	// g.keyLayers = append(g.keyLayers, keyMap)

	return g
}

// Update the Game by simulating the state by one frame.
func (g *Game) Update() {
	updateStart := time.Now()
	dt := g.dt(updateStart)

	g.keymap.Update()

	// onGround := g.player.canJump
	// g.player.update(dt)
	//
	// // player just contacted the ground
	// if !onGround && g.player.canJump {
	// 	g.camera.StartShake()
	// }

	s := g.states[g.state]
	next := s.nextState()
	if next != g.state {
		log.Println("Change state from", g.state, "to", next)
		s.end()
		s = g.states[next]
		s.begin(g.state)
		g.state = next
	}

	s.update(dt)

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

	// g.player.draw(dst, g.camera)
	//
	// if g.state == mainMenu {
	// 	for _, kl := range g.keyLabels {
	// 		kl.draw(dst, g.camera)
	// 	}
	// }

	g.states[g.state].draw(dst, g.camera)

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
	// Don't want negative dt somehow
	if dt < 0 {
		dt = 0
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

// func (g *Game) handlePlayerMoveLeft(down bool) bool {
// 	g.player.Left = down
// 	g.keyLabels[left].active = down
// 	return false
// }
//
// func (g *Game) handlePlayerMoveRight(down bool) bool {
// 	g.player.Right = down
// 	g.keyLabels[right].active = down
// 	return false
// }
//
// func (g *Game) handlePlayerMove(val float64) bool {
// 	g.player.Move = val
// 	g.keyLabels[move].active = val != 0
// 	return false
// }
//
// func (g *Game) handlePlayerJump(down bool) bool {
// 	g.player.Jump = down
// 	g.keyLabels[jump].active = down
// 	return false
// }
