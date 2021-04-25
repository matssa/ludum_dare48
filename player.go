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
		p.hasTurned = true
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
		p.hasTurned = true
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
	if inpututil.IsKeyJustPressed(ebiten.KeyW) && p.canJump {
		p.canJump = false
		p.isDigging = false
		p.jump()
	} else {
		for _, tile := range tiles {
			if tile.PlayerCollide(p) {
				p.canJump = true
				p.vy16 = 0
				if ebiten.IsKeyPressed(ebiten.KeyS) {
					p.isDigging = true
					p.dig()
				} else {
					p.isDigging = false
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
