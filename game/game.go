package game

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Bredgren/game1/game/comp"
	"github.com/Bredgren/game1/game/gamestate"
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
	move       = "left/right"
	jump       = "jump"
	uppercut   = "uppercut"
	slam       = "slam"
	punch      = "punch"
	launch     = "launch"
	punchH     = "punch horizontal"
	punchV     = "punch vertical"
	fullscreen = "fullscreen"
	pause      = "pause"
	leftClick  = "left click"
)

const (
	generalLayer   = iota
	remapLayer     // Handles key remapping
	leftClickLayer // Handles UI left clicks
	uiLayer        // Handles other UI keys
	playerLayer    // Handles player controls
	numInputLayers
)

// Game manages the overall state of the game.
type Game struct {
	state         gamestate.State
	states        map[gamestate.State]gamestate.GameState
	showDebugInfo bool
	timeScale     float64
	lastUpdate    time.Time
	background    *background

	canToggleFullscreen bool
	canTogglePause      bool

	keymap keymap.Layers
	input  input

	entityState *state
	camera      entity
	player      entity

	// Fields only for debugging
	lastUpdateTime time.Duration
	lastDrawTime   time.Duration
	lastTimeSample time.Time

	//
	// test        map[string]*sprite.Desc
	// testSprites []sprite.Sprite
	// counter     time.Duration
}

// New creates, initializes, and returns a new Game.
func New(screenWidth, screenHeight int) *Game {
	// cam := camera.New(screenWidth, screenHeight)
	// // cam.MaxDist = 100
	// // cam.MaxSpeed = 600
	// // cam.Ease = geo.EaseOutExpo
	// //
	// // cam.MaxDist = 80
	// // cam.MaxSpeed = 600
	// // cam.Ease = geo.EaseInExpo
	//
	// cam.Shaker.Amplitude = 30
	// cam.Shaker.Duration = 1 * time.Second
	// cam.Shaker.Frequency = 10
	// cam.Shaker.Falloff = geo.EaseOutQuad

	es := newState(100)

	camera, err := es.newEntity()
	if err != nil {
		log.Fatal(err)
	}

	es.Mask[camera] = comp.Position | comp.BoundingBox
	es.Position[camera] = geo.Vec0
	es.BoundingBox[camera] = geo.RectXYWH(
		float64(-screenWidth/2), float64(-screenHeight/2), float64(screenWidth), float64(screenHeight))

	player, err := es.newEntity()
	if err != nil {
		log.Fatal(err)
	}

	es.Mask[player] = comp.Camera
	es.Camera[player] = camera

	bg := newBackground()

	g := &Game{
		state:         gamestate.Intro,
		showDebugInfo: true,
		timeScale:     1.0,
		background:    bg,

		keymap: make(keymap.Layers, numInputLayers),

		entityState: es,
		player:      player,
		camera:      camera,

		// test: map[string]*sprite.Desc{},
	}

	generalActions := keymap.ButtonHandlerMap{
		pause: func(down bool) bool {
			if down && g.canTogglePause {
				log.Println("pause not implement yet")
				g.canTogglePause = false
			} else if !down {
				g.canTogglePause = true
			}
			return true
		},
		fullscreen: func(down bool) bool {
			if down && g.canToggleFullscreen {
				ebiten.SetFullscreen(!ebiten.IsFullscreen())
				g.canToggleFullscreen = false
			} else if !down {
				g.canToggleFullscreen = true
			}
			return true
		},
	}

	g.keymap[generalLayer] = keymap.New(generalActions, nil)
	g.keymap[generalLayer].KeyMouse.Set(button.FromKey(ebiten.KeyEscape), pause)
	g.keymap[generalLayer].KeyMouse.Set(button.FromKey(ebiten.KeyF11), fullscreen)
	g.keymap[generalLayer].GamepadBtn.Set(ebiten.GamepadButton7, pause)
	g.keymap[generalLayer].GamepadBtn.Set(ebiten.GamepadButton6, fullscreen)

	playerActions := keymap.ButtonHandlerMap{
		left:   g.input.handleLeft,
		right:  g.input.handleRight,
		jump:   g.input.handleJump,
		punch:  g.input.handlePunch,
		launch: g.input.handleLaunch,
	}
	playerAxisActions := keymap.AxisHandlerMap{
		move:   g.input.handleMove,
		punchH: g.input.handlePunchH,
		punchV: g.input.handlePunchV,
	}
	g.keymap[playerLayer] = keymap.New(playerActions, playerAxisActions)
	setDefaultKeyMap(g.keymap[playerLayer])

	g.states = map[gamestate.State]gamestate.GameState{
		gamestate.Intro: newIntroState(g),
		// mainMenu: newMainMenu(p, screenHeight, screenWidth, cam, bg, g.keymap),
		// play:     newPlayState(p, screenHeight, cam, bg),
	}
	// descs, err := sprite.Psd(asset.Psd("test"))
	// if err != nil {
	// 	log.Fatalf("Adding PSD asset 'test' to collection: %v", err)
	// }
	// for i := range descs {
	// 	g.test[descs[i].Name] = &descs[i]
	// }
	// g.testSprites = []sprite.Sprite{
	// 	sprite.Sprite{
	// 		Desc: g.test["white"],
	// 	},
	// 	sprite.Sprite{
	// 		Desc: g.test["white"],
	// 		Loop: true,
	// 	},
	// 	sprite.Sprite{
	// 		Desc: g.test["green"],
	// 	},
	// }

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

	return g
}

