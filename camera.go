package main

import (
	"fmt"
	"math"

	"golang.org/x/image/math/f64"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Camera struct {
	ViewPort   f64.Vec2
	Position   f64.Vec2
	ZoomFactor int
	Rotation   int

	drawDebug    bool
	manualCamera bool
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %d, S: %d",
		c.Position, c.Rotation, c.ZoomFactor,
	)
}

func (c *Camera) viewportCenter() f64.Vec2 {
	return f64.Vec2{
		c.ViewPort[0] * 0.5,
		c.ViewPort[1] * 0.5,
	}
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Position[0], -c.Position[1])
	// We want to scale and rotate around center of image / screen
	m.Translate(-c.viewportCenter()[0], -c.viewportCenter()[1])
	m.Scale(
		math.Pow(1.01, float64(c.ZoomFactor)),
		math.Pow(1.01, float64(c.ZoomFactor)),
	)
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])
	return m
}

func (c *Camera) Render(world, screen *ebiten.Image) {
	screen.DrawImage(world, &ebiten.DrawImageOptions{
		GeoM: c.worldMatrix(),
	})

	if c.drawDebug {
		c.renderCameraDebug(screen)
	}
}

func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.worldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		// When scaling it can happend that matrix is not invertable
		return math.NaN(), math.NaN()
	}
}

func (c *Camera) Reset() {
	c.Position[0] = 0
	c.Position[1] = 0
	c.Rotation = 0
	c.ZoomFactor = 0
}

func (c *Camera) followCharacter(g *Game) {
	cx := g.player.x16
	cy := g.player.y16

	worldMatrix := c.worldMatrix()
	characterScreenPosX, characterScreenPosY := worldMatrix.Apply(float64(cx), float64(cy))

	if characterScreenPosX > 0.7*screenWidth {
		c.Position[0] += 4
	}
	if characterScreenPosX < 0.3*screenWidth {
		c.Position[0] -= 4
	}
	if characterScreenPosY > 0.5*screenHeight {
		c.Position[1] += 4
	}
	if characterScreenPosY < 0.3*screenHeight {
		c.Position[1] -= 4
	}
}

func (c *Camera) toggleDebug() {
	c.drawDebug = !c.drawDebug
}

func (c *Camera) toggleCameraControls() {
	c.manualCamera = !c.manualCamera
}

func (c *Camera) renderCameraDebug(screen *ebiten.Image) {
	worldX, worldY := c.ScreenToWorld(ebiten.CursorPosition())
	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("%s\nCursor World Pos: %.2f,%.2f\ncameraControls (V to toggle): %t",
			c.String(),
			worldX, worldY,
			c.manualCamera),
		0, screenHeight-48,
	)
}

func (c *Camera) manuallyControl() {
	// Manual controls in debug mode
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		c.Position[0] -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		c.Position[0] += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		c.Position[1] -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		c.Position[1] += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		if c.ZoomFactor > -2400 {
			c.ZoomFactor -= 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		if c.ZoomFactor < 2400 {
			c.ZoomFactor += 1
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		c.Rotation += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		c.Reset()
	}
}

func (c *Camera) update(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		c.toggleDebug()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		c.toggleCameraControls()
	}

	if !c.manualCamera {
		c.followCharacter(g)
	} else {
		c.manuallyControl()
	}

}
