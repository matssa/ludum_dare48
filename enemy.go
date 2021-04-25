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
	size               float32
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

func newEnemy(s float32, a int, b int, g *Game) *Enemy {
	animatedSprite := NewAnimatedSprite(
		0,
		0,
		int(32*s),
		int(32*s),
		5,
		runnerEnemyImage)
	animatedIdleSprite := NewAnimatedSprite(
		0,
		0,
		int(32*s),
		int(32*s),
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
	moveDuration := rand.Float32() * 4
	e.changeActionAfter = time.Now().Add(time.Second * time.Duration(moveDuration))
	switch actionType {
	case moveLeft:
		e.moveLeft()
	case moveRight:
		e.moveRight()
	case idle:
		e.rest()
	default:
		e.jump()
	}
}

func (e *Enemy) canChangeAction() bool {
	return time.Now().After(e.changeActionAfter)
}

func (g *Game) executeEnemyMovement() {
	for i := range g.enemies {
		// Gravity
		g.enemies[i].vy16 += gravity
		if g.enemies[i].vy16 > maxVelocityY {
			g.enemies[i].vy16 = maxVelocityY
		}

		for _, tile := range tiles {
			if tile.EnemyCollide(g.enemies[i]) {
				if g.enemies[i].vy16 >= 0 {
					g.enemies[i].vy16 = 0
				}
				g.enemies[i].y16 = tile.posy - 22 // TODO Need to offset the tile y pos ofcourse, but why does 22 work?
			}
		}
		if g.enemies[i].canChangeAction() {
			g.enemies[i].changeAction()
		}
		g.enemies[i].y16 += int(g.enemies[i].vy16)
		g.enemies[i].x16 += int(g.enemies[i].vx16)
		g.enemies[i].count++
	}
}

func (g *Game) drawEnemies() {
	for i := range g.enemies {
		op := &ebiten.DrawImageOptions{}
		if g.enemies[i].looksLeft {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(float64(g.enemies[i].animatedSprite.frameWidth), 0)
		}
		op.GeoM.Translate(float64(g.enemies[i].x16), float64(g.enemies[i].y16))
		op.Filter = ebiten.FilterLinear
		if g.enemies[i].isResting {
			g.world.DrawImage(g.enemies[i].animatedIdleSprite.GetCurrFrame(), op)
			if g.enemies[i].restingCount >= 10 {
				g.enemies[i].restingCount = 0
				g.enemies[i].animatedIdleSprite.NextFrame()
			}
		} else {
			g.world.DrawImage(g.enemies[i].animatedSprite.GetCurrFrame(), op)
			if g.enemies[i].count >= 5 {
				g.enemies[i].count = 0
				g.enemies[i].animatedSprite.NextFrame()
			}
		}
	}
}
