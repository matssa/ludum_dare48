package main

import (
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const TILE_SIZE = 16

var (
	tileImage *ebiten.Image
)

type Tile struct {
	posx            int
	posy            int
	sprite_selector string // String for selecting what tile to draw. Example: "top" or "top-right" etc.
}

func init() {
	tmp, err := ebitenutil.OpenFile("./resources/sprites/tile-map.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(tmp)
	if err != nil {
		log.Fatal(err)
	}
	tileImage = ebiten.NewImageFromImage(img)
}

func NewTile(posx int, posy int, sprite_selector string) *Tile {
	return &Tile{
		posx:            posx,
		posy:            posy,
		sprite_selector: sprite_selector,
	}
}

func (t Tile) DrawTile(screen *ebiten.Image) {
	var startx int
	var starty int

	// Choose subimage based on sprite_selector
	switch t.sprite_selector {
	case "top-left":
		startx = 0
		starty = 0
	case "top":
		startx = TILE_SIZE
		starty = 0
	case "top-right":
		startx = TILE_SIZE * 2
		starty = 0
	default:
		log.Fatal("sprite selector not implemented yet")
	}

	subimg := tileImage.SubImage(image.Rect(
		startx,
		starty,
		startx+TILE_SIZE,
		starty+TILE_SIZE)).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(t.posx), float64(t.posy))
	screen.DrawImage(subimg, op)
}

// Check if tile is colliding with player or vice-versa
func (tile Tile) PlayerCollide(p *Player) bool {
	playerLeftBorderInTile := p.x16 >= tile.posx && p.x16 <= tile.posx+TILE_SIZE
	playerRightBorderInTile := p.x16+16 >= tile.posx && p.x16+16 <= tile.posx+TILE_SIZE
	// Not sure why the numbers 8 and 4 work here.. the idea is to get a piece of the top portion of the tile. Those numbers seem to do the job.
	playerBottomInTopPortionOfTile := p.y16+16 >= tile.posy-8 && p.y16+16 <= tile.posy-4

	return ((playerLeftBorderInTile || playerRightBorderInTile) && playerBottomInTopPortionOfTile)
}

func (tile Tile) EnemyCollide(e *Enemy) bool {
	enemyLeftBorderInTile := e.x16 >= tile.posx && e.x16 <= tile.posx+TILE_SIZE
	enemyRightBorderInTile := e.x16+16 >= tile.posx && e.x16+16 <= tile.posx+TILE_SIZE
	// Not sure why the numbers 8 and 4 work here.. the idea is to get a piece of the top portion of the tile. Those numbers seem to do the job.
	enemyBottomInTopPortionOfTile := e.y16+16 >= tile.posy-8 && e.y16+16 <= tile.posy-4

	return ((enemyLeftBorderInTile || enemyRightBorderInTile) && enemyBottomInTopPortionOfTile)
}
