package main

import (
	"fmt"
	"time"

	// "log"

	"github.com/hajimehoshi/ebiten/v2"
)

type EnemyBullet struct {
	posx           int
	posy           int
	goingLeft      bool
	count          int
	animatedSprite *AnimatedSprite
	isDestroyed    bool
	expireAt       time.Time
}

const (
	BULLET_VELOCITY = 7
)

var (
	bulletAnimatedSprite *AnimatedSprite
)

func init() {
	bulletAnimatedSprite = NewAnimatedSprite(
		0,
		0,
		16,
		16,
		13,
		enemyBulletImage)
}

func (g *Game) CreateBullet(posx int, posy int, goingLeft bool) {
	bulletAnimatedSprite = NewAnimatedSprite(
		0,
		0,
		16,
		16,
		13,
		enemyBulletImage)
	myNewBullet := &EnemyBullet{
		posx:           posx,
		posy:           posy,
		goingLeft:      goingLeft,
		animatedSprite: bulletAnimatedSprite,
		expireAt:       time.Now().Add(time.Second * 3),
	}
	g.enemyBullets = append(g.enemyBullets, myNewBullet)
}

func (g *Game) UpdateBullets() {
	for _, bullet := range g.enemyBullets {
		if time.Now().After(bullet.expireAt) {
			bullet.isDestroyed = true
		}
		if bullet.isDestroyed {
			continue
		}
		// Update the bullets
		if bullet.goingLeft {
			bullet.posx -= BULLET_VELOCITY
		} else {
			bullet.posx += BULLET_VELOCITY
		}

		if bullet.IsPlayerHit(g.player) {
			fmt.Printf("Its a hit!")
			g.player.health -= 1
			bullet.isDestroyed = true
		}
	}
}

func (g *Game) DrawBullets() {
	for _, bullet := range g.enemyBullets {
		if bullet.isDestroyed {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(bullet.posx), float64(bullet.posy))
		g.world.DrawImage(bullet.animatedSprite.GetCurrFrame(), op)
		if bullet.count >= 5 {
			bullet.count = 0
			bullet.animatedSprite.NextFrame()
		}
		bullet.count += 1
	}
}

func (b *EnemyBullet) IsPlayerHit(p Player) bool {
	aggr := 6

	isInX := b.posx >= p.x16 && b.posx+aggr <= p.x16+32
	isInY := b.posy >= p.y16 && b.posy+aggr <= p.y16+32

	return isInY && isInX
}
