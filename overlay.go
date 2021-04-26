package main

// Copyright 2020 The Ebiten Authors
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


import (
	"fmt"
	"image/color"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const maxHealth = 100

var (
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func createInfoString(health int, maxHealth int, enemies int) string {
	return fmt.Sprintf("health %d / %d, enemies left: %d", health, maxHealth, enemies)
}

func DrawOverlay(screen *ebiten.Image, health int, enemies int) {
	transparentGray := color.RGBA{0x80, 0x80, 0x80, 0x55}
	{
		const x, y = 10, 20
		myText := createInfoString(health, maxHealth, enemies)
		b := text.BoundString(mplusNormalFont, myText)
		ebitenutil.DrawRect(screen, float64(b.Min.X+x), float64(b.Min.Y+y), float64(b.Dx()), float64(b.Dy()), transparentGray)
		text.Draw(screen, myText, mplusNormalFont, x, y, color.White)
	}
	{
		_, height := screen.Size()
		x, y := 20, height-10
		myText := "move: wasd, attack: y, camera: v, debug: c"
		b := text.BoundString(mplusNormalFont, myText)
		ebitenutil.DrawRect(screen, float64(b.Min.X+x), float64(b.Min.Y+y), float64(b.Dx()), float64(b.Dy()), transparentGray)
		text.Draw(screen, myText, mplusNormalFont, x, y, color.White)
	}
}
