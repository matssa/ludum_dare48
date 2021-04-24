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
	animatedSprite     *AnimatedSprite
	animatedIdleSprite *AnimatedSprite
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
	img, _, err := image.Decode(running)
	if err != nil {
		log.Fatal(err)
	}
	imgIdle, _, err := image.Decode(idle)
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebiten.NewImageFromImage(img)
	idleImage = ebiten.NewImageFromImage(imgIdle)
	animatedSprite = NewAnimatedSprite(
		0,
		0,
		32,
		32,
		5,
		runnerImage)
	animatedIdleSprite = NewAnimatedSprite(
		0,
		0,
		32,
		32,
		3,
		idleImage)
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
