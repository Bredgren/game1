package game

import (
	"fmt"
	"time"

	"github.com/Bredgren/game1/game/comp"
	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

type entity int

func (e entity) String() string {
	return fmt.Sprintf("Entity(%d)", int(e))
}

type entityPool struct {
	taken map[entity]bool
	free  []entity
}

func newEntityPool(maxEntities int) *entityPool {
	free := make([]entity, maxEntities)
	for i := 0; i < maxEntities; i++ {
		free[i] = entity(i)
	}
	return &entityPool{
		taken: make(map[entity]bool, maxEntities),
		free:  free,
	}
}

func (ep *entityPool) newEntity() (entity, error) {
	if len(ep.free) == 0 {
		return -1, fmt.Errorf("Max entities reached")
	}
	old := ep.free
	e, new := old[len(old)-1], old[:len(old)-1]
	ep.free = new
	ep.taken[e] = true
	return e, nil
}

func (ep *entityPool) delEntity(e entity) {
	if !ep.taken[e] {
		return
	}
	ep.taken[e] = false
	ep.free = append(ep.free, e)
}

type state struct {
	*entityPool
	Mask        []comp.Mask
	Position    []geo.Vec
	Velocity    []geo.Vec
	Rotation    []float64
	Gravity     []geo.Vec
	BoundingBox []geo.Rect
	Camera      []entity
	Follow      []followParams
	Shake       []shakeParams
	Sprite      []*ebiten.Image
	Hitbox      [][]geo.Rect
	Hurtbox     [][]geo.Rect
	Animation   []animationParams
}

func newState(maxEntities int) *state {
	return &state{
		entityPool:  newEntityPool(maxEntities),
		Mask:        make([]comp.Mask, maxEntities),
		Position:    make([]geo.Vec, maxEntities),
		Velocity:    make([]geo.Vec, maxEntities),
		Rotation:    make([]float64, maxEntities),
		Gravity:     make([]geo.Vec, maxEntities),
		BoundingBox: make([]geo.Rect, maxEntities),
		Camera:      make([]entity, maxEntities),
		Follow:      make([]followParams, maxEntities),
		Shake:       make([]shakeParams, maxEntities),
		Sprite:      make([]*ebiten.Image, maxEntities),
		Hitbox:      make([][]geo.Rect, maxEntities),
		Hurtbox:     make([][]geo.Rect, maxEntities),
		Animation:   make([]animationParams, maxEntities),
	}
}

func (s *state) delEntity(e entity) {
	s.Mask[e] = comp.None
	s.entityPool.delEntity(e)
}

type followParams struct {
	Target   entity
	MaxDist  float64
	MaxSpeed float64
	Ease     geo.EaseFn
}

type shakeParams struct {
	Shaker geo.Shaker
	Time   time.Time
}

type animationParams struct {
	currentAnimation string
	currentFrame     int
	timeLeft         time.Duration
	Animations       map[string][]frameDesc
}

func (a *animationParams) Play(anim string) {
	a.currentFrame = 0
	a.currentAnimation = anim
	a.timeLeft = a.Animations[anim][0].Duration
}

func (a *animationParams) AdvanceFrame() {
	maxFrame := len(a.Animations[a.currentAnimation]) - 1
	if a.currentFrame >= maxFrame {
		return
	}
	a.currentFrame++
	a.timeLeft = a.Animations[a.currentAnimation][a.currentFrame].Duration
}

func (a *animationParams) Ended() bool {
	return a.currentFrame == len(a.Animations[a.currentAnimation])-1 && a.timeLeft <= 0
}

func (a *animationParams) CurrentImg() *ebiten.Image {
	return a.Animations[a.currentAnimation][a.currentFrame].Img
}

func (a *animationParams) CurrentHitbox() []geo.Rect {
	return a.Animations[a.currentAnimation][a.currentFrame].Hitbox
}

func (a *animationParams) CurrentHurtbox() []geo.Rect {
	return a.Animations[a.currentAnimation][a.currentFrame].Hurtbox
}

func (a *animationParams) CurrentAnchor() geo.Vec {
	return a.Animations[a.currentAnimation][a.currentFrame].Anchor
}

type frameDesc struct {
	Img      *ebiten.Image
	Duration time.Duration
	Hitbox   []geo.Rect
	Hurtbox  []geo.Rect
	Anchor   geo.Vec
}
