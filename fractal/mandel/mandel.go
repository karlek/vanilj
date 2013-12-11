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
func divergence(c complex128, iterations float64) float64 {
	z := complex(0, 0)
	for i := 0.0; i < iterations; i++ {
		z = z*z + c
		if !isMemberOfSet(z) {
			return i
		}
	}
	return iterations
}

// divergencePrim returns the number of iterations it takes for a complex point
// to leave the mandelbrot set and also returns the point last point (which could be outside the mandelbrot set).
func divergencePrim(c complex128, iterations float64) (float64, complex128) {
	z := complex(0, 0)
	for i := 0.0; i < iterations; i++ {
		z = z*z + c
		if !isMemberOfSet(z) {
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
func Draw(rgba *image.RGBA, zoom float64, center complex128, iterations float64, gradient fractal.Gradient) {
	w, h := float64(rgba.Bounds().Size().X), float64(rgba.Bounds().Size().Y)
	for x := 0.0; x < w; x++ {
		for y := 0.0; y < h; y++ {				
			// We multiply the real value with the ratio since we divide by the boundary values (width and height).
			// Division is equal to zoom. Since 'w' is 'ratio' times bigger than 'h'. If we didn't multiply, the image wouldn't zoom equally on the x- and y-axis rendering
			// strange images.
			ratio := w / h
			
			// (x-w/2.0) is used to make the x-axis (and also the origo) run through the middle of the screen, horizontally.
			// (0.2 * zoom * w) is used to make the zoom proportional to the width of the image.
			// However since we want to zoom on the x-/y-axis equally we need to make the real value proportional to the imaginary value.
			// Because 'w' is 'ratio' times bigger than 'h' we multiply the real value with ratio.
			// Now the values are proportional.
			pr := ratio * (x - w/2.0) / (0.2 * zoom * w)
			
			// (y - h/2.0) is used to make the y-axis run through the middle of the screen, vertically.
			pi := (y - h/2.0) / (0.2 * zoom * h)
			
			// Center is the complex point were we will zoom in on the mandelbrot. 
			p := complex(pr, pi) + center

			// Don't draw the points outside the Mandelbrot set.
			if !isMemberOfSet(p) {
				continue
			}

			// Get the speed of divergence.
			mVal := divergence(p, iterations)
			rgba.Set(int(x), int(y), gradient.DivergenceToColor(int(mVal)))
		}
	}
}

// DrawSmooth draws the mandelbrot fractal to an image with smooth coloring.
func DrawSmooth(rgba *image.RGBA, zoom float64, center complex128, iterations float64) {
	w, h := float64(rgba.Bounds().Size().X), float64(rgba.Bounds().Size().Y)
	for y := 0.0; y < h; y++ {
		for x := 0.0; x < w; x++ {
			ratio := w / h
			pr := ratio * (x - w/2.0) / (0.2 * zoom * w)
			pi := (y - h/2.0) / (0.2 * zoom * h)
			p := complex(pr, pi) + center

			// Don't draw the points outside the Mandelbrot set.
			if !isMemberOfSet(p) {
				continue
			}

			// Get the speed of divergence.
			mVal, z := divergencePrim(p, iterations)
			nsmooth := (float64(mVal) + float64(1) - math.Log(math.Log(cmplx.Abs(z)))/math.Log(2)) / iterations
			rgba.Set(int(x), int(y), smoothColor(nsmooth))
		}
	}
}

/// Add credit.
// smoothColor returns a color from the smooth color formula.
func smoothColor(val float64) color.Color {
	return colorful.Hsv(val*360, 1, 1)
}
