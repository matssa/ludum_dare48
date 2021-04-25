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
	screenWidth  = 640
	screenHeight = 480

	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameNum    = 5

	play     = 1
	menu     = 2
	gameOver = 3

	gravity      = 0.3
	maxVelocityY = 5
)

var (
	tiles = make([]*Tile, 0, 0)
	runClouds = false
)

type Game struct {
	player   Player
	enemies  []*Enemy
	gameMode int

	layers [][]int
	world  *ebiten.Image
	camera Camera
	ominousClouds OminousClouds
}

func init() {
	initAnimation()
	initBackgroundImg()
	initWorldImg()
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyO) {
		runClouds = true
		g.ominousClouds.StartClouds()
	}
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		runClouds = false
		g.ominousClouds.StopClouds()
	}

	g.camera.update(g)

	switch g.gameMode {
	case play:

		g.player.executeMovement()
		g.executeEnemyMovement()

	default:
		g.gameMode = play
	}


	g.ominousClouds.UpdateClouds()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBackground()
	g.drawWorld()
	g.drawCharacter()
	g.drawEnemies()

	// sx, sy := frameOX+i*frameWidth, 0
	if g.gameMode == play {
		DrawOverlay(g.world, 5)
		for _, tile := range tiles {
			tile.DrawTile(g.world)
		}
	}

	g.ominousClouds.DrawClouds(g.world)

	// Anything relative to world must be drawn on g.world before calling
	// Render()
	g.camera.Render(g.world, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		camera: Camera{ViewPort: f64.Vec2{screenWidth, screenHeight}},
		player: Player{count: 0, hasTurned: false},
	}
	g.createEnemies(3)
	buildWorld(g)

	// Create some tiles
	posx := 0
	for i := 0; i < 10; i++ {
		tiles = append(tiles, NewTile(posx, 100, "top"))
		posx += TILE_SIZE
	}
	tiles = append(tiles, NewTile(posx, 100, "top-right"))
	posx = 64
	tiles = append(tiles, NewTile(posx, 150, "top-left"))
	for i := 0; i < 20; i++ {
		posx += TILE_SIZE
		tiles = append(tiles, NewTile(posx, 150, "top"))
	}
	posx = 0
	for i := 0; i < 10; i++ {
		tiles = append(tiles, NewTile(posx, 200, "top"))
		posx += TILE_SIZE
	}
	tiles = append(tiles, NewTile(posx, 200, "top-right"))

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
