package main

import (
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)


type Tile struct {
	posx int
	posy int
	sprite_selector string

	open_image *ebiten.Image // PRIVATE FIELD
}

func NewTile(posx int, posy int, sprite_selector string) *Tile {
	tileimage, err := ebitenutil.OpenFile("./resources/sprites/tile-map.png");
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(tileimage)
	if err != nil {
		log.Fatal(err)
	}
	tileImage := ebiten.NewImageFromImage(img)

	return &Tile{
		posx:   posx,
		posy:   posy,
		sprite_selector: sprite_selector,
		open_image: tileImage,
	}
}

func (t Tile) DrawTile(screen *ebiten.Image) {
	var startx int
	var starty int
	switch t.sprite_selector {
	case "top-left":
		startx = 0
		starty = 0
	case "top":
		startx = 16
		starty = 0
	case "top-right":
		startx = 32
		starty = 0
	default:
		log.Fatal("sprite selector not implemented yet");
	}

	subimg := t.open_image.SubImage(image.Rect(
		startx,
		starty,
		startx + 16,
		starty + 16)).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(t.posx), float64(t.posy))
	screen.DrawImage(subimg, op)
}

func (tile Tile) PlayerCollide(p *Player) bool {
	playerLeftBorderInTile := p.x16 >= tile.posx && p.x16 <= tile.posx + 16
	playerRightBorderInTile := p.x16 + 16 >= tile.posx && p.x16 + 16 <= tile.posx + 16
	playerBottomInTopPortionOfTile := p.y16 + 16 >= tile.posy - 8 && p.y16 + 16 <= tile.posy - 4

	return ((playerLeftBorderInTile || playerRightBorderInTile) && playerBottomInTopPortionOfTile)
}
