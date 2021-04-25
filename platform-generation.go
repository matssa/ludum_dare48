package main

import (
	"fmt"
	"time"
	"math/rand"
)


func init() {
	rand.Seed(time.Now().UnixNano())
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
    for j := 1; j < length - 1; j++ {
	    myTiles = append(myTiles, NewTile(
		    startX + (j * 16),
		    startY,
		    "top"))
    }

    lastLineLength := length - 1
    myTiles = append(myTiles, NewTile(
	    startX + ((lastLineLength) * 16),
	    startY,
	    "top-right"))
    return myTiles;
}

func createMap() [][]*Tile {
	lines := make([][]*Tile, 0, 0)
	startLine := make([]*Tile, 0, 0)
	for i := 0; i < 10; i++ {
		startLine = append(startLine, NewTile(
			0 + i *16,
			300,
			"top"))
	}
	lines = append(lines, startLine)

	numLines := rand.Intn(100 - 80) + 80;

	lineLengths := make([]int, 0, 0);
	for i := 0; i < numLines; i++ {
		lineLengths = append(lineLengths, rand.Intn(10 - 3) + 3)
	}
	fmt.Printf("lengths %v", lineLengths);

	prevx := 0
	prevy := 300
	lineXStarts := make([]int, 0, 0)
	lineYStarts := make([]int, 0, 0)
	for i := 0; i < numLines; i++ {
		maxNewX := (prevx + (lineLengths[i] * 16)) + (5 * 16)
		minNewX := (prevx + (lineLengths[i] * 16)) - (3 * 16)
		maxNewY := prevy + (3 * 16)
		minNewY := prevy - (3 * 16)

		borderLimit := 20

		var newX int
		if minNewX < borderLimit {
			newX = maxNewX
		} else if maxNewX > screenWidth - borderLimit {
			newX = minNewX
		} else {
			newX = rand.Intn(maxNewX - minNewX) + minNewX
		}
		var newY int
		if minNewY < borderLimit {
			newY = maxNewY
		} else if maxNewY > screenHeight - borderLimit {
			newY = minNewY
		} else {
			newY = rand.Intn(maxNewY - minNewY) + minNewY
		}

		lineXStarts = append(lineXStarts, newX)
		lineYStarts = append(lineYStarts, newY)
		prevx = newX
		prevy = newY
	}
	fmt.Printf("startx %v", lineXStarts);
	fmt.Printf("starty %v", lineYStarts);

	for i := 0; i < numLines; i++ {
		lines = append(lines, createLine(lineLengths[i], lineXStarts[i], lineYStarts[i]))
	}


	return lines
}
