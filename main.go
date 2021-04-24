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

type Game struct {
	count    int
	gameMode int

	// Character position
	x16  int
	y16  int
	vy16 int
	vx16 int

	// Camera position
	cameraX int
	cameraY int
}

func init() {
	w := keyboardKey{false, ebiten.KeyW}
	a := keyboardKey{false, ebiten.KeyA}
	s := keyboardKey{false, ebiten.KeyS}
	d := keyboardKey{false, ebiten.KeyD}
	space := keyboardKey{false, ebiten.KeySpace}
	myKeys = []keyboardKey{w, a, s, d, space}

	img, _, err := image.Decode(bytes.NewReader(resources.Gopher_png))
	if err != nil {
		log.Fatal(err)
	}
	gopherImage = ebiten.NewImageFromImage(img)
}

func (g *Game) init() {
	g.x16 = 0
	g.y16 = 100 * 16
	g.cameraX = -240
	g.cameraY = 0
}

func newGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func jump(g *Game) {
	g.vy16 = -96
}

func execueMovement(g *Game) {
	for i, key := range myKeys {
		if inpututil.IsKeyJustPressed(key.key) {
			myKeys[i].isPressed = true
			if key.key == ebiten.KeyW {
				jump(g)
			}
		}
		if inpututil.IsKeyJustReleased(key.key) {
			myKeys[i].isPressed = false
		}
	}

	if myKeys[1].isPressed {
		g.vx16 = -40
	} else if myKeys[3].isPressed {
		g.vx16 = 40
	} else {
		g.vx16 = 0
	}
	g.y16 += g.vy16
	g.x16 += g.vx16
	g.count++
}

func (g *Game) Update() error {
	switch g.gameMode {
	case play:
		// Gravity
		g.vy16 += 4
		if g.vy16 > 96 {
			g.vy16 = 96
		}

		execueMovement(g)
		// fmt.Printf("key pressed w a s d space: %v, %v, %v, %v, %v", myKeys[0].isPressed, myKeys[1].isPressed, myKeys[2].isPressed, myKeys[3].isPressed, myKeys[4].isPressed)

	default:
		g.gameMode = play
	}
	return nil
}

func (g *Game) drawCharacter(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.x16/16.0)-float64(g.cameraX), float64(g.y16/16.0)-float64(g.cameraY))
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	op.Filter = ebiten.FilterLinear
	if g.count == 5 {
		g.count = 0
		animatedSprite.NextFrame()
	}
	screen.DrawImage(animatedSprite.GetCurrFrame(), op)
}

func (g *Game) Draw(screen *ebiten.Image) {
	// sx, sy := frameOX+i*frameWidth, 0
	if g.gameMode == 1 {
		g.drawCharacter(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
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

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
