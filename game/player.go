package game

import (
	"image/color"
	"time"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/game1/game/sprite"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

const (
	playerMoveSpeed = 500
	playerGravity   = 50
	playerJumpSpeed = 700
	playerJumpTime  = 500 * time.Millisecond
)

type playerState int

const (
	awaken playerState = iota
	idle
)

type player struct {
	pos    geo.Vec
	vel    geo.Vec
	bounds geo.Rect
	img    *ebiten.Image

	Left  bool    // Move left button is down
	Right bool    // Move right button is down
	Move  float64 // Gampad axis for movement
	Jump  bool    // Jump button is down

	canJump   bool
	isJumping bool
	jumpTime  time.Duration

	state playerState

	currentSprite *sprite.Sprite
	awakenSprite  sprite.Sprite
	idleSprite    sprite.Sprite
}

func newPlayer() *player {
	img, _ := ebiten.NewImage(16, 16, ebiten.FilterNearest)
	img.Fill(color.Black)
	p := &player{
		bounds: geo.RectWH(16, 16),
		img:    img,

		canJump:   true,
		isJumping: false,
		jumpTime:  0,

		state:        idle,
		awakenSprite: sprite.Get("awaken"),
		idleSprite:   sprite.Get("idle"),
	}

	p.currentSprite = &p.idleSprite

	return p
}

func (p *player) awaken() {
	p.currentSprite = &p.awakenSprite
	p.currentSprite.Begin(false)
	p.state = awaken
}

func (p *player) awoke() bool {
	return p.awakenSprite.Ended()
}

func (p *player) update(dt time.Duration) {
	switch p.state {
	case awaken:
		if p.awoke() {
			p.state = idle
			p.currentSprite = &p.idleSprite
		}
	case idle:
		p.updateMovement(dt)
	}

	p.currentSprite.Update(dt)
}

func (p *player) draw(dst *ebiten.Image, cam *camera.Camera) {
	pos := cam.ScreenCoords(geo.VecXY(p.bounds.BottomMid()))
	p.currentSprite.Draw(dst, pos)
}

func (p *player) updateMovement(dt time.Duration) {
	if p.Move == 0 { // If gamepad axis isn't being used, check left/right buttons.
		if p.Left {
			p.Move = -1
		}
		if p.Right {
			p.Move = 1
		}
	}

	p.vel.X = p.Move * playerMoveSpeed

	// Check if it's time to jump before handling jump the jump state so that we start
	// jumping as soon as possible
	if !p.isJumping && p.Jump && p.canJump {
		p.isJumping = true
		p.jumpTime = playerJumpTime
	}

	if p.isJumping {
		if p.jumpTime <= 0 || !p.Jump {
			p.isJumping = false
		} else {
			p.jumpTime -= dt
			p.vel.Y = geo.Lerp(0, -playerJumpSpeed, p.jumpTime.Seconds()/playerJumpTime.Seconds())
		}
	} else {
		p.vel.Y += playerGravity
	}

	p.pos.Add(p.vel.Times(dt.Seconds()))

	// Ground collision
	if p.pos.Y > 0 {
		p.pos.Y = 0
		p.canJump = true
	} else {
		p.canJump = false
	}

	p.bounds.SetBottomMid(p.pos.XY())

	// Reset Move for next frame, it will be set each frame by user input.
	p.Move = 0
}

func (p *player) Pos() geo.Vec {
	return p.pos
}

func (p *player) SetPos(pos geo.Vec) {
	p.pos = pos
	p.bounds.SetBottomMid(p.pos.XY())
}

func (p *player) handleLeft(down bool) bool {
	p.Left = down
	return false
}

func (p *player) handleRight(down bool) bool {
	p.Right = down
	return false
}

func (p *player) handleMove(val float64) bool {
	p.Move = val
	return false
}

func (p *player) handleJump(down bool) bool {
	p.Jump = down
	return false
}
