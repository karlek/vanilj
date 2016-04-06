package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"math/rand"

	"os"
	"runtime"
	"sync"
	"time"

	"github.com/0xC3/progress/barcli"
	"github.com/disintegration/imaging"
	"github.com/dustin/randbo"
	"github.com/karlek/profile"
	"github.com/karlek/verbose"
)

var (
	w          = 2000
	h          = 2000
	offset     = 0.6 + 0i
	zoom       = float64(w) / 2.8
	iterations = 20000
	filename   = "a.png"
)

var random = randbo.NewFrom(rand.NewSource(time.Now().UnixNano()))

func main() {
	verbose.Verbose = true
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := play()
	if err != nil {
		log.Fatalln(err)
	}
}

func initalize() (img *image.NRGBA, visited [][]int) {
	// Output image with black background.
	img = image.NewNRGBA(image.Rect(0, 0, w, h))
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	visited = make([][]int, h)
	for y := range visited {
		visited[y] = make([]int, w)
	}
	return img, visited
}

func play() (err error) {
	defer profile.Start(profile.CPUProfile).Stop()

	verbose.Println("[.]    Initializing.")
	img, visited := initalize()

	verbose.Println("[-]    Calculating visited points.")
	b, err := barcli.New(100)
	if err != nil {
		return err
	}
	for i := 0; i < 100; i++ {
		for j := 0; j < 5000000; j++ {
			rpoint(visited)
		}
		err = b.Inc()
		if err != nil {
			log.Println(err)
		}
		err = b.Print()
		if err != nil {
			log.Println(err)
		}
	}
	verbose.Println("[/]    Creating image.")
	plot(img, visited)
	return save(img)
}

func rpoint(visited [][]int) {
	c := complex(sign()*2*randfloat(), sign()*2*randfloat())
	// Converges.
	points := divergencePrim(c)
	if points == nil {
		return
	}
	for _, z := range points {
		p := ptoi(z)
		// Ignore points outside image.
		if p.X >= w || p.Y >= h || p.X < 0 || p.Y < 0 {
			continue
		}
		visited[p.Y][p.X]++
	}
}

var p = make([]byte, 1)

func sign() float64 {
	random.Read(p)
	r := int(p[0] % 2)
	if r == 1 {
		return -1.0
	}
	return 1.0
}

var p1 = make([]byte, 4)

func randfloat() float64 {
	random.Read(p1)
	b0, b1, b2, b3 := float64(p1[0]), float64(p1[1]), float64(p1[2]), float64(p1[3])
	return (1 / 256.0) * (b0 + (1/256.0)*(b1+(1/256.0)*(b2+(1/256.0)*b3)))
}

func plot(img *image.NRGBA, visited [][]int) {
	max, min := -1, math.MaxInt64
	for _, row := range visited {
		for _, v := range row {
			if v > max {
				max = v
			}
			if v < min && v != 0 {
				min = v
			}
		}
	}
	scale := 0.0
	if max > 255 {
		scale = float64(max) / 255.0
	} else {
		scale = 255.0 / float64(max)
	}
	for y, row := range visited {
		for x, v := range row {
			plotPoint(PlotPoint{p: image.Pt(x, y), v: float64(v), scale: scale}, img)
		}
	}
}

// Point to index.
func ptoi(c complex128) (p image.Point) {
	r, i := real(c), imag(c)

	p.X = int(zoom*(r+real(offset))) + w/2
	p.Y = int(zoom*(i+imag(offset))) + h/2

	return p
}

// add adds two colors together to make a brighter color.
func add(c1 color.Color, c2 color.RGBA) color.RGBA {
	r, g, b, _ := c1.RGBA()
	return color.RGBA{uint8(r) + c2.R, uint8(g) + c2.G, uint8(b) + c2.B, 255}
}

// save creates an output image file.
func save(img *image.NRGBA) (err error) {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	img = imaging.Rotate270(img)

	verbose.Println("[!]    Done:", filename)
	return png.Encode(out, img)
}

func isInBulb(c complex128) bool {
	x, y := real(c), imag(c)
	q := (x-1.0/4.0)*(x-1.0/4.0) + y*y
	return q*(q+(x-1.0/4.0)) < 1.0/4.0*y*y
}

var points = make([]complex128, iterations)

func divergencePrim(c complex128) []complex128 {
	if isInBulb(c) {
		return nil
	}

	z := complex(0, 0)
	var num int
	for i := 0; i < iterations; i++ {
		z = z*z + c
		// Diverges.
		if x, y := real(z), imag(z); 4 < x*x+y*y {
			return points[:num]
		}
		points[num] = z
		num++
	}
	// Converges.
	return nil
}

// isMemberOfSet determines if the complex point z is member of the mandelbrot
// set.
func isMemberOfSet(z complex128) bool {
	// same as 2 > cmplx.Abs(z)
	return real(z)*real(z)+imag(z)*imag(z) < 4
}

// PlotPoint contains a color and a coordinate.
type PlotPoint struct {
	smooth color.Color
	p      image.Point
	scale  float64
	v      float64
}

func plotPoint(p PlotPoint, img *image.NRGBA) {
	var v uint8
	if p.scale > 1 {
		v = uint8(p.v / p.scale)
	} else {
		v = uint8(p.v * p.scale)
	}
	img.Set(p.p.X, p.p.Y, color.RGBA{v, v, v, 255})
}

func worker(rgba *image.NRGBA, tasks <-chan PlotPoint, quit <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case p, ok := <-tasks:
			if !ok {
				return
			}
			plotPoint(p, rgba)
		case <-quit:
			return
		}
	}
}
