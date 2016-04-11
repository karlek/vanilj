package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"log"
	"math"
	"math/rand"

	"os"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/disintegration/imaging"
	"github.com/dustin/randbo"
	"github.com/karlek/profile"
)

const (
	w = 1024
	h = 1024
)

var (
	offset     = 0.4 + 0i
	zoom       = float64(w) / 2.8
	iterations = 200
	filename   = "a.jpeg"
	step       = 0.001
)

type Visit [w][h]int

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := play()
	if err != nil {
		log.Fatalln(err)
	}
}

func initialize() (img *image.NRGBA, v *Visit) {
	// Output image with black background.
	img = image.NewNRGBA(image.Rect(0, 0, w, h))
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	return img, &Visit{}
}

func play() (err error) {
	defer profile.Start(profile.CPUProfile).Stop()

	logrus.Println("[.]    Initializing.")
	img, visited := initialize()

	logrus.Println("[-]    Calculating visited points.")
	// b, err := barcli.New(100)
	// if err != nil {
	// 	return err
	// }
	ordered(visited)
	// for i := 0; i < 100; i++ {
	// 	for j := 0; j < 500000; j++ {
	// 		rpoint(visited)
	// 	}
	// 	err = b.Inc()
	// 	if err != nil {
	// 		logrus.Println(err)
	// 	}
	// 	err = b.Print()
	// 	if err != nil {
	// 		logrus.Println(err)
	// 	}
	// }
	logrus.Println("[/]    Creating image.")
	plot(img, visited)
	return save(img)
}

func ordered(v *Visit) {
	incChan := make(chan image.Point, 20)
	go func(v *Visit, incChan chan image.Point) {
		for p := range incChan {
			if p.X >= w || p.Y >= h || p.X < 0 || p.Y < 0 {
				continue
			}
			v[p.Y][p.X]++
		}
	}(v, incChan)

	// wg := new(sync.WaitGroup)
	// wg.Add(int(4*(1.0/step) + 1))

	for x := -2.0; x <= 2; x += step {
		// go func(x float64, incChan chan image.Point, wg *sync.WaitGroup) {
		for y := -2.0; y <= 2; y += step {
			divergencePrim(complex(x, y), incChan)
		}
		// wg.Done()
		// }(x, incChan, wg)
	}

	// wg.Wait()
	close(incChan)
}

// func point(c complex128, incChan chan image.Point) {
// 	for _, z := range divergencePrim(c) {
// 		p := ptoi(z)
// 		// Ignore points outside image.
// 		if p.X >= w || p.Y >= h || p.X < 0 || p.Y < 0 {
// 			continue
// 		}
// 		incChan <- p
// 	}
// }

func plot(img *image.NRGBA, v *Visit) {
	max, min := -1, math.MaxInt64
	for _, row := range v {
		for _, v := range row {
			if v > max {
				max = v
			}
			if v < min && v != 0 {
				min = v
			}
		}
	}
	logrus.Println("Vistiations:", max, min)
	// s := scale(max, min)
	for y, row := range v {
		for x, v := range row {
			if v == 0 {
				continue
			}
			s := scale(v, max, min)
			plotPoint(PlotPoint{p: image.Pt(x, y), scale: s, v: float64(v)}, img)
		}
	}
}

func scale(v, max, min int) float64 {
	return math.Min((125.0 * float64(v) * (math.Sqrt(float64(v)) / float64(max))), 255.0)
}

// Point to index.
func ptoi(c complex128) (p image.Point) {
	r, i := real(c), imag(c)

	p.X = int(zoom*(r+real(offset))) + w/2
	p.Y = int(zoom*(i+imag(offset))) + h/2

	return p
}

// save creates an output image file.
func save(img *image.NRGBA) (err error) {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	img = imaging.Rotate270(img)

	logrus.Println("[!]    Done:", filename)
	return jpeg.Encode(out, img, &jpeg.Options{Quality: 100})
}

func isInBulb(c complex128) bool {
	x, y := real(c), imag(c)
	q := (x-1.0/4.0)*(x-1.0/4.0) + y*y
	return q*(q+(x-1.0/4.0)) < 1.0/4.0*y*y
}

var points = make([]complex128, iterations)

func divergencePrim(c complex128, incChan chan image.Point) {
	if isInBulb(c) {
		return
	}

	z := complex(0, 0)
	for i := 0; i < iterations; i++ {
		z = z*z + c
		// Diverges.
		if x, y := real(z), imag(z); 4 < x*x+y*y {
			return
		}
		incChan <- ptoi(z)
	}
}

// PlotPoint contains a color and a coordinate.
type PlotPoint struct {
	smooth color.Color
	p      image.Point
	v      float64
	scale  float64
}

func plotPoint(p PlotPoint, img *image.NRGBA) {
	// var r, g, b uint8
	// if p.v < 500 {
	// 	r = uint8(p.scale)
	// } else if p.v < 5000 {
	// 	g = uint8(p.scale)
	// } else if p.v <= 20000 {
	// 	b = uint8(p.scale)
	// }
	// c := color.RGBA{r, g, b, 255}
	c := color.RGBA{uint8(p.scale), uint8(p.scale), uint8(p.scale), 255}
	// c := colorful.Hsv(p.scale, (p.scale / 360), (p.scale / 360))
	img.Set(p.p.X, p.p.Y, c)
}

var random = randbo.NewFrom(rand.NewSource(time.Now().UnixNano()))

// func rpoint(v *Visit) {
// 	c := complex(sign(random)*2*rand.Float64(), sign(random)*2*rand.Float64())
// 	point(c, v)
// }

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
