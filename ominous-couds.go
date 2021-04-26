package main

import (
	"container/list"
	"image"
	// 	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var smokeImage *ebiten.Image

const (
	CLOUDS_MOVEMENT_DAMPING = 0.3
	CLOUDS_MAX_NUM          = 500
)

func init() {
	myimg, err := ebitenutil.OpenFile("./resources/sprites/ominous-cloud.png")
	img, _, err := image.Decode(myimg)
	if err != nil {
		log.Fatal(err)
	}
	smokeImage = ebiten.NewImageFromImage(img)

}

func randInRange(min int, max int) int {
	return rand.Intn(max-min) + min
}

type OminousCloud struct {
	count    int
	maxCount int
	dir      float64
	posx     int
	posy     int

	img   *ebiten.Image
	scale float64
	angle float64
	alpha float64
}

type OminousClouds struct {
	clouds          *list.List
	createNewClouds bool
}

func (s *OminousCloud) update() {
	if s.count == 0 {
		return
	}
	s.count--
}

func (s *OminousCloud) terminated() bool {
	return s.count == 0
}

func (s *OminousCloud) draw(screen *ebiten.Image) {
	if s.count == 0 {
		return
	}

	const (
		ox = screenWidth / 2
		oy = screenHeight / 2
	)
	x := math.Cos(s.dir) * (float64(s.maxCount-s.count) * CLOUDS_MOVEMENT_DAMPING)
	y := math.Sin(s.dir) * (float64(s.maxCount-s.count) * CLOUDS_MOVEMENT_DAMPING)

	op := &ebiten.DrawImageOptions{}

	sx, sy := s.img.Size()
	op.GeoM.Translate(-float64(sx)/2, -float64(sy)/2)
	op.GeoM.Rotate(s.angle)
	op.GeoM.Scale(s.scale, s.scale)
	op.GeoM.Translate(x, y)
	op.GeoM.Translate(ox, oy)

	op.GeoM.Translate(float64(s.posx), float64(s.posy))

	rate := float64(s.count) / float64(s.maxCount)
	alpha := 0.0
	if rate < 0.2 {
		alpha = rate / 0.2
	} else if rate > 0.8 {
		alpha = (1 - rate) / 0.2
	} else {
		alpha = 1
	}
	alpha *= s.alpha
	op.ColorM.Scale(1, 1, 1, alpha)

	screen.DrawImage(s.img, op)
}

func NewOminousCloud(img *ebiten.Image) *OminousCloud {
	c := rand.Intn(50) + CLOUDS_MAX_NUM
	dir := rand.Float64() * 2 * math.Pi
	a := rand.Float64() * 2 * math.Pi
	s := rand.Float64()*1 + 0.4
	posx := randInRange(-300, 300)
	posy := randInRange(-300, 300)
	return &OminousCloud{
		img: img,

		maxCount: c,
		count:    c,
		dir:      dir,

		posx: posx,
		posy: posy,

		angle: a,
		scale: s,
		alpha: 0.1,
	}
}

func (o *OminousClouds) UpdateClouds() error {
	if o.clouds == nil {
		o.clouds = list.New()
	}

	if o.createNewClouds {
		if o.clouds.Len() < 500 && rand.Intn(4) < 3 {
			// Emit
			o.clouds.PushBack(NewOminousCloud(smokeImage))
		}
	}

	for e := o.clouds.Front(); e != nil; e = e.Next() {
		s := e.Value.(*OminousCloud)
		s.update()
		if s.terminated() {
			defer o.clouds.Remove(e)
		}
	}
	return nil
}

func (o OminousClouds) DrawClouds(screen *ebiten.Image) {
	for e := o.clouds.Front(); e != nil; e = e.Next() {
		s := e.Value.(*OminousCloud)
		s.draw(screen)
	}
}

func (o *OminousClouds) StopClouds() {
	o.createNewClouds = false
}
func (o *OminousClouds) StartClouds() {
	o.createNewClouds = true
}
