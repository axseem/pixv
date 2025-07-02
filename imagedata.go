package pixv

import (
	"image"
	"image/color"
	"sync/atomic"
)

// ImageData holds a pre-processed, optimized representation of an image.
type ImageData struct {
	pixels  []color.RGBA
	bounds  image.Rectangle
	visited []uint32 // Use uint32 for atomic operations (0=false, 1=true)
}

// NewImageData creates a new ImageData object from a standard image.Image.
func NewImageData(img image.Image) *ImageData {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	pixels := make([]color.RGBA, width*height)

	for y := range height {
		for x := range width {
			c := color.RGBAModel.Convert(img.At(x+bounds.Min.X, y+bounds.Min.Y)).(color.RGBA)
			pixels[y*width+x] = c
		}
	}

	return &ImageData{
		pixels:  pixels,
		bounds:  image.Rect(0, 0, width, height),
		visited: make([]uint32, width*height),
	}
}

// At returns the color of the pixel at (x, y).
func (d *ImageData) At(x, y int) color.RGBA {
	if !image.Pt(x, y).In(d.bounds) {
		return color.RGBA{}
	}
	return d.pixels[y*d.bounds.Dx()+x]
}

// Bounds returns the image dimensions.
func (d *ImageData) Bounds() image.Rectangle {
	return d.bounds
}

// IsVisited checks if the pixel at (x, y) has been visited. This is thread-safe.
func (d *ImageData) IsVisited(x, y int) bool {
	return atomic.LoadUint32(&d.visited[y*d.bounds.Dx()+x]) == 1
}

// Visit marks the pixel at (x, y) as visited. It returns true if the pixel was successfully marked. This is thread-safe.
func (d *ImageData) Visit(x, y int) bool {
	return atomic.CompareAndSwapUint32(&d.visited[y*d.bounds.Dx()+x], 0, 1)
}

// visitRect marks a rectangular area as visited. This is thread-safe.
func (d *ImageData) visitRect(r image.Rectangle) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			atomic.StoreUint32(&d.visited[y*d.bounds.Dx()+x], 1)
		}
	}
}
