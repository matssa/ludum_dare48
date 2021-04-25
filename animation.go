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
	playerSprite     *AnimatedSprite
	playerIdleSprite *AnimatedSprite
	playerAttackSprite *AnimatedSprite
	runnerEnemyImage   *ebiten.Image
	idleEnemyImage     *ebiten.Image
	shootingEnemyImage     *ebiten.Image
	enemyBulletImage *ebiten.Image
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

func initAnimation() {
	running, err := ebitenutil.OpenFile("./resources/sprites/Squirrel-running.png")
	if err != nil {
		log.Fatal(err)
	}
	idle, err := ebitenutil.OpenFile("./resources/sprites/Squirrel-idle.png")
	if err != nil {
		log.Fatal(err)
	}
	attack, err := ebitenutil.OpenFile("./resources/sprites/Squirrel-attack.png")
	if err != nil {
		log.Fatal(err)
	}
	runningEnemy, err := ebitenutil.OpenFile("./resources/sprites/enemy-chipmunk-running.png")
	if err != nil {
		log.Fatal(err)
	}
	idleEnemy, err := ebitenutil.OpenFile("./resources/sprites/enemy-chipmunk-idle.png")
	if err != nil {
		log.Fatal(err)
	}
	shootingEnemy, err := ebitenutil.OpenFile("./resources/sprites/enemy-chipmunk-attack.png")
	if err != nil {
		log.Fatal(err)
	}
	enemyBullet, err := ebitenutil.OpenFile("./resources/sprites/enemy-bullet.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(running)
	if err != nil {
		log.Fatal(err)
	}
	imgIdle, _, err := image.Decode(idle)
	if err != nil {
		log.Fatal(err)
	}
	imgAttack, _, err := image.Decode(attack)
	if err != nil {
		log.Fatal(err)
	}
	imgEnemy, _, err := image.Decode(runningEnemy)
	if err != nil {
		log.Fatal(err)
	}
	imgIdleEnemy, _, err := image.Decode(idleEnemy)
	if err != nil {
		log.Fatal(err)
	}
	imgShootingEnemy, _, err := image.Decode(shootingEnemy)
	if err != nil {
		log.Fatal(err)
	}
	imgEnemyBullet, _, err := image.Decode(enemyBullet)
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebiten.NewImageFromImage(img)
	idleImage = ebiten.NewImageFromImage(imgIdle)
	attackImage = ebiten.NewImageFromImage(imgAttack)
	runnerEnemyImage = ebiten.NewImageFromImage(imgEnemy)
	idleEnemyImage = ebiten.NewImageFromImage(imgIdleEnemy)
	shootingEnemyImage = ebiten.NewImageFromImage(imgShootingEnemy)
	enemyBulletImage = ebiten.NewImageFromImage(imgEnemyBullet)
	playerSprite = NewAnimatedSprite(
		0,
		0,
		32,
		32,
		5,
		runnerImage)
	playerIdleSprite = NewAnimatedSprite(
		0,
		0,
		32,
		32,
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
	a.currFrameNum = 0;
}
