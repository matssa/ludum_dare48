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

	// Character position
	x16  int
	y16  int
	vy16 int
	vx16 int
}

func (p *Player) jump() {
	p.vy16 = -60
}

func (p *Player) moveLeft() {
	if !p.looksLeft {
		p.looksLeft = true
		p.hasTurned = true
	}
	p.vx16 = -40
	p.restingCount = 0
	p.isResting = false
}

func (p *Player) moveRight() {
	if p.looksLeft {
		p.looksLeft = false
		p.hasTurned = true
	}
	p.vx16 = 40
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

func (p *Player) executeMovement() {
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		p.jump()
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.moveLeft()
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.moveRight()
	} else {
		p.rest()
	}
	p.y16 += p.vy16
	p.x16 += p.vx16
	p.count++
}
