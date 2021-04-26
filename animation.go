package main

import (
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	runnerImage        *ebiten.Image
	idleImage          *ebiten.Image
	attackImage        *ebiten.Image
	playerSprite       *AnimatedSprite
	playerIdleSprite   *AnimatedSprite
	playerAttackSprite *AnimatedSprite
	runnerEnemyImage   *ebiten.Image
	idleEnemyImage     *ebiten.Image
	shootingEnemyImage *ebiten.Image
	enemyBulletImage   *ebiten.Image
	portalImage        *ebiten.Image
)

type AnimatedSprite struct {
	frameOX     int // should probably be 0
	frameOY     int // should probably be 0
	frameWidth  int
	frameHeight int
	frameNum    int

	SpriteSheet  *ebiten.Image
	currFrameNum int
}

func ebitenImageFromPath(path string) *ebiten.Image {
	openFile, err := ebitenutil.OpenFile(path)
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(openFile)
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

func initAnimation() {
	runnerImage = ebitenImageFromPath("./resources/sprites/Squirrel-running.png")
	idleImage = ebitenImageFromPath("./resources/sprites/Squirrel-idle.png")
	attackImage = ebitenImageFromPath("./resources/sprites/Squirrel-attack.png")
	runnerEnemyImage = ebitenImageFromPath("./resources/sprites/enemy-chipmunk-running.png")
	idleEnemyImage = ebitenImageFromPath("./resources/sprites/enemy-chipmunk-idle.png")
	shootingEnemyImage = ebitenImageFromPath("./resources/sprites/enemy-chipmunk-attack.png")
	enemyBulletImage = ebitenImageFromPath("./resources/sprites/enemy-bullet.png")
	portalImage = ebitenImageFromPath("./resources/sprites/portal.png")

	playerSprite = NewAnimatedSprite(
		0,
		0,
		chipmunkSize,
		chipmunkSize,
		5,
		runnerImage)
	playerIdleSprite = NewAnimatedSprite(
		0,
		0,
		chipmunkSize,
		chipmunkSize,
		3,
		idleImage)
	playerAttackSprite = NewAnimatedSprite(
		0,
		0,
		32,
		32,
		6,
		attackImage)
}

// frameOS and frameOY should probably be 0
func NewAnimatedSprite(frameOX int, frameOY int, frameWidth int, frameHeight int, frameNum int, spriteSheet *ebiten.Image) *AnimatedSprite {
	return &AnimatedSprite{
		frameOX:      frameOX,
		frameOY:      frameOY,
		frameWidth:   frameWidth,
		frameHeight:  frameHeight,
		frameNum:     frameNum,
		SpriteSheet:  spriteSheet,
		currFrameNum: 0,
	}
}

func (a *AnimatedSprite) GetCurrFrame() *ebiten.Image {
	sx, sy := a.frameOX+a.currFrameNum*a.frameWidth, 0
	return a.SpriteSheet.SubImage(image.Rect(sx, sy, sx+a.frameWidth, sy+a.frameHeight)).(*ebiten.Image)
}

func (a *AnimatedSprite) NextFrame() {
	if a.currFrameNum+1 >= a.frameNum {
		a.currFrameNum = 0
	} else {
		a.currFrameNum++
	}
}

func (a *AnimatedSprite) ResetSprite() {
	a.currFrameNum = 0
}
