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

func NewRandomGradient() Gradient {
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	g := Gradient{
		color.RGBA{uint8(r.Intn(255)), uint8(r.Intn(255)), uint8(r.Intn(255)), 0xff},
		color.RGBA{uint8(r.Intn(255)), uint8(r.Intn(255)), uint8(r.Intn(255)), 0xff},
		color.RGBA{uint8(r.Intn(255)), uint8(r.Intn(255)), uint8(r.Intn(255)), 0xff},
		color.RGBA{uint8(r.Intn(255)), uint8(r.Intn(255)), uint8(r.Intn(255)), 0xff},
		color.RGBA{uint8(r.Intn(255)), uint8(r.Intn(255)), uint8(r.Intn(255)), 0xff},
	}
	return g
}

func (g Gradient) DivergenceToColor(escapedIn, iterations int) color.Color {
	// If escapedIn is equal to the number of iterations it can be considered a
	// member of the fractal set.
	switch {
	case (escapedIn == iterations):
		return g[0]
	case (escapedIn%4 == 0):
		return g[1]
	case (escapedIn%3 == 0):
		return g[2]
	case (escapedIn%2 == 0):
		return g[3]
	case (escapedIn%2 == 1):
		return g[4]
	default:
		return g[0]
	}
}
