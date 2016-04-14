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
	overexposure float64
	factor       float64
	f            func(float64) float64
	fun          string
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
	w          = 4096
	h          = 4096
	iterations = 1000000
	// Camera options.
	zoom = float64(w) / 2.8
)

func init() {
	rand.Seed(time.Now().UnixNano())
	flag.BoolVar(&load, "load", false, "use pre-computed values.")
	flag.BoolVar(&save, "save", false, "save orbits.")
	flag.BoolVar(&rotate, "rotate", false, "rotate the fractal to an upright position.")
	flag.Float64Var(&overexposure, "oe", 1.0, "over exposure")
	flag.Float64Var(&factor, "f", 1.0, "factor")
	flag.StringVar(&fun, "fun", "exp", "color scaling function")
	flag.StringVar(&filename, "o", "a.jpg", "output filename")
	flag.IntVar(&tries, "t", 12000000, "number of orbits attempts")
	flag.Float64Var(&bailout, "b", 4, "bailout value")
	flag.Float64Var(&offsetReal, "r", 0.4, "offsetReal")
	flag.Float64Var(&offsetImag, "i", 0, "offsetImag")

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

	if err := play(); err != nil {
		logrus.Fatalln(err)
	}
}

func initialize() (img *image.RGBA, r, g, b *Visit) {
	// Output image with black background.
	return image.NewRGBA(image.Rect(0, 0, w, h)), &Visit{}, &Visit{}, &Visit{}
}

func play() (err error) {
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

	workers := 16
	wg := new(sync.WaitGroup)
	wg.Add(workers)
	for n := 0; n < workers; n++ {
		incChan := make(chan hit, tries/workers/1000)
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

func arbitrary(n int, incChan chan hit, wg *sync.WaitGroup) {
	var random = randbo.NewFrom(rand.NewSource(rand.Int63()))
	var send [iterations]image.Point
	for i := 0; i < n; i++ {
		c := complex(sign(random)*2*randfloat(random), sign(random)*2*randfloat(random))
		orbit(c, incChan, &send)
	}
	wg.Done()
	close(incChan)
}

func orbit(c complex128, incChan chan hit, send *[iterations]image.Point) {
	points := divergencePrim(c)
	num := 0
	for _, z := range points {
		p := ptoi(z)
		// Ignore points outside image.
		if p.X >= w || p.Y >= h || p.X < 0 || p.Y < 0 {
			continue
		}
		send[num] = p
		num++
	}
	incChan <- hit{send[:num], len(points)}
}

type hit struct {
	ps []image.Point
	it int
}

func max(v *Visit) (max float64) {
	max = -1
	for _, row := range v {
		for _, v := range row {
			if v > max {
				max = v
			}
		}
	}
	return max
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func plot(img *image.RGBA, r, g, b *Visit) {
	rMax := max(r)
	gMax := max(g)
	bMax := max(b)
	logrus.Println("[i] Visitations:", rMax, gMax, bMax)
	logrus.Printf("[i] Function: %s, factor: %.2f, overexposure: %.2f", GetFunctionName(f), factor, overexposure)
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
	return (255 * overexposure) / f(max)
}

// Point to index.
func ptoi(c complex128) (p image.Point) {
	r, i := real(c), imag(c)

	p.X = int(zoom*(r+real(offset))) + w/2
	p.Y = int(zoom*(i+imag(offset))) + h/2

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
	return jpeg.Encode(out, img, &jpeg.Options{Quality: 100})
}

func itoc(r, g, b *Visit, incChan chan hit) {
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

type Visit [w][h]float64
