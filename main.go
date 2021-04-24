// Copyright 2018 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build example

package main

import (
	_ "image/png"
	"log"

	"golang.org/x/image/math/f64"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240

	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameNum    = 5

	play     = 1
	menu     = 2
	gameOver = 3
)

var (
	tiles = make([]*Tile, 0, 0)
)

type Game struct {
	player   Player
	gameMode int

	layers [][]int
	world  *ebiten.Image
	camera Camera

}

func init() {
	initAnimation()
	initWorldImg()
}

func (g *Game) init() {
	g.player.x16 = 0
	g.player.y16 = 100 * 16
}

func newGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) Update() error {

	controlCamera(g)

	switch g.gameMode {
	case play:
		// Gravity
		g.player.vy16 += 0.3
		if g.player.vy16 > 10 {
			g.player.vy16 = 10
		}

		g.player.executeMovement()

	default:
		g.gameMode = play
	}
	return nil
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

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawWorld()
	g.drawCharacter()

	g.camera.Render(g.world, screen)

	// sx, sy := frameOX+i*frameWidth, 0
	if g.gameMode == play {
		DrawOverlay(screen, 5)
		for _, tile := range tiles {
			tile.DrawTile(screen)
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		camera: Camera{ViewPort: f64.Vec2{screenWidth, screenHeight}},
		player: Player{count: 0, hasTurned: false},
	}
	buildWorld(g)


	// Create some tiles
	posx := 0
	for i := 0; i < 10; i++ {
		tiles = append(tiles, NewTile(posx, 100, "top"))
		posx += 16;
	}
	tiles = append(tiles, NewTile(posx, 100, "top-right"))
	posx = 64
	tiles = append(tiles, NewTile(posx, 150, "top-left"))
	for i := 0; i < 20; i++ {
		posx += 16;
		tiles = append(tiles, NewTile(posx, 150, "top"))
	}
	posx = 0
	for i := 0; i < 10; i++ {
		tiles = append(tiles, NewTile(posx, 200, "top"))
		posx += 16;
	}
	tiles = append(tiles, NewTile(posx, 200, "top-right"))



	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")

	newGame()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
