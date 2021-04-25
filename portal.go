package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const portalSize = 64

type Portal struct {
	animatedSprite *AnimatedSprite
	count          int

	// Portal position
	x16 int
	y16 int
}

func findPlacement(g *Game) {
	return 0, 0
}

func (g *Game) newPortal() *Portal {
	portalSprite = NewAnimatedSprite(
		0,
		0,
		portalSize,
		portalSize,
		2,
		portalImage)

	x, y := findPlacement(g)

	return &Portal{
		animatedSprite: portalSprite,
		x16:            x,
		y16:            y,
	}
}

func (g *Game) drawPortal() {
	op := &ebiten.DrawImageOptions{}
	if g.player.looksLeft {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(playerSprite.frameWidth), 0)
	}
	op.GeoM.Translate(float64(g.player.x16), float64(g.player.y16))
	op.Filter = ebiten.FilterLinear
	g.world.DrawImage(portalImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
}
