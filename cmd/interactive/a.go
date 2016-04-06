// simple demonstrates how to draw images using the Draw and DrawRect methods.
// It also gives an example of a basic event loop.
package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"runtime"
	"time"

	"github.com/karlek/profile"
	"github.com/karlek/vanilj/canvas"
	"github.com/karlek/vanilj/fractal"
	"github.com/karlek/vanilj/fractal/mandel"
	"github.com/mewmew/sdl/win"
	"github.com/mewmew/we"
)

func main() {
	err := simple()
	if err != nil {
		log.Fatalln(err)
	}
}

// var width, height int = 2560, 1440

var width, height int = 640, 480

// simple demonstrates how to draw images using the Draw and DrawRect methods.
// It also gives an example of a basic event loop.
func simple() (err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	defer profile.Start(profile.CPUProfile).Stop()

	// Open the window.
	err = win.Open(width, height, win.Resizeable)
	if err != nil {
		return err
	}
	defer win.Close()

	// Load image resources.
	err = loadResources()
	if err != nil {
		return err
	}
	defer freeResources()

	// start and frames will be used to calculate the average FPS of the
	// application.
	start := time.Now()
	frames := 0

	// Render the images onto the window.
	err = render()
	if err != nil {
		return err
	}

	// Update and event loop.
	for {
		// Display window updates on screen.
		err = win.Update()
		if err != nil {
			return err
		}
		frames++

		// Poll events until the event queue is empty.
		for e := win.PollEvent(); e != nil; e = win.PollEvent() {
			fmt.Printf("%T event: %v\n", e, e)
			switch obj := e.(type) {
			case we.KeyPress:
				r := obj.Key.String()
				err = zz(r)
			case we.KeyRepeat:
				r := obj.Key.String()
				err = zz(r)
			case we.KeyRune:
				r := obj.String()
				err = zz(r)
			}
			if err != nil {
				return err
			}
			switch e.(type) {
			case we.Close:
				displayFPS(start, frames)
				// Close the application.
				return nil
			case we.Resize:
				// Rerender the images onto the window after resize events.
				err = render()
				if err != nil {
					return err
				}
			}
		}

		// Cap refresh rate at 30 FPS.
		time.Sleep(time.Second / 30)
	}
}

func zz(s string) (err error) {
	if "a" == s {
		zoom *= 2
	}
	if "b" == s {
		zoom /= 2
	}
	if "c" == s {
		iterations *= 2
	}
	if "d" == s {
		iterations /= 2
	}
	fmt.Printf("asdf: %v\n", s)
	if "[up]" == s {
		centerImag = (centerImag - 1/zoom)
	}
	if "[down]" == s {
		centerImag = (centerImag + 1/zoom)
	}
	if "[left]" == s {
		centerReal = (centerReal - 1/zoom)
	}
	if "[right]" == s {
		centerReal = (centerReal + 1/zoom)
	}
	if "s" == s {
		fmt.Println(zoom, centerReal, centerImag)
	}
	err = loadResources()
	if err != nil {
		return err
	}
	return render()
}

// render renders the background and foreground images onto the window.
func render() (err error) {
	// Draw the entire background image onto the screen starting at the top left
	// point (0, 0).
	dp := image.ZP
	err = win.Draw(dp, fractalImg)
	if err != nil {
		return err
	}

	return nil
}

// Background and foreground images.
var fractalImg *win.Image

var zoom float64 = 1
var iterations float64 = 10.0
var centerReal, centerImag float64 = 0.35, 0.5

// loadResources loads the background and foreground images.
func loadResources() (err error) {
	c := canvas.NewCanvas(width, height)
	for x := 0; x < c.Bounds().Size().X; x++ {
		for y := 0; y < c.Bounds().Size().Y; y++ {
			c.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	f := fractal.Fractal{
		Src:    c.RGBA,
		Iter:   iterations,
		Center: complex(centerReal, centerImag),
		Zoom:   zoom,
	}
	mandel.Smooth(&f)

	// Load background image.
	fractalImg, err = win.ReadImage(f.Src)
	if err != nil {
		return err
	}
	return nil
}

// freeResources frees the memory of the background and foreground images.
func freeResources() {
	fractalImg.Free()
}

// displayFPS calculates and displays the average FPS based on the provided
// frame count.
func displayFPS(start time.Time, frames int) {
	seconds := float64(time.Since(start)) / float64(time.Second)
	fps := float64(frames) / seconds
	fmt.Println()
	fmt.Println("=== [ statistics ] =============================================================")
	fmt.Println()
	fmt.Printf("   Total runtime: %.2f seconds.\n", seconds)
	fmt.Printf("   Frame count:   %d frames\n", frames)
	fmt.Printf("   Average FPS:   %.2f frames/second\n", fps)
}
