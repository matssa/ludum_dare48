package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Player struct {
	count        int
	restingCount int
	isResting    bool
	hasTurned    bool
	looksLeft    bool
	canJump      bool
	isDigging    bool

	// Character position
	x16  int
	y16  int
	vy16 float64
	vx16 float64
}

func (p *Player) jump() {
	p.vy16 = -7
}

func (p *Player) moveLeft() {
	if !p.looksLeft {
		p.looksLeft = true
	}
	if p.isDigging {
		p.vx16 = -1
	} else {
		p.vx16 = -3
	}
	p.restingCount = 0
	p.isResting = false
}

func (p *Player) moveRight() {
	if p.looksLeft {
		p.looksLeft = false
	}
	if p.isDigging {
		p.vx16 = 1
	} else {
		p.vx16 = 3
	}
	p.restingCount = 0
	p.isResting = false
}

func (p *Player) rest() {
	p.vx16 = 0
	p.restingCount++
	if p.restingCount >= 5 {
		p.isResting = true
	}
}

func (p *Player) dig() {
	p.vy16 = 1
	p.restingCount++
	if p.restingCount >= 5 {
		p.isResting = true
	}
}

func (p *Player) executeMovement() {
	// Gravity
	p.vy16 += gravity
	if p.vy16 > maxVelocityY {
		p.vy16 = maxVelocityY
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) && p.canJump {
		p.canJump = false
		p.isDigging = false
		p.jump()
	} else {
		for _, tile := range tiles {
			if tile.PlayerCollide(p) {
				p.canJump = true
				if ebiten.IsKeyPressed(ebiten.KeyS) {
					p.isDigging = true
					p.dig()
				} else {
					p.isDigging = false
					if p.vy16 >= 0 {
						p.vy16 = 0
					}
					p.y16 = tile.posy - 22 // TODO Need to offset the tile y pos ofcourse, but why does 22 work?
				}
			}
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.moveLeft()
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.moveRight()
	} else {
		p.rest()
	}
	p.y16 += int(p.vy16)
	p.x16 += int(p.vx16)
	p.count++
}

func (g *Game) drawCharacter() {
	op := &ebiten.DrawImageOptions{}
	if g.player.looksLeft {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(animatedSprite.frameWidth), 0)
	}
	op.GeoM.Translate(float64(g.player.x16), float64(g.player.y16))
	op.Filter = ebiten.FilterLinear
	if g.player.isResting {
		g.world.DrawImage(animatedIdleSprite.GetCurrFrame(), op)
		if g.player.restingCount >= 10 {
			g.player.restingCount = 0
			animatedIdleSprite.NextFrame()
		}
	} else {
		g.world.DrawImage(animatedSprite.GetCurrFrame(), op)
		if g.player.count >= 5 {
			g.player.count = 0
			animatedSprite.NextFrame()
		}
	}
}
