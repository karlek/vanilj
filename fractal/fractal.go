// Package fractal contains utility functions useable when drawing any fractal.
package fractal

import (
	"image"
	"image/color"
	"math/rand"
	"time"
)

// Fractal is a general fractal representation. Have information about:
// iterations, zoom value and centering point.
type Fractal struct {
	Src    *image.RGBA // Image to write fractal on.
	Iter   float64     // Number of iterations to perform.
	Center complex128  // Point to center the fractal on.
	Zoom   float64     // Zoom value.
}

// Gradient is a list of colors.
type Gradient []color.Color

var (
	// PedagogicalGradient have a fixed transformation between colors for easier
	// visualization of divergence.s
	PedagogicalGradient = Gradient{
		color.RGBA{0, 0, 0, 0xff},       // Black.
		color.RGBA{0xff, 0xf0, 0, 0xff}, // Yellow.
		color.RGBA{0, 0, 0xff, 0xff},    // Blue.
		color.RGBA{0, 0xff, 0, 0xff},    // Green.
		color.RGBA{0xff, 0, 0, 0xff},    // Red.
	}
)

// NewRandomGradient creates a gradient of colors proportional to the number of
// iterations.
func NewRandomGradient(iterations float64) Gradient {
	seed := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	grad := make(Gradient, int64(iterations))
	for n := range grad {
		grad[n] = randomColor(seed)
	}
	return grad
}

// randomColor returns a random RGB color from a random seed.
func randomColor(seed *rand.Rand) color.RGBA {
	return color.RGBA{
		uint8(seed.Intn(255)),
		uint8(seed.Intn(255)),
		uint8(seed.Intn(255)),
		0xff} // No alpha.
}

// NewPrettyGradient creates a gradient of colors fading between purple and
// white. The smoothness is proportional to the number of iterations
func NewPrettyGradient(iterations float64) Gradient {
	grad := make(Gradient, int64(iterations))
	var col color.Color
	for n := range grad {
		// val ranges from [0..255]
		val := uint8(float64(n) / float64(iterations) * 255)
		if int64(n) < int64(iterations/2) {
			col = color.RGBA{val * 2, 0x00, val * 2, 0xff} // Shade of purple.
		} else {
			col = color.RGBA{val, val, val, 0xff} // Shade of white.
		}
		grad[n] = col
	}
	return grad
}

// DivergenceToColor returns a color depending on the number of iterations it
// took for the fractal to escape the fractal set.
func (g Gradient) DivergenceToColor(escapedIn int) color.Color {
	return g[escapedIn%len(g)]
}

const (
	Modulo int = iota
	IterationCount
)

func (g *Gradient) Get(i, it, mode int) (float64, float64, float64) {
	switch mode {
	case Modulo:
		return g.modulo(i)
	case IterationCount:
		return g.iteration(i, it)
	default:
		return g.modulo(i)
	}
}

func (grad *Gradient) modulo(i int) (float64, float64, float64) {
	if i >= len(*grad) {
		i %= len(*grad)
	}
	r, g, b, _ := (*grad)[i].RGBA()
	return float64(r>>8) / 256, float64(g>>8) / 256, float64(b>>8) / 256
}

var Keys = []int{
	0,
	0,
	0,
}

func (grad *Gradient) iteration(i, it int) (float64, float64, float64) {
	ranges := []float64{
		0.05,
		0.15,
		0.25,
	}
	key := 0
	for rId := len(ranges) - 1; rId >= 0; rId-- {
		if float64(i)/float64(it) >= ranges[rId] {
			key = rId
			break
		}
	}
	Keys[key] += i
	r, g, b, _ := (*grad)[key].RGBA()
	return float64(r>>8) / 256, float64(g>>8) / 256, float64(b>>8) / 256
}

func (g *Gradient) AddColor(c color.Color) {
	(*g) = append((*g), c)
}
