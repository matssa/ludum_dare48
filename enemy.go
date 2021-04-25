package main

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	jump      = 0
	moveLeft  = 1
	moveRight = 2
	idle      = 3
)

type Enemy struct {
	count        int
	restingCount int
	looksLeft    bool
	isResting    bool

	// Enemy stuff
	isAlive            bool
	size               float64
	ability            int
	behaviour          int
	animatedSprite     *AnimatedSprite
	animatedIdleSprite *AnimatedSprite
	action             int
	changeActionAfter  time.Time

	// Enemy position
	x16  int
	y16  int
	vy16 float64
	vx16 float64
}

func spawnPosition(g *Game) (int, int) {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(tileSize * tileXNum), 20
}

func newEnemy(s float64, a int, b int, g *Game) *Enemy {
	animatedSprite := NewAnimatedSprite(
		0,
		0,
		chipmunkSize,
		chipmunkSize,
		5,
		runnerEnemyImage)
	animatedIdleSprite := NewAnimatedSprite(
		0,
		0,
		chipmunkSize,
		chipmunkSize,
		3,
		idleEnemyImage)
	x, y := spawnPosition(g)
	return &Enemy{
		isAlive: true,
		size:    s,
		ability: a, behaviour: b,
		animatedSprite:     animatedSprite,
		animatedIdleSprite: animatedIdleSprite,
		x16:                x,
		y16:                y,
	}
}

func (g *Game) createEnemies(num int) {
	for i := 0; i < num; i++ {
		g.enemies = append(g.enemies, newEnemy(1, 1, 1, g))
	}
}

func (e *Enemy) jump() {
	e.vy16 = -7
}

func (e *Enemy) moveLeft() {
	if !e.looksLeft {
		e.looksLeft = true
	}
	e.vx16 = -1
	e.restingCount = 0
	e.isResting = false
}

func (e *Enemy) moveRight() {
	if e.looksLeft {
		e.looksLeft = false
	}
	e.vx16 = 1
	e.restingCount = 0
	e.isResting = false
}

func (e *Enemy) rest() {
	e.vx16 = 0
	e.restingCount++
	e.isResting = true
}

func (e *Enemy) changeAction() {
	rand.Seed(time.Now().UnixNano())
	actionType := rand.Intn(4)
	moveDuration := float32(0)
	switch actionType {
	case moveLeft:
		moveDuration = rand.Float32() * 2
		e.moveLeft()
	case moveRight:
		moveDuration = rand.Float32() * 2
		e.moveRight()
	case idle:
		moveDuration = rand.Float32() * 3
		e.rest()
	default:
		moveDuration = rand.Float32() * 2
		e.jump()
	}
	e.changeActionAfter = time.Now().Add(time.Second * time.Duration(moveDuration))
}

func (e *Enemy) canChangeAction() bool {
	return time.Now().After(e.changeActionAfter)
}

func (g *Game) executeEnemyMovement() {
	for _, e := range g.enemies {
		// Gravity
		e.vy16 += gravity
		if e.vy16 > maxVelocityY {
			e.vy16 = maxVelocityY
		}

		for _, tile := range tiles {
			if tile.EnemyCollide(e) {
				if e.vy16 >= 0 {
					e.vy16 = 0
				}
				e.y16 = tile.posy - 22 // TODO Need to offset the tile y pos ofcourse, but why does 22 work?
			}
		}
		if e.canChangeAction() {
			e.changeAction()
		}
		e.y16 += int(e.vy16)
		e.x16 += int(e.vx16)
		e.count++
	}
}

func (g *Game) drawEnemies() {
	for _, e := range g.enemies {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(e.size, e.size)
		if e.looksLeft {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(float64(e.animatedSprite.frameWidth), 0)
		}
		op.GeoM.Translate(float64(e.x16), float64(e.y16))
		op.Filter = ebiten.FilterLinear
		if e.isResting {
			g.world.DrawImage(e.animatedIdleSprite.GetCurrFrame(), op)
			if e.restingCount >= 10 {
				e.restingCount = 0
				e.animatedIdleSprite.NextFrame()
			}
		} else {
			g.world.DrawImage(e.animatedSprite.GetCurrFrame(), op)
			if e.count >= 5 {
				e.count = 0
				e.animatedSprite.NextFrame()
			}
		}
	}
}
