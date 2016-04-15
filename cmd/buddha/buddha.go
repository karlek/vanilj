// buddha renders buddhabrot fractals.
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"os"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/disintegration/imaging"
	"github.com/dustin/randbo"

	"github.com/karlek/profile"
)

var (
	// Color scaling.
	exposure float64
	factor   float64
	f        func(float64) float64
	fun      string
	// File options.
	filename string
	load     bool
	save     bool
	rotate   bool
	tries    int
	bailout  float64

	offsetReal float64
	offsetImag float64
	offset     complex128
)

const (
	width      = 4096
	height     = 4096
	iterations = 10000000
	// Camera options.
	zoom = float64(width) / 2.8
)

func init() {
	rand.Seed(time.Now().UnixNano())
	flag.BoolVar(&load, "load", false, "use pre-computed values.")
	flag.BoolVar(&save, "save", false, "save orbits.")
	flag.BoolVar(&rotate, "rotate", false, "rotate the fractal to an upright position.")
	flag.Float64Var(&exposure, "exposure", 1.0, "over exposure")
	flag.Float64Var(&factor, "factor", 1.0, "factor")
	flag.StringVar(&fun, "function", "exp", "color scaling function")
	flag.StringVar(&filename, "out", "a.jpeg", "output filename")
	flag.IntVar(&tries, "tries", 12000000, "number of orbits attempts")
	flag.Float64Var(&bailout, "bail", 4, "bailout value")
	flag.Float64Var(&offsetReal, "real", 0.4, "offsetReal")
	flag.Float64Var(&offsetImag, "imag", 0, "offsetImag")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "%s [OPTIONS],,,", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	switch fun {
	case "exp":
		f = exp
	case "log":
		f = log
	case "sqrt":
		f = sqrt
	case "lin":
		f = lin
	default:
		logrus.Fatalln("invalid color scaling function:", fun)
	}
	offset = complex(offsetReal, offsetImag)

	if err := buddha(); err != nil {
		logrus.Fatalln(err)
	}
}

// Initialize allocates memory for our image and histograms.
func initialize() (img *image.RGBA, r, g, b *Histo) {
	// Output image with black background.
	return image.NewRGBA(image.Rect(0, 0, width, height)), &Histo{}, &Histo{}, &Histo{}
}

func buddha() (err error) {
	defer profile.Start(profile.CPUProfile).Stop()

	logrus.Println("[.] Initializing.")
	img, r, g, b := initialize()

	if load {
		logrus.Println("[-] Loading visits.")
		r, g, b, err = loadVisits()
		if err != nil {
			return err
		}
		plot(img, r, g, b)
		return render(img)
	}

	logrus.Println("[-] Calculating visited points.")

	workers := 32
	wg := new(sync.WaitGroup)
	wg.Add(workers)
	for n := 0; n < workers; n++ {
		incChan := make(chan orbit, tries/workers/1000)
		go itoc(r, g, b, incChan)
		go arbitrary(tries/workers, incChan, wg)
	}
	wg.Wait()

	if save {
		logrus.Println("[i] Saving r, g, b channels")
		if err := gobVisits(r, g, b); err != nil {
			return err
		}
	}

	logrus.Println("[/] Creating image.")
	plot(img, r, g, b)
	return render(img)
}

// arbitrary will try to find orbits from n random complex points.
func arbitrary(n int, incChan chan orbit, wg *sync.WaitGroup) {
	var random = randbo.NewFrom(rand.NewSource(rand.Int63()))
	var send [iterations]image.Point
	for i := 0; i < n; i++ {
		c := complex(sign(random)*2*randfloat(random), sign(random)*2*randfloat(random))
		findOrbit(c, incChan, &send)
	}
	wg.Done()
	close(incChan)
}

func findOrbit(c complex128, incChan chan orbit, send *[iterations]image.Point) {
	points := escaped(c)
	num := 0
	for _, z := range points {
		p := ptoc(z)
		// Ignore points outside image.
		if p.X >= width || p.Y >= height || p.X < 0 || p.Y < 0 {
			continue
		}
		send[num] = p
		num++
	}
	incChan <- orbit{send[:num], len(points)}
}

// orbit is a collection of points which escaped after an iteration.
type orbit struct {
	ps []image.Point
	it int
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func plot(img *image.RGBA, r, g, b *Histo) {
	rMax := max(r)
	gMax := max(g)
	bMax := max(b)
	logrus.Println("[i] Visitations:", rMax, gMax, bMax)
	logrus.Printf("[i] Function: %s, factor: %.2f, exposure: %.2f", getFunctionName(f), factor, exposure)
	for x, col := range r {
		for y := range col {
			if r[x][y] == 0 &&
				g[x][y] == 0 &&
				b[x][y] == 0 {
				continue
			}
			c := color.RGBA{
				uint8(value(r[x][y], rMax)),
				uint8(value(g[x][y], gMax)),
				uint8(value(b[x][y], bMax)),
				255}
			img.Set(x, y, c)
		}
	}
}

func exp(x float64) float64 {
	return (1 - math.Exp(-factor*x))
}
func log(x float64) float64 {
	return math.Log1p(factor * x)
}
func sqrt(x float64) float64 {
	return math.Sqrt(factor * x)
}
func lin(x float64) float64 {
	return x
}
func value(v, max float64) float64 {
	return f(v) * scale(max)
}
func scale(max float64) float64 {
	return (255 * exposure) / f(max)
}

// ptoc converts a point from the complex function to a pixel coordinate.
func ptoc(c complex128) (p image.Point) {
	r, i := real(c), imag(c)

	p.X = int(zoom*(r+real(offset))) + width/2
	p.Y = int(zoom*(i+imag(offset))) + height/2

	return p
}

// render creates an output image file.
func render(img image.Image) (err error) {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	if rotate {
		logrus.Println("[/] Rotating image")
		img = imaging.Rotate270(img)
	}
	logrus.Println("[!] Done:", filename)
	return jpeg.Encode(out, img, &jpeg.Options{jpeg.DefaultQuality})
}

//
func itoc(r, g, b *Histo, incChan chan orbit) {
	for hi := range incChan {
		if hi.it < 10000 {
			continue
		}
		ps := hi.ps
		if hi.it <= 30000 {
			for _, p := range ps {
				r[p.X][p.Y]++
			}
		}
		if hi.it >= 15000 && hi.it <= 50000 {
			for _, p := range ps {
				b[p.X][p.Y]++
			}
		}
		if hi.it >= 30000 && hi.it <= 100000 {
			for _, p := range ps {
				g[p.X][p.Y]++
			}
		}
	}
}
