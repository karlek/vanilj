// 42 is a program to create mandelbrot images.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/karlek/42/canvas"
	"github.com/karlek/42/fractal/mandel"
	"github.com/mewkiz/pkg/errutil"
)

var (
	filename   string
	width      int
	iterations int
	height     int
	zoom       float64
	centerReal float64
	centerImag float64
)

func init() {
	flag.Usage = usage

	flag.StringVar(&filename, "o", "fractal.png", "output filename.")
	flag.IntVar(&width, "width", 1920, "image width.")
	flag.IntVar(&height, "height", 1080, "image height.")
	flag.IntVar(&iterations, "i", 100, "number of iterations.")
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
	mandel.Draw(c.RGBA, zoom, iterations, complex(centerReal, centerImag))

	err = c.Save(filename)
	if err != nil {
		return errutil.Err(err)
	}

	return nil
}