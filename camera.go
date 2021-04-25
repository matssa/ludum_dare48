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
	ViewPort     f64.Vec2
	Position     f64.Vec2
	Velocity     f64.Vec2
	Acceleration f64.Vec2
	ZoomFactor   int
	Rotation     int

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
	c.Velocity[0] = 0
	c.Velocity[1] = 0
	c.Acceleration[0] = 0
	c.Acceleration[1] = 0
	c.Rotation = 0
	c.ZoomFactor = 0
}

func (c *Camera) physicsTick(g *Game) {
	var damping float64 = 0.01

	c.Acceleration[0] -= c.Velocity[0] * damping
	c.Acceleration[1] -= c.Velocity[1] * damping

	c.Velocity[0] += c.Acceleration[0]
	c.Velocity[1] += c.Acceleration[1]

	c.Position[0] += c.Velocity[0]
	c.Position[1] += c.Velocity[1]
}

func (c *Camera) followCharacter(g *Game) {
	cx := g.player.x16
	cy := g.player.y16

	worldMatrix := c.worldMatrix()
	characterScreenPosX, characterScreenPosY := worldMatrix.Apply(float64(cx), float64(cy))

	var rulerLeft = 0.4 * screenWidth
	var rulerRight = 0.6 * screenWidth

	// Y is from top, i.e. larger Y -> lower screen pos
	var rulerTop = 0.4 * screenHeight
	var rulerBottom = 0.6 * screenHeight

	var overStepX float64 = 0
	var overStepY float64 = 0

	// Reset acceleration
	c.Acceleration[0] = 0
	c.Acceleration[0] = 0

	if characterScreenPosX > rulerRight {
		overStepX = characterScreenPosX - rulerRight
		c.Acceleration[0] = overStepX / 16
	}
	if characterScreenPosX < rulerLeft {
		overStepX = rulerLeft - characterScreenPosX
		c.Acceleration[0] = -overStepX / 16
	}
	if characterScreenPosY > rulerBottom {
		overStepY = characterScreenPosY - rulerBottom
		c.Acceleration[1] = overStepY / 16
	}
	if characterScreenPosY < rulerTop {
		overStepY = rulerTop - characterScreenPosY
		c.Acceleration[1] = -overStepY / 16
	}

	// var damping float64 = 1
	c.Velocity[0] = c.Velocity[0] - 0.5*c.Velocity[0]
	c.Velocity[1] = c.Velocity[1] - 0.5*c.Velocity[1]
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

// manuallyControl controls the camera based on key input
func (c *Camera) manuallyControl() {
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
		c.physicsTick(g)
		c.followCharacter(g)
	} else {
		c.manuallyControl()
	}

}
