// Vanilj is a mandelbrot explorer.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/karlek/vanilj/canvas"
	// "github.com/karlek/vanilj/fractal"
	"github.com/karlek/vanilj/fractal/mandel"
	"github.com/mewkiz/pkg/errutil"
)

var (
	filename   string
	width      int
	height     int
	iterations float64
	zoom       float64
	centerReal float64
	centerImag float64
)

func init() {
	flag.Usage = usage

	flag.StringVar(&filename, "o", "fractal.png", "output filename.")
	flag.IntVar(&width, "width", 1920, "image width.")
	flag.IntVar(&height, "height", 1080, "image height.")
	flag.Float64Var(&iterations, "i", 100, "number of iterations.")
	flag.Float64Var(&zoom, "z", 250, "zoom level.")
	flag.Float64Var(&centerReal, "cr", 0, "real value of center offset.")
	flag.Float64Var(&centerImag, "ci", 0, "imaginary value of center offset.")
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [OPTIONS]...\n\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	err := renderMandelbrot()
	if err != nil {
		log.Fatalln(err)
	}
}

func renderMandelbrot() (err error) {
	flag.Parse()

	c := canvas.NewCanvas(width, height)

	// mandel.Draw(c.RGBA, zoom, complex(centerReal, centerImag), iterations, fractal.NewRandomGradient(iterations))
	// mandel.Draw(c.RGBA, zoom, complex(centerReal, centerImag), iterations, fractal.NewPrettyGradient(iterations))
	// mandel.Draw(c.RGBA, zoom, complex(centerReal, centerImag), iterations, fractal.PedagogicalGradient)
	mandel.DrawSmooth(c.RGBA, zoom, complex(centerReal, centerImag), iterations)

	err = c.Save(filename)
	if err != nil {
		return errutil.Err(err)
	}

	return nil
}
