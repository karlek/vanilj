// Package fractal contains utility functions useable when drawing any fractal.
package fractal

import (
	"image/color"
	"math/rand"
	"time"
)

type Gradient []color.Color

var (
	PedagogicalGradient = Gradient{
		color.RGBA{0, 0, 0, 0xff},       // Black.
		color.RGBA{0xff, 0xf0, 0, 0xff}, // Yellow.
		color.RGBA{0, 0, 0xff, 0xff},    // Blue.
		color.RGBA{0, 0xff, 0, 0xff},    // Green.
		color.RGBA{0xff, 0, 0, 0xff},    // Red.
	}
)

func NewRandomGradient(iterations int) Gradient {
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	g := make(Gradient, iterations)
	for n := range g {
		g[n] = color.RGBA{uint8(r.Intn(255)), uint8(r.Intn(255)), uint8(r.Intn(255)), 0xff}
	}
	return g

}

func NewPrettyGradient(iterations int) Gradient {
	g := make(Gradient, iterations)
	var col color.Color
	for n := range g {
		val := uint8(float64(n) / float64(iterations) * 255)
		if n < (iterations / 2) {
			col = color.RGBA{val * 2, 0x00, val * 2, 0xff}
		} else {
			col = color.RGBA{val, val, val, 0xff}
		}
		g[n] = col
	}
	return g
}

func (g Gradient) DivergenceToColor(escapedIn int) color.Color {
	return g[escapedIn%len(g)]
}
