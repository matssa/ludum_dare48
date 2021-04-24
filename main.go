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
	"bytes"
	"image"
	_ "image/png"
	"log"

	"golang.org/x/image/math/f64"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	resources "github.com/hajimehoshi/ebiten/v2/examples/resources/images/flappy"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	runnerImage    *ebiten.Image
	animatedSprite *AnimatedSprite
	keys           = []ebiten.Key{
		ebiten.KeyA,
		ebiten.KeyW,
		ebiten.KeyS,
		ebiten.KeyD,
		ebiten.KeySpace,
	}
	myKeys      = []keyboardKey{}
	gopherImage *ebiten.Image
)

type keyboardKey struct {
	isPressed bool
	key       ebiten.Key
}

type Player struct {
	count     int
	hasTurned bool
	looksLeft bool

	// Character position
	x16  int
	y16  int
	vy16 int
	vx16 int
}

type Game struct {
	player   Player
	gameMode int

	layers [][]int
	world  *ebiten.Image
	camera Camera

	// Camera position
	cameraX int
	cameraY int
}

func init() {
	img, _, err := image.Decode(bytes.NewReader(resources.Gopher_png))
	if err != nil {
		log.Fatal(err)
	}
	gopherImage = ebiten.NewImageFromImage(img)

	initWorldImg()
}

func (g *Game) init() {
	g.player.x16 = 0
	g.player.y16 = 100 * 16
	g.cameraX = -240
	g.cameraY = 0
}

func newGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (p *Player) jump() {
	p.vy16 = -60
}

func (p *Player) executeMovement() {
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		p.jump()
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		if !p.looksLeft {
			p.looksLeft = true
			p.hasTurned = true
		}
		p.vx16 = -40
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		if p.looksLeft {
			p.looksLeft = false
			p.hasTurned = true
		}
		p.vx16 = 40
	} else {
		p.vx16 = 0
	}

	p.y16 += p.vy16
	p.x16 += p.vx16
	p.count++
}

func (g *Game) Update() error {

	controlCamera(g)

	switch g.gameMode {
	case play:
		// Gravity
		g.player.vy16 += 4
		if g.player.vy16 > 96 {
			g.player.vy16 = 96
		}

		g.player.executeMovement()

	default:
		g.gameMode = play
	}
	return nil
}

func (g *Game) drawCharacter(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	if g.player.looksLeft {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(animatedSprite.frameWidth), 0)
	}
	op.GeoM.Translate(float64(g.player.x16/16.0)-float64(g.cameraX), float64(g.player.y16/16.0)-float64(g.cameraY))
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	op.Filter = ebiten.FilterLinear
	if g.player.count >= 5 {
		g.player.count = 0
		animatedSprite.NextFrame()
	}
	screen.DrawImage(animatedSprite.GetCurrFrame(), op)
}

func (g *Game) Draw(screen *ebiten.Image) {
	renderWorld(g, screen)

	// sx, sy := frameOX+i*frameWidth, 0
	if g.gameMode == 1 {
		g.drawCharacter(screen)
		DrawOverlay(screen, 5)
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

	res, err := ebitenutil.OpenFile("./Squirrel-Sheet.png")
	if err != nil {
		log.Fatal(err)
	}
	// img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	img, _, err := image.Decode(res)
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebiten.NewImageFromImage(img)
	animatedSprite = NewAnimatedSprite(
		0,
		0,
		32,
		32,
		5,
		runnerImage)

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")
	newGame()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
