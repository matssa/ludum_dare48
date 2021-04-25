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
	shooting  = 4
)

type Enemy struct {
	count        int
	restingCount int
	looksLeft    bool
	isResting    bool
	isShooting   bool

	// Enemy stuff
	isAlive                bool
	size                   float64
	ability                int
	behaviour              int
	animatedSprite         *AnimatedSprite
	animatedIdleSprite     *AnimatedSprite
	animatedShootingSprite *AnimatedSprite
	shootFrameCount        int
	action                 int
	changeActionAfter      time.Time

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
	animatedShootingSprite := NewAnimatedSprite(
		0,
		0,
		32,
		32,
		6,
		shootingEnemyImage)
	x, y := spawnPosition(g)
	return &Enemy{
		isAlive: true,
		size:    s,
		ability: a, behaviour: b,
		animatedSprite:         animatedSprite,
		animatedIdleSprite:     animatedIdleSprite,
		animatedShootingSprite: animatedShootingSprite,
		x16:                    x,
		y16:                    y,
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

func (e *Enemy) shoot() {
	e.vx16 = 0
	e.isShooting = true
	e.isResting = false
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

func (e *Enemy) shouldShoot(p Player) bool {
	aggr := 16

	isInY := e.y16+aggr >= p.y16 && e.y16-aggr <= p.y16+32
	var isOnCorrectSide bool
	if e.looksLeft {
		isOnCorrectSide = e.x16 > p.x16
	} else {
		isOnCorrectSide = e.x16 < p.x16
	}
	return isInY && isOnCorrectSide
}

func (g *Game) UpdateEnemies() {
	for i := range g.enemies {
		if (!g.enemies[i].isAlive) {
			continue;
		}
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
		if g.enemies[i].shouldShoot(g.player) {
			g.enemies[i].shoot()
		} else {
			if g.enemies[i].canChangeAction() {
				g.enemies[i].changeAction()
			}
		}
		g.enemies[i].y16 += int(g.enemies[i].vy16)
		g.enemies[i].x16 += int(g.enemies[i].vx16)
		g.enemies[i].count++
	}
}

func (g *Game) drawEnemies() {
	for i, e := range g.enemies {
		if (!g.enemies[i].isAlive) {
			continue;
		}
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
		} else if g.enemies[i].isShooting {
			e := g.enemies[i]
			g.world.DrawImage(e.animatedShootingSprite.GetCurrFrame(), op)
			if e.shootFrameCount >= 6 {
				g.CreateBullet(e.x16, e.y16+4, e.looksLeft)
				e.isShooting = false
				e.changeAction()
				e.shootFrameCount = 0
				e.animatedShootingSprite.ResetSprite()
			} else if e.count >= 5 {
				e.shootFrameCount += 1
				e.count = 0
				e.animatedShootingSprite.NextFrame()
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
