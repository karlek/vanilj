package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
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
	"github.com/karlek/progress/barcli"
)

var (
	// Color scaling.
	overexposure = 1.0
	factor       = 10.0
	f            = exp
	// File options.
	filename = "a.jpg"
	load     = false
	rotate   = false
)

const (
	w          = 4096
	h          = 4096
	iterations = 10000000
	bailout    = 4
	step       = 0.001
	tries      = 10000000
	// Camera options.
	offset = 0.4 + 0i
	zoom   = float64(w) / 2.8
)

func itoc(r, g, b *Visit, incChan chan hit) {
	for hi := range incChan {
		if hi.it < 10000 {
			continue
		}
		p := hi.p
		switch {
		case hi.it%3 == 0 && hi.it >= 10000 && hi.it <= 30000:
			r[p.X+p.Y*h]++
		case hi.it%5 == 0 && hi.it >= 30000 && hi.it <= 300000:
			g[p.X+p.Y*h]++
		case hi.it%7 == 0 && hi.it >= 300000:
			b[p.X+p.Y*h]++
		}
	}
}

var fun string

type Visit [w * h]float64

func init() {
	rand.Seed(time.Now().UnixNano())
	flag.BoolVar(&load, "load", false, "use pre-computed values.")
	flag.BoolVar(&rotate, "r", false, "rotate the fractal to an upright position.")
	flag.Float64Var(&overexposure, "o", 3.0, "over exposure")
	flag.Float64Var(&factor, "f", 10.0, "factor")
	flag.StringVar(&fun, "fun", "exp", "color scaling function")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "%s [OPTIONS],,,", os.Args[0])
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
		return save(img)
	}

	logrus.Println("[-] Calculating visited points.")
	incChan := make(chan hit, 100)
	go itoc(r, g, b, incChan)
	// ordered(r, g, b, incChan)
	arbitrary(r, g, b, incChan)
	// close(incChan)
	logrus.Println("[i] Saving r, g, b channels")
	if err := gobVisits(r, g, b); err != nil {
		return err
	}

	logrus.Println("[/] Creating image.")
	plot(img, r, g, b)
	return save(img)
}

func gobVisits(r, g, b *Visit) (err error) {
	file, err := os.Create("r-g-b.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	err = enc.Encode(r)
	if err != nil {
		return err
	}
	err = enc.Encode(g)
	if err != nil {
		return err
	}
	return enc.Encode(b)
}

func loadVisits() (r, g, b *Visit, err error) {
	file, err := os.Open("r-g-b.gob")
	if err != nil {
		return nil, nil, nil, err
	}
	defer file.Close()
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&r); err != nil {
		return nil, nil, nil, err
	}
	if err := dec.Decode(&g); err != nil {
		return nil, nil, nil, err
	}
	if err := dec.Decode(&b); err != nil {
		return nil, nil, nil, err
	}
	return r, g, b, nil
}

func ordered(r, g, b *Visit, incChan chan hit) {
	bar, _ := barcli.New(4 * (1 / step))
	for x := -2.0; x <= 2; x += step {
		for y := -2.0; y <= 2; y += step {
			orbit(complex(x, y), incChan)
		}
		bar.Inc()
		bar.Print()
	}
}

func arbitrary(r, g, b *Visit, incChan chan hit) {
	// bar, _ := barcli.New(100)
	workers := 2000
	wg := new(sync.WaitGroup)
	wg.Add(workers)
	// for i := 0; i < 100; i++ {
	for n := 0; n < workers; n++ {
		go func(incChan chan hit, wg *sync.WaitGroup, iters int) {
			var random = randbo.NewFrom(rand.NewSource(rand.Int63()))
			for j := 0; j < iters; j++ {
				c := complex(sign(random)*2*randfloat(random), sign(random)*2*randfloat(random))
				orbit(c, incChan)
			}
			wg.Done()
		}(incChan, wg, tries/workers)
	}
	// 	bar.Inc()
	// 	bar.Print()
	// }
	wg.Wait()
	close(incChan)
}

