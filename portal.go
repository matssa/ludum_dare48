package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const portalSize = 64

type Portal struct {
	animatedSprite *AnimatedSprite
	counter        int

	// Portal position
	x16 int
	y16 int
}

func findPlacement() (int, int) {
	for _, line := range tiles {
		if line.posx > worldWidth-1000 {
			return line.posx, line.posy - 64
		}
	}
	return 9500, 1150
}

func (g *Game) createPortal() {
	portalSprite := NewAnimatedSprite(
		0,
		0,
		portalSize,
		portalSize,
		4,
		portalImage)

	x, y := findPlacement()

	g.portal = Portal{
		animatedSprite: portalSprite,
		x16:            x,
		y16:            y,
	}
}

func (g *Game) drawPortal() {
	g.portal.counter++
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.portal.x16), float64(g.portal.y16))
	op.Filter = ebiten.FilterLinear
	g.world.DrawImage(g.portal.animatedSprite.GetCurrFrame(), op)
	if g.portal.counter >= 5 {
		g.portal.animatedSprite.NextFrame()
		g.portal.counter = 0
	}
}
