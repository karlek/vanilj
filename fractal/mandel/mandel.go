// Package mandel contains functions for computing mandelbrot specific tasks.
// Such as the speed of divergence in a complex point.
package mandel

import (
	"image"
	"image/color"
	"math"
	"math/cmplx"

	"github.com/karlek/vanilj/fractal"
	"github.com/lucasb-eyer/go-colorful"
)

// divergence returns the number of iterations it takes for a complex point to
// leave the mandelbrot set.
func divergence(c complex128, iterations int) int {
	z := complex(0, 0)
	for i := 0; i < iterations; i++ {
		z = z*z + c
		if cmplx.Abs(z) > 2 {
			return i
		}
	}
	return iterations
}

// divergencePrim returns the number of iterations it takes for a complex point
// to leave the mandelbrot set and also returns the last point.
func divergencePrim(c complex128, iterations int) (int, complex128) {
	z := complex(0, 0)
	for i := 0; i < iterations; i++ {
		z = z*z + c
		if cmplx.Abs(z) > 2 {
			return i, z
		}
	}
	return iterations, z
}

// isMemberOfSet determines if the complex point z is member of the mandelbrot
// set.
func isMemberOfSet(z complex128) bool {
	return 2 >= cmplx.Abs(z)
}

// Draw draws the mandelbrot fractal to an image.
func Draw(rgba *image.RGBA, zoom float64, center complex128, iterations int, gradient fractal.Gradient) {
	size := rgba.Bounds().Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			p := complex(float64(size.X/2-x)/zoom-real(center), float64(size.Y/2-y)/zoom+imag(center))

			// Don't draw the points outside the Mandelbrot set.
			if !isMemberOfSet(p) {
				continue
			}

			// Get the speed of divergence.
			mVal := divergence(p, iterations)
			rgba.Set(x, y, gradient.DivergenceToColor(mVal))
		}
	}
}

// DrawSmooth draws the mandelbrot fractal to an image with smooth coloring.
func DrawSmooth(rgba *image.RGBA, zoom float64, center complex128, iterations int) {
	size := rgba.Bounds().Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			p := complex(float64(size.X/2-x)/zoom-real(center), float64(size.Y/2-y)/zoom+imag(center))

			// Don't draw the points outside the Mandelbrot set.
			if !isMemberOfSet(p) {
				continue
			}

			// Get the speed of divergence.
			mVal, z := divergencePrim(p, iterations)
			nsmooth := (float64(mVal) + float64(1) - math.Log(math.Log(cmplx.Abs(z)))/math.Log(1.5)) / float64(iterations)
			rgba.Set(x, y, smoothColor(nsmooth, mVal))
		}
	}
}

func smoothColor(val float64, mVal int) color.Color {
	return colorful.Hsv(val*360, 1, 1)
}
