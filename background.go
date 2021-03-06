package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	bgImage *ebiten.Image
)

func initBackgroundImg() {
	bgImageFile, err := ebitenutil.OpenFile("./resources/sprites/background-640x480-repeated.png")
	if err != nil {
		log.Fatal(err)
	}
	bgImageRaw, _, err := image.Decode(bgImageFile)
	if err != nil {
		log.Fatal(err)
	}
	bgImage = ebiten.NewImageFromImage(bgImageRaw)
}

func (g *Game) drawBackground() {
	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(g.camera.Position[0]/2, g.camera.Position[1]/2)

		g.world.DrawImage(bgImage, op)
	}
	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate((g.camera.Position[0]/2)+(640*2), g.camera.Position[1]/2)

		g.world.DrawImage(bgImage, op)
	}
	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate((g.camera.Position[0]/2)+(640*4), g.camera.Position[1]/2)

		g.world.DrawImage(bgImage, op)
	}
}
