package main

import (
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
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
