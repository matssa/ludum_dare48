package main


import (
	_ "image/png"
	"log"

	"golang.org/x/image/math/f64"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := &Game{
		camera: Camera{ViewPort: f64.Vec2{screenWidth, screenHeight}},
		player: Player{health: 10, count: 0, hasTurned: false, x16: 500, y16: 500},
	}
	g.createEnemies(totNumEnemies) // totNumEnemies defined in game.go
	buildWorld(g)

	tileLines := createMap()
	for _, line := range tileLines {
		for _, tile := range line {
			tiles = append(tiles, tile)
		}
	}
	g.createPortal()

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("CyberSchmunk 2077")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
