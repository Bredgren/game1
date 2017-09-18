package game

import (
	"image/color"
	"log"
	"math"
	"time"

	"github.com/Bredgren/game1/game/camera"
	"github.com/Bredgren/game1/game/sprite"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
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
	cam *camera.Camera
	pos geo.Vec
	vel geo.Vec

	left             bool    // Move left button is down
	right            bool    // Move right button is down
	move             float64 // Gampad axis for movement
	jump             bool    // Jump button is down
	punch            bool    // Punch button is down
	punchAxis        geo.Vec
	punchWithGamepad bool

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

	coreHitbox   hitbox
	attackHitbox hitbox
}

func newPlayer(cam *camera.Camera) *player {
	p := &player{
		cam:       cam,
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

	p.coreHitbox = hitbox{
		Label:    "PlayerCore",
		Bounds:   geo.RectWH(2, 2),
		Callback: p.coreHit,
		Active:   true,
		Owner:    &p,
	}

	p.attackHitbox = hitbox{
		Label:    "PlayerAttack",
		Callback: p.attackHit,
		Owner:    &p,
	}

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

func (p *player) doIdle() {
	p.state = idle
	p.currentSprite = &p.idleSprite
}

func (p *player) doMove() {
	p.state = playerMove
	p.currentSprite = &p.moveSprite
}

func (p *player) shouldStartPunch() bool {
	p.punchWithGamepad = p.punchAxis.Len2() > 0.25
	return p.punchGap <= 0 && (p.punch || p.punchWithGamepad)
}

func (p *player) doPunch() {
	p.punch = false
	p.state = playerPunch
	p.currentSprite = &p.punchSprite
	p.punchTime = playerPunchTime
}

func (p *player) update(dt time.Duration) {
	p.punchGap -= dt

	switch p.state {
	case awaken:
		if p.awoke() {
			p.doIdle()
		}
	case idle:
		p.updateMove()
		if p.move != 0 {
			p.doMove()
		}
		if p.shouldStartPunch() {
			p.doPunch()
		}
	case playerMove:
		p.updateMove()
		if p.move == 0 {
			p.doIdle()
		}
		if p.shouldStartPunch() {
			p.doPunch()
		}
	case playerPunch:
		p.updateMove()
		p.punchTime -= dt
		if p.punchTime <= 0 {
			p.punchGap = playerPunchGap
			if p.move == 0 {
				p.doMove()
			} else {
				p.doIdle()
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

	p.updateHitboxes()

	p.currentSprite.Update(dt)

	// Reset Move for next frame, it will be set each frame by user input.
	p.move = 0
}

func (p *player) updateHitboxes() {
	// Since the player is being updated before the camera this will be off by a frame.
	// I don't think it will be noticeable though.
	mousePos := p.cam.WorldCoords(geo.VecXYi(ebiten.CursorPosition()))

	centerY := p.currentSprite.Size().Y / 2
	switch p.state {
	case awaken, idle, playerMove:
		p.attackHitbox.Active = false
	case playerPunch:
		p.attackHitbox.Active = true
		p.attackHitbox.Bounds.SetSize(8, 8)
		center := p.pos.Plus(geo.VecXY(0, -centerY))
		toMouse := mousePos.Minus(center)
		if p.punchWithGamepad {
			toMouse = p.punchAxis
		}
		toMouse.SetLen(9)
		p.attackHitbox.Bounds.SetMid(center.Plus(toMouse).XY())
	}

	p.coreHitbox.Bounds.SetTopLeft(p.pos.Plus(geo.VecXY(0, -centerY)).XY())
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
		pos = geo.VecXY(bounds.Mid())

		if p.punchWithGamepad {
			mPos = pos.Plus(p.punchAxis)
		}

		p.flipDir = mPos.X < pos.X
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

	// debug draw hitboxes
	if p.attackHitbox.Active {
		x, y := cam.ScreenCoords(geo.VecXY(p.attackHitbox.Bounds.TopLeft())).XY()
		w, h := p.attackHitbox.Bounds.Size()
		ebitenutil.DrawRect(dst, x, y, w, h, color.RGBA{0xFF, 0x00, 0x00, 0x88})
	}
	// if p.coreHitbox.Active {
	// 	x, y := cam.ScreenCoords(geo.VecXY(p.coreHitbox.Bounds.TopLeft())).XY()
	// 	w, h := p.coreHitbox.Bounds.Size()
	// 	ebitenutil.DrawRect(dst, x, y, w, h, color.RGBA{0x00, 0xFF, 0x00, 0x88})
	// }
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

func (p *player) handlePunchH(val float64) bool {
	p.punchAxis.X = val
	return false
}

func (p *player) handlePunchV(val float64) bool {
	p.punchAxis.Y = -val
	return false
}

func (p *player) hitboxes() []*hitbox {
	return []*hitbox{&p.coreHitbox, &p.attackHitbox}
}

func (p *player) coreHit(other *hitbox) {
	log.Println("coreHit:", other.Label)
}

func (p *player) attackHit(other *hitbox) {
	log.Println("attackHit:", other.Label)
	switch other.Label {
	case "ground":
		p.vel.Y = -125
		p.cam.Shaker.Amplitude = 10
		p.cam.Shaker.Duration = 500 * time.Millisecond
		p.cam.Shaker.Frequency = 5
		p.cam.Shaker.Falloff = geo.EaseOutQuad
		p.cam.StartShake()
	}
}
