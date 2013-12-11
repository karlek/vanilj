// Package fractal contains utility functions useable when drawing any fractal.
package fractal

import (
	"image/color"
	"math/rand"
	"time"
)

// Gradient is a list of colors.
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

// NewRandomGradient creates a gradient of colors proportional to the number of iterations.
func NewRandomGradient(iterations float64) Gradient {
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	g := make(Gradient, int64(iterations))
	for n := range g {
		g[n] = color.RGBA{uint8(r.Intn(255)), uint8(r.Intn(255)), uint8(r.Intn(255)), 0xff}
	}
	return g

}

// NewPrettyGradient creates a gradient of colors fading between purple and white. The smoothness is proportional to the number of iterations
func NewPrettyGradient(iterations float64) Gradient {
	g := make(Gradient, int64(iterations))
	var col color.Color
	for n := range g {
		val := uint8(float64(n) / float64(iterations) * 255)
		if int64(n) < int64(iterations/2) {
			col = color.RGBA{val * 2, 0x00, val * 2, 0xff}
		} else {
			col = color.RGBA{val, val, val, 0xff}
		}
		g[n] = col
	}
	return g
}

// DivergenceToColor returns a color depending on the number of iterations it took for the fractal to escape the fractal set.
func (g Gradient) DivergenceToColor(escapedIn int) color.Color {
	return g[escapedIn%len(g)]
}
