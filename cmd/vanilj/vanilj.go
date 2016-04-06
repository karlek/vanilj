// Vanilj is a mandelbrot explorer.
package main

import (
	"flag"
	"fmt"
	"image/color"
	// "image/png"
	"log"
	"os"
	"runtime"

	"github.com/karlek/profile"
	"github.com/karlek/vanilj/canvas"
	"github.com/karlek/vanilj/fractal"
	"github.com/karlek/vanilj/fractal/mandel"
	"github.com/mewkiz/pkg/errutil"
)

var (
	filename    string
	colorScheme string
	width       int
	height      int
	iterations  float64
	zoom        float64
	centerReal  float64
	centerImag  float64
)

func init() {
	flag.Usage = usage

	flag.StringVar(&filename, "o", "fractal.png", "output filename.")
	flag.IntVar(&width, "width", 1920, "image width.")
	flag.IntVar(&height, "height", 1080, "image height.")
	flag.Float64Var(&iterations, "i", 1000, "number of iterations.")
	flag.Float64Var(&zoom, "z", 1, "zoom level.")
	flag.Float64Var(&centerReal, "cr", 0, "real value of center offset.")
	flag.Float64Var(&centerImag, "ci", 0, "imaginary value of center offset.")
	flag.StringVar(&colorScheme, "t", "smooth", "color scheme")
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [OPTIONS]...\n\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	defer profile.Start(profile.CPUProfile).Stop()

	err := renderMandelbrot()
	if err != nil {
		log.Fatalln(err)
	}
}

func renderMandelbrot() (err error) {
	flag.Parse()

	c := canvas.NewCanvas(width, height)
	for x := 0; x < c.Bounds().Size().X; x++ {
		for y := 0; y < c.Bounds().Size().Y; y++ {
			c.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}
	switch colorScheme {
	case "smooth":
		f := fractal.Fractal{
			Src:    c.RGBA,
			Iter:   iterations,
			Center: complex(centerReal, centerImag),
			Zoom:   zoom,
		}
		mandel.Smooth(&f)
	case "random":
		mandel.Draw(c.RGBA, zoom, complex(centerReal, centerImag), iterations, fractal.NewRandomGradient(iterations))
	case "pretty":
		mandel.Draw(c.RGBA, zoom, complex(centerReal, centerImag), iterations, fractal.NewPrettyGradient(iterations))
	case "pedagogical":
		mandel.Draw(c.RGBA, zoom, complex(centerReal, centerImag), iterations, fractal.PedagogicalGradient)
	case "orbit":
		// f, err := os.Open("z.png")
		// if err != nil {
		// 	return err
		// }
		// defer f.Close()
		// src, err := png.Decode(f)
		// if err != nil {
		// 	return err
		// }
		// mandel.DrawOrbit(src, c.RGBA, zoom, complex(centerReal, centerImag), iterations)
	case "julia":
		// mandel.DrawJulia(src, c.RGBA, zoom, complex(centerReal, centerImag), iterations)
	default:
		log.Fatalln(errutil.Newf("undefined color scheme: %s", colorScheme))
	}

	err = c.Save(filename)
	if err != nil {
		return errutil.Err(err)
	}

	return nil
}
