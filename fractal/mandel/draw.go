// Package mandel contains functions for computing mandelbrot specific tasks.
// Such as the speed of divergence in a complex point.
package mandel

import (
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/karlek/vanilj/fractal"
	"github.com/lucasb-eyer/go-colorful"
)

func Smooth(f *fractal.Fractal) {
	h := float64(f.Src.Bounds().Size().Y)
	wg := new(sync.WaitGroup)

	// For each row.
	for y := 0.0; y < h; y++ {
		wg.Add(1)
		go func(f *fractal.Fractal, y float64, wg *sync.WaitGroup) {
			row(f, y)
			wg.Done()
		}(f, y, wg)
	}
	wg.Wait()
}

func row(f *fractal.Fractal, y float64) {
	w, h := float64(f.Src.Bounds().Dx()), float64(f.Src.Bounds().Dy())
	for x := 0.0; x < w; x++ {
		ratio := w / h
		pr := ratio * (x - w/2.0) / (0.2 * f.Zoom * w)
		pi := (y - h/2.0) / (0.2 * f.Zoom * h)
		p := complex(pr, pi) + f.Center

		// Don't draw the points outside the Mandelbrot set.
		if !isMemberOfSet(p) {
			continue
		}
		calc(f, p, int(x), int(y))
	}
}

func calc(f *fractal.Fractal, p complex128, x, y int) {
	abs := func(z complex128) float64 {
		return real(z)*real(z) + imag(z)*imag(z)
	}

	// Get the speed of divergence.
	escape, z := divergence(p, f.Iter)
	f.Src.Set(x, y, smoothColor((escape+1.0-math.Log(math.Log(abs(z)))/math.Log(2))/f.Iter))
}

/// Add credit.
// smoothColor returns a color from the smooth color formula.
func smoothColor(val float64) color.Color {
	return colorful.Hsv(val*360, 1, 1)
}

// Draw draws the mandelbrot fractal to an image.
func Draw(rgba *image.RGBA, zoom float64, center complex128, iterations float64, gradient fractal.Gradient) {
	wg := new(sync.WaitGroup)
	w, h := float64(rgba.Bounds().Size().X), float64(rgba.Bounds().Size().Y)
	wg.Add(int(w))
	ratio := w / h
	for x := 0.0; x < w; x++ {
		// (x-w/2.0) is used to make the x-axis (and also the origo) run through the middle of the screen, horizontally.
		// (0.2 * zoom * w) is used to make the zoom proportional to the width of the image.
		// However since we want to zoom on the x-/y-axis equally we need to make the real value proportional to the imaginary value.
		// Because 'w' is 'ratio' times bigger than 'h' we multiply the real value with ratio.
		// Now the values are proportional.
		pr := ratio * (x - w/2.0) / (0.2 * zoom * w)
		go rowprim(rgba, w, h, x, zoom, iterations, pr, center, gradient, wg)
	}
	wg.Wait()
}

func rowprim(rgba *image.RGBA, w, h, x, zoom, iterations, pr float64, center complex128, gradient fractal.Gradient, wg *sync.WaitGroup) {
	// We multiply the real value with the ratio since we divide by the boundary values (width and height).
	// Division is equal to zoom. Since 'w' is 'ratio' times bigger than 'h'. If we didn't multiply, the image wouldn't zoom equally on the x- and y-axis rendering
	// strange images.
	for y := 0.0; y < h; y++ {
		// (y - h/2.0) is used to make the y-axis run through the middle of the screen, vertically.
		pi := (y - h/2.0) / (0.2 * zoom * h)

		// Center is the complex point were we will zoom in on the mandelbrot.
		p := complex(pr, pi) + center

		// Don't draw the points outside the Mandelbrot set.
		if !isMemberOfSet(p) {
			continue
		}

		// Get the speed of divergence.
		mVal, _ := divergence(p, iterations)
		rgba.Set(int(x), int(y), gradient.DivergenceToColor(int(mVal)))
	}
	wg.Done()
}
