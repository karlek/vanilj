// Package mandel contains functions for computing mandelbrot specific tasks.
// Such as the speed of divergence in a complex point.
package mandel

import (
	"image"
	"math/cmplx"

	"github.com/karlek/42/fractal"
)

// Divergence returns the number of iterations it takes for a complex point to
// leave the mandelbrot set.
func Divergence(c complex128, iterations int) int {
	z := complex(0, 0)
	for i := 0; i < iterations; i++ {
		z = z*z + c
		if cmplx.Abs(z) > 2 {
			return i
		}
	}
	return iterations
}

// IsMemberOfSet determines if the complex point z is member of the mandelbrot
// set.
func IsMemberOfSet(z complex128) bool {
	return 2 >= cmplx.Abs(z)
}

// Draw draws the mandelbrot fractal to an image.
func Draw(rgba *image.RGBA, zoom float64, iterations int, center complex128) {
	g := fractal.NewRandomGradient()
	size := rgba.Bounds().Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			p := complex(float64(size.X/2-x)/zoom-real(center), float64(size.Y/2-y)/zoom+imag(center))

			// Don't draw the points outside the Mandelbrot set.
			if !IsMemberOfSet(p) {
				continue
			}

			// Get the speed of divergence.
			m := Divergence(p, iterations)
			rgba.Set(x, y, g.DivergenceToColor(m, iterations))
		}
	}
}
