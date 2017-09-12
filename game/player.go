package game

import (
	"math"
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
	playerPunchTime = 200 * time.Millisecond
	playerPunchGap  = 100 * time.Millisecond
)

type playerState int

const (
	awaken playerState = iota
	idle
	playerMove
	playerPunch
)

type player struct {
	pos geo.Vec
	vel geo.Vec

	left  bool    // Move left button is down
	right bool    // Move right button is down
	move  float64 // Gampad axis for movement
	jump  bool    // Jump button is down
	punch bool    // Punch button is down

	canJump   bool
	isJumping bool
	jumpTime  time.Duration
	flipDir   bool

	punchTime time.Duration
	punchGap  time.Duration

	state playerState

	currentSprite *sprite.Sprite
	awakenSprite  sprite.Sprite
	idleSprite    sprite.Sprite
	moveSprite    sprite.Sprite
	punchSprite   sprite.Sprite
}

func newPlayer() *player {
	p := &player{
		canJump:   true,
		isJumping: false,
		jumpTime:  0,

		state:        idle,
		awakenSprite: sprite.Get("awaken"),
		idleSprite:   sprite.Get("idle"),
		moveSprite:   sprite.Get("move"),
		punchSprite:  sprite.Get("punch"),
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

func (p *player) idle() {
	p.state = idle
	p.currentSprite = &p.idleSprite
}

func (p *player) update(dt time.Duration) {
	p.punchGap -= dt

	switch p.state {
	case awaken:
		if p.awoke() {
			p.state = idle
			p.currentSprite = &p.idleSprite
		}
	case idle:
		p.updateMove()
		if p.move != 0 {
			p.state = playerMove
			p.currentSprite = &p.moveSprite
		}
		if p.punch && p.punchGap <= 0 {
			p.state = playerPunch
			p.currentSprite = &p.punchSprite
			p.punchTime = playerPunchTime
		}
	case playerMove:
		p.updateMove()
		if p.move == 0 {
			p.idle()
		}
		if p.punch && p.punchGap <= 0 {
			p.punch = false
			p.state = playerPunch
			p.currentSprite = &p.punchSprite
			p.punchTime = playerPunchTime
		}
	case playerPunch:
		p.updateMove()
		p.punchTime -= dt
		if p.punchTime <= 0 {
			p.punchGap = playerPunchGap
			if p.move == 0 {
				p.state = playerMove
				p.currentSprite = &p.moveSprite
			} else {
				p.idle()
			}
		}
	}

	switch p.state {
	case awaken:
	case idle:
		p.updateMovement(dt)
	case playerMove:
		p.updateMovement(dt)
	case playerPunch:
		p.updateMovement(dt)
	}

	p.currentSprite.Update(dt)

	// Reset Move for next frame, it will be set each frame by user input.
	p.move = 0
}

func (p *player) draw(dst *ebiten.Image, cam *camera.Camera) {
	pos := cam.ScreenCoords(p.pos)
	geom := ebiten.GeoM{}

	size := p.currentSprite.Size()
	bounds := geo.RectWH(size.XY())
	bounds.SetBottomMid(pos.XY())

	switch p.state {
	case awaken:
	case idle:
	case playerMove:
		p.flipDir = p.vel.X < 0
	case playerPunch:
		mPos := geo.VecXYi(ebiten.CursorPosition())
		p.flipDir = mPos.X < pos.X
		pos = geo.VecXY(bounds.Mid())
		angle := mPos.Minus(pos).Angle()
		if p.flipDir {
			angle += math.Pi
			angle *= -1
		}
		geom.Translate(-size.X/2, -size.Y/2)
		geom.Rotate(-angle)
		geom.Translate(size.X/2, size.Y/2)
	}

	if p.flipDir {
		geom.Scale(-1, 1)
		geom.Translate(size.X, 0)
	}
	p.currentSprite.Draw(dst, pos, geom)
}

func (p *player) updateMove() {
	if p.move == 0 { // If gamepad axis isn't being used, check left/right buttons.
		if p.left {
			p.move = -1
		}
		if p.right {
			p.move = 1
		}
	}
}

func (p *player) updateMovement(dt time.Duration) {
	p.vel.X = p.move * playerMoveSpeed

	// Check if it's time to jump before handling jump the jump state so that we start
	// jumping as soon as possible
	if !p.isJumping && p.jump && p.canJump {
		p.isJumping = true
		p.jumpTime = playerJumpTime
	}

	if p.isJumping {
		if p.jumpTime <= 0 || !p.jump {
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
}

func (p *player) Pos() geo.Vec {
	return p.pos
}

func (p *player) SetPos(pos geo.Vec) {
	p.pos = pos
}

func (p *player) handleLeft(down bool) bool {
	p.left = down
	return false
}

func (p *player) handleRight(down bool) bool {
	p.right = down
	return false
}

func (p *player) handleMove(val float64) bool {
	p.move = val
	return false
}

func (p *player) handleJump(down bool) bool {
	p.jump = down
	return false
}

func (p *player) handlePunch(down bool) bool {
	p.punch = down
	return false
}
