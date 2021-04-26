package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func roundToTileSize(num int) int {
	return int(math.Round(float64(num)/16.0) * 16)
}

func createLine(length int, startX int, startY int) []*Tile {
	myTiles := make(
		[]*Tile,
		0,
		0)
	myTiles = append(myTiles, NewTile(
		startX,
		startY,
		"top-left"))
	for j := 1; j < length-1; j++ {
		myTiles = append(myTiles, NewTile(
			startX+(j*16),
			startY,
			"top"))
	}

	lastLineLength := length - 1
	myTiles = append(myTiles, NewTile(
		startX+((lastLineLength)*16),
		startY,
		"top-right"))
	return myTiles
}

func createMap() [][]*Tile {
	firstLineX := 500
	firstLineY := 1200
	firstLineLength := 10
	lines := make([][]*Tile, 0, 0)
	startLine := make([]*Tile, 0, 0)
	for i := 0; i < firstLineLength; i++ {
		startLine = append(startLine, NewTile(
			firstLineX+i*16,
			firstLineY,
			"top"))
	}
	lines = append(lines, startLine)

	numLines := rand.Intn(200-150) + 150


	prevLineLength := firstLineLength
	prevx := firstLineX
	prevy := firstLineY
	lineXStarts := make([]int, 0, 0)
	lineYStarts := make([]int, 0, 0)
	for i := 0; i < numLines; i++ {
		lineLength := rand.Intn(20-3)+3
		maxNewX := (prevx + (prevLineLength * 16)) + (5 * 16)
		minNewX := (prevx + (prevLineLength * 16)) - (3 * 16)
		maxNewY := prevy + (4 * 16)
		minNewY := prevy - (4 * 16)

		borderLimit := 300

		var newX int
		if minNewX < borderLimit {
			newX = maxNewX
		} else if maxNewX > worldWidth-borderLimit {
			newX = minNewX
		} else {
			newX = rand.Intn(maxNewX-minNewX) + minNewX
		}
		var newY int
		if minNewY < borderLimit {
			newY = maxNewY
		} else if maxNewY > (worldHeight-400)-borderLimit {
			newY = minNewY
		} else {
			newY = rand.Intn(maxNewY-minNewY) + minNewY
		}

		newX = roundToTileSize(newX);
		newY = roundToTileSize(newY);

		lineXStarts = append(lineXStarts, newX)
		lineYStarts = append(lineYStarts, newY)
		prevx = newX
		prevy = newY
		prevLineLength = lineLength
		lines = append(lines, createLine(lineLength, newX, newY))

		extraLines := rand.Intn(5);
		for i := 0; i < extraLines; i++ {
			extraLineLength := rand.Intn(25-1)+2
			extraMaxNewX := (prevx + (extraLineLength * 16)) + (25 * 16)
			extraMinNewX := (prevx + (extraLineLength * 16)) - (25 * 16)
			extraMaxNewY := prevy + (25 * 16)
			extraMinNewY := prevy - (25 * 16)

			extraNewX := rand.Intn(extraMaxNewX-extraMinNewX) + extraMinNewX
			if extraNewX < borderLimit || extraNewX > worldWidth-borderLimit {
				fmt.Printf("happened x\n");
				continue
			}
			extraNewY := rand.Intn(extraMaxNewY-extraMinNewY) + extraMinNewY
			if extraNewY < borderLimit || extraNewY > (worldHeight-400-borderLimit) {
				fmt.Printf("happened x\n");
				continue;
			}
			extraNewX = roundToTileSize(extraNewX)
			extraNewY = roundToTileSize(extraNewY)
			lines = append(lines, createLine(extraLineLength, extraNewX, extraNewY))
		}
	}

	return lines
}
