// Code partly stolen from Snilsson with love <3
package julia

import (
	"image"
	"image/color"
	"sync"
)

func DrawJulia(src image.Image, rgba *image.RGBA, zoom float64, center complex128, iterations float64) {
	w, h := float64(rgba.Bounds().Size().X), float64(rgba.Bounds().Size().Y)

	ratio := w / h

	fillMapper(src)

	wg := new(sync.WaitGroup)
	wg.Add(int(w))

	for x := 0.0; x < w; x++ {
		// (x-w/2.0) is used to make the x-axis (and also the origo) run through the middle of the screen, horizontally.
		// (0.2 * zoom * w) is used to make the zoom proportional to the width of the image.
		// However since we want to zoom on the x-/y-axis equally we need to make the real value proportional to the imaginary value.
		// Because 'w' is 'ratio' times bigger than 'h' we multiply the real value with ratio.
		// Now the values are proportional.
		pr := ratio * (x - w/2.0) / (0.2 * zoom * w)
		go roworbit(src, rgba, w, h, x, zoom, iterations, pr, center, wg)
	}
	wg.Wait()
}

// Julia returns an image of size n x n of the Julia set for f.
func Julia(f ComplexFunc, n int) image.Image {
	cols := make(chan column)
	bounds := image.Rect(-n/2, -n/2, n/2, n/2)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		go produceCols(x, n, f, bounds, cols)
	}

	return drawCols(bounds, cols, n)
}

func produceCols(x, n int, f ComplexFunc, bounds image.Rectangle, cols chan<- column) {
	col := make([]color.RGBA, bounds.Max.Y-bounds.Min.Y)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		col[y+n/2] = point(x, y, n, f)
	}
	cols <- column{x, col}
}

func point(x, y, n int, f ComplexFunc) color.RGBA {
	s := float64(n / 4)

	// function point in julia set relative to the image resolution.
	rel := complex(float64(x)/s, float64(y)/s)

	// Julia value, i.e divergence rate.
	jv := Iterate(f, rel, 256)
	return color.RGBA{
		0,
		0,
		uint8(jv % 32 * 8),
		255,
	}
}

type column struct {
	x      int
	pixels []color.RGBA
}

func drawCols(bounds image.Rectangle, cols chan column, n int) (img *image.RGBA) {
	img = image.NewRGBA(bounds)

	num := 0
	for col := range cols {
		if num == bounds.Max.Y-bounds.Min.Y-1 {
			close(cols)
		}
		for y, pixel := range col.pixels {
			img.Set(col.x, y-n/2, pixel)
		}
		num++
	}
	return
}

// Iterate sets z_0 = z, and repeatedly computes z_n = f(z_{n-1}), n â‰¥ 1,
// until |z_n| > 2  or n = max and returns this n.
func Iterate(f ComplexFunc, z complex128, max int) (n int) {
	for ; n < max; n++ {
		if real(z)*real(z)+imag(z)*imag(z) > 4 {
			break
		}
		z = f(z)
	}
	return
}
