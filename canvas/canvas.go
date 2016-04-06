// Package canvas implements utility functions for RGBA images.
package canvas

import (
	"image"
	"image/png"
	"os"

	"github.com/mewkiz/pkg/errutil"
)

// Canvas is drawable RGBA rectangle.
type Canvas struct {
	*image.RGBA
}

// Save saves the canvas with filename to disk.
func (c *Canvas) Save(filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return errutil.Err(err)
	}
	defer f.Close()
	return png.Encode(f, c)
}

// NewCanvas initalizes a new canvas with the dimensions supplied on function
// call.
func NewCanvas(width, height int) Canvas {
	return Canvas{image.NewRGBA(image.Rect(0, 0, width, height))}
}
