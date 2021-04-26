package main

import (
	"fmt"
	_ "image/png"
	"os"

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

	chipmunkSize = 32

	totNumEnemies = 50
)

var (
	tiles     = make([]*Tile, 0, 0)
	runClouds = false
)

type Game struct {
	player       Player
	enemies      []*Enemy
	enemyBullets []*EnemyBullet
	gameMode     int

	layers        [][]int
	world         *ebiten.Image
	camera        Camera
	ominousClouds OminousClouds
	portal        Portal
}

func init() {
	initAnimation()
	initBackgroundImg()
	initWorldImg()
}

func calcAliveEnemies(enemies []*Enemy) int {
	numAliveEnemies := 0
	for _, enemy := range enemies {
		if enemy.isAlive {
			numAliveEnemies += 1
		}
	}
	return numAliveEnemies
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

	if ebiten.IsKeyPressed(ebiten.KeyY) {
		g.player.isResting = false
		g.player.isAttacking = true
	}

	g.camera.update(g)

	switch g.gameMode {
	case play:

		g.player.executeMovement()
		g.UpdateEnemies()
		g.UpdateBullets()
		if g.isPlayerHit() {
			if g.player.health <= 0 {
				g.gameMode = gameOver
			}
		}

	case gameOver:
		fmt.Printf("Game Over! :(")
		g.player.health = 100
		g.gameMode = play

	default:
		g.gameMode = play
	}

	g.ominousClouds.UpdateClouds()


	aliveEnemies := calcAliveEnemies(g.enemies);
	if aliveEnemies < totNumEnemies {
		g.createEnemies(totNumEnemies - aliveEnemies);
	}

	if g.player.y16 > 1700 {
		fmt.Printf("\n\nYou lost...\n\n")
		os.Exit(0)
	}

	if g.player.health <= 0 {
		fmt.Printf("\n\nYou lost...\n\n")
		os.Exit(0)
	}

	if g.player.x16 >= g.portal.x16 {
		fmt.Printf("\n\nYou won!\n\n")
		os.Exit(0)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBackground()
	g.drawWorld()
	g.drawPortal()
	g.drawCharacter()
	g.drawEnemies()
	g.DrawBullets()

	// sx, sy := frameOX+i*frameWidth, 0
	if g.gameMode == play {
		for _, tile := range tiles {
			tile.DrawTile(g.world)
		}
	}

	g.ominousClouds.DrawClouds(g.world)

	// Anything relative to world must be drawn on g.world before calling
	// Render()
	g.camera.Render(g.world, screen)
	numAliveEnemies := calcAliveEnemies(g.enemies)
	DrawOverlay(screen, g.player.health, numAliveEnemies)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