func orbit(c complex128, incChan chan hit) {
	points := divergencePrim(c)
	for _, z := range points {
		p := ptoi(z)
		// Ignore points outside image.
		if p.X >= w || p.Y >= h || p.X < 0 || p.Y < 0 {
			continue
		}
		incChan <- hit{p, len(points)}
	}
}

var p = make([]byte, 1)

func sign(random io.Reader) float64 {
	random.Read(p)
	r := int(p[0] % 2)
	if r == 1 {
		return -1.0
	}
	return 1.0
}

var p1 = make([]byte, 4)

func randfloat(random io.Reader) float64 {
	random.Read(p1)
	b0, b1, b2, b3 := float64(p1[0]), float64(p1[1]), float64(p1[2]), float64(p1[3])
	return (1 / 256.0) * (b0 + (1/256.0)*(b1+(1/256.0)*(b2+(1/256.0)*b3)))
}

type hit struct {
	p  image.Point
	it int
}

func max(v *Visit) (max float64) {
	max = -1
	for _, v := range v {
		if v > max {
			max = v
		}
	}
	return max
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func twodim(n float64) (int, int) {
	y := math.Floor(n / 4096)
	x := n - y*4096
	return int(x), int(y)
}

func plot(img *image.RGBA, r, g, b *Visit) {
	rMax := max(r)
	gMax := max(g)
	bMax := max(b)
	logrus.Println("[i] Visitations:", rMax, gMax, bMax)
	logrus.Printf("[i] Function: %s, factor: %.2f, overexposure: %.2f", GetFunctionName(f), factor, overexposure)
	for _, n := range r {
		x, y := twodim(n)
		if r[x+y*h] == 0 &&
			g[x+y*h] == 0 &&
			b[x+y*h] == 0 {
			continue
		}
		c := color.RGBA{
			uint8(value(r[x+y*h], rMax)),
			uint8(value(g[x+y*h], gMax)),
			uint8(value(b[x+y*h], bMax)),
			255}
		img.Set(x, y, c)
	}
}

type pixel struct {
	p image.Point
	c color.RGBA
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

// save creates an output image file.
func save(img image.Image) (err error) {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	if rotate {
		img = imaging.Rotate270(img)
	}
	logrus.Println("[!] Done:", filename)
	return jpeg.Encode(out, img, &jpeg.Options{Quality: 100})
}

// Credits: https://github.com/morcmarc/buddhabrot/blob/master/buddhabrot.go
func isInBulb(c complex128) bool {
	Cr, Ci := real(c), imag(c)
	// Main cardioid
	if !(((Cr-0.25)*(Cr-0.25)+(Ci*Ci))*(((Cr-0.25)*(Cr-0.25)+(Ci*Ci))+(Cr-0.25)) < 0.25*Ci*Ci) {
		// 2nd order period bulb
		if !((Cr+1.0)*(Cr+1.0)+(Ci*Ci) < 0.0625) {
			// smaller bulb left of the period-2 bulb
			if !((((Cr + 1.309) * (Cr + 1.309)) + Ci*Ci) < 0.00345) {
				// smaller bulb bottom of the main cardioid
				if !((((Cr + 0.125) * (Cr + 0.125)) + (Ci-0.744)*(Ci-0.744)) < 0.0088) {
					// smaller bulb top of the main cardioid
					if !((((Cr + 0.125) * (Cr + 0.125)) + (Ci+0.744)*(Ci+0.744)) < 0.0088) {
						return false
					}
				}
			}
		}
	}

	return true
}

var points [iterations]complex128

func divergencePrim(c complex128) []complex128 {
	if isInBulb(c) {
		return nil
	}

	var brent complex128
	z := complex(0, 0)
	var num int
	for i := 0; i < iterations; i++ {
		z = z*z + c
		// Cycle detection.
		if (i-1)&i == 0 && i > 1 {
			brent = z
		} else if z == brent {
			return nil
		}
		// Diverges.
		if x, y := real(z), imag(z); x*x+y*y >= bailout {
			return points[:num]
		}
		points[num] = z
		num++
	}
	// Converges.
	return nil
}
