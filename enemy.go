package main

import (
	"github.com/hajimehoshi/ebiten/v2"
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

	// Enemy position
	x16  int
	y16  int
	vy16 float64
	vx16 float64
}

func spawnPosition(g *Game) (int, int) {
	return 50, 20
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