// Update the Game by simulating the state by one frame.
func (g *Game) Update() {
	updateStart := time.Now()
	dt := g.dt(updateStart)

	g.keymap.Update()

	s := g.states[g.state]
	next := s.NextState()
	if next != g.state {
		log.Println("Change state from", g.state, "to", next)
		s.End()
		s = g.states[next]
		s.Begin(g.state)
		g.state = next
	}

	s.Update(dt)

	// g.counter += dt
	// if g.counter > 5*time.Second {
	// 	g.counter = 0
	// 	for i := range g.testSprites {
	// 		g.testSprites[i].Start()
	// 	}
	// }
	// for i := range g.testSprites {
	// 	g.testSprites[i].Update(dt)
	// }

	// g.camera.Update(dt)

	// g.handleCollisions()

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

	g.states[g.state].Draw(dst)

	// opts := ebiten.DrawImageOptions{}
	// opts.GeoM.Translate(geo.VecXY(50, 100).Minus(g.testSprites[0].Points("anchor")[0]).XY())
	// g.testSprites[0].Draw(dst, &opts)
	//
	// opts.GeoM.Reset()
	// opts.GeoM.Translate(geo.VecXY(50, 120).Minus(g.testSprites[1].Points("anchor")[0]).XY())
	// g.testSprites[1].Draw(dst, &opts)
	//
	// opts.GeoM.Reset()
	// opts.GeoM.Translate(geo.VecXY(50, 140).Minus(g.testSprites[2].Points("anchor")[0]).XY())
	// g.testSprites[2].Draw(dst, &opts)

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

// func (g *Game) handleCollisions() {
// 	// for _, box := range g.player.hitboxes() {
// 	// 	if !box.Active {
// 	// 		continue
// 	// 	}
// 	// 	if box.Bounds.Bottom() > 0 {
// 	// 		box.Callback(&g.groundHB)
// 	// 	}
// 	// }
// }

func (g *Game) render(dst *ebiten.Image) {
	entComp := comp.Position // | comp.Sprite
	cameraComp := comp.Position | comp.BoundingBox
	state := g.entityState
	for i, m := range state.Mask {
		e := entity(i)
		if !m.Contains(entComp) {
			continue
		}

		if m.Contains(comp.Camera) {
			if !state.Mask[state.Camera[e]].Contains(cameraComp) {
				continue
			}

			pos := state.Position[e]

			rotation := 0.0
			if m.Contains(comp.Rotation) {
				rotation = state.Rotation[e]
			}

			camera := state.Camera[e]
			cameraPos := state.Position[camera]
			cameraBounds := state.BoundingBox[camera].Moved(cameraPos.XY())

			if m.Contains(comp.BoundingBox) {
				entityBounds := state.BoundingBox[e].Moved(pos.XY())
				if cameraBounds.CollideRect(entityBounds) {
					screenPos := pos.Minus(geo.VecXY(cameraBounds.TopLeft()))
					_ = screenPos
					_ = rotation
				}
			} else {
				if cameraBounds.CollidePoint(pos.XY()) {
					screenPos := pos.Minus(geo.VecXY(cameraBounds.TopLeft()))
					_ = screenPos
					_ = rotation
				}
			}
		} else {
			// No camera, draw directly to screen
			pos := state.Position[e]

			rotation := 0.0
			if m.Contains(comp.Rotation) {
				rotation = state.Rotation[e]
			}

			_ = pos
			_ = rotation
		}
	}
}

func (g *Game) followUpdate(dt time.Duration) {
	entComp := comp.Position | comp.Follow
	targetComp := comp.Position
	state := g.entityState
	for i, m := range state.Mask {
		e := entity(i)
		if !m.Contains(entComp) || !state.Mask[state.Follow[e].Target].Contains(targetComp) {
			continue
		}
		pos := state.Position[e]

		params := state.Follow[e]
		targetPos := state.Position[params.Target]

		distToTarget2 := targetPos.Dist2(pos)
		max2 := params.MaxDist * params.MaxDist
		if distToTarget2 > max2 {
			state.Position[e] = targetPos.Plus(pos.Minus(targetPos).WithLen(params.MaxDist))
		} else {
			ratio := distToTarget2 / max2
			speed := params.Ease(ratio) * params.MaxSpeed
			vel := targetPos.Minus(pos).WithLen(speed)
			state.Position[e].Add(vel.Times(dt.Seconds()))
		}
	}
}

func (g *Game) shakeUpdate(dt time.Duration) {
	entComp := comp.Position | comp.Shake
	state := g.entityState
	for i, m := range state.Mask {
		e := entity(i)
		if !m.Contains(entComp) {
			continue
		}
		pos := &state.Position[e]
		params := &state.Shake[e]
		params.Time.Add(dt)
		if params.Shaker.Falloff != nil {
			pos.Add(params.Shaker.Shake(params.Time))
		} else {
			pos.Add(params.Shaker.ShakeConst(params.Time))
		}
	}
}
