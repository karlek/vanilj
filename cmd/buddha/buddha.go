package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"math"
	"math/rand"

	"os"
	"runtime"
	"sync"
	"time"

	"github.com/0xC3/progress/barcli"
	"github.com/Sirupsen/logrus"
	"github.com/disintegration/imaging"
	"github.com/dustin/randbo"
	"github.com/karlek/profile"
)

var (
	w      = 512
	h      = 512
	offset = 0.0 + 0i
	// offset     = 0.6 + 0i
	// zoom       = float64(w) / 2.8
	zoom       = float64(w) / 5
	iterations = 200
	filename   = "a.png"
)

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

	v = new(Visit)
	v.m = make([][]int, h)
	for y := range v.m {
		v.m[y] = make([]int, w)
	}
	return img, v
}

func play() (err error) {
	defer profile.Start(profile.CPUProfile).Stop()

	logrus.Println("[.]    Initializing.")
	img, visited := initialize()

	logrus.Println("[-]    Calculating visited points.")
	b, err := barcli.New(100)
	if err != nil {
		return err
	}
	for i := 0; i < 100; i++ {
		for j := 0; j < 500000; j++ {
			rpoint(visited)
		}
		err = b.Inc()
		if err != nil {
			logrus.Println(err)
		}
		err = b.Print()
		if err != nil {
			logrus.Println(err)
		}
	}
	logrus.Println("[/]    Creating image.")
	plot(img, visited)
	return save(img)
}

type WorkRequest struct {
	Name  string
	Delay time.Duration
}

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				// Receive a work request.
				fmt.Printf("worker%d: Received work request, delaying for %f seconds\n", w.ID, work.Delay.Seconds())

				time.Sleep(work.Delay)
				fmt.Printf("worker%d: Hello, %s!\n", w.ID, work.Name)

			case <-w.QuitChan:
				// We have been asked to stop.
				fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

func rpointprim() {
	name := "asdf"
	delay := time.Second * 1
	// Now, we take the delay, and the person's name, and make a WorkRequest out of them.
	work := WorkRequest{Name: name, Delay: delay}

	// Push the work onto the queue.
	WorkQueue <- work
	fmt.Println("Work request queued")
}

var WorkerQueue chan chan WorkRequest

func StartDispatcher(nworkers int) {
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				fmt.Println("Received work requeust")
				go func() {
					worker := <-WorkerQueue

					fmt.Println("Dispatching work request")
					worker <- work
				}()
			}
		}
	}()
}

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 100)

var random = randbo.NewFrom(rand.NewSource(time.Now().UnixNano()))

func rpoint(v *Visit) {
	c := complex(sign(random)*2*rand.Float64(), sign(random)*2*rand.Float64())
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
		v.Inc(p.X, p.Y)
	}
}

type Visit struct {
	m [][]int
	sync.Mutex
}

func (v *Visit) Inc(x, y int) {
	v.Lock()
	v.m[y][x]++
	v.Unlock()
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

func plot(img *image.NRGBA, visited *Visit) {
	max, min := -1, math.MaxInt64
	for _, row := range visited.m {
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
	for y, row := range visited.m {
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

	logrus.Println("[!]    Done:", filename)
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
