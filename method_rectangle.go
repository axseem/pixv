package pixv

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"runtime"
	"sync"
)

// vectorizeRectangle finds the largest possible rectangular chunks of color.
func vectorizeRectangle(ctx context.Context, data *ImageData) (map[color.RGBA][]string, error) {
	bounds := data.Bounds()
	rectsByColor := make(map[color.RGBA][]image.Rectangle)
	mu := &sync.Mutex{}

	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	rowsPerWorker := (bounds.Dy() + numWorkers - 1) / numWorkers

	for i := range numWorkers {
		yStart := i * rowsPerWorker
		yEnd := min(bounds.Dy(), yStart+rowsPerWorker)
		if yStart >= yEnd {
			continue
		}

		wg.Add(1)
		go func(yStart, yEnd int) {
			defer wg.Done()
			localRects := make(map[color.RGBA][]image.Rectangle)

			for y := yStart; y < yEnd; y++ {
				for x := range bounds.Dx() {
					if ctx.Err() != nil {
						return
					}

					if !data.Visit(x, y) {
						continue
					}

					c := data.At(x, y)
					if c.A == 0 {
						continue
					}

					rect := findRect(data, x, y, c)
					data.visitRect(rect)
					localRects[c] = append(localRects[c], rect)
				}
			}

			mu.Lock()
			defer mu.Unlock()
			for c, rects := range localRects {
				rectsByColor[c] = append(rectsByColor[c], rects...)
			}
		}(yStart, yEnd)
	}

	wg.Wait()
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	pathsByColor := make(map[color.RGBA][]string)
	for c, rects := range rectsByColor {
		paths := make([]string, len(rects))
		for i, r := range rects {
			paths[i] = fmt.Sprintf("M%d,%dh%dv%dh-%d",
				r.Min.X, r.Min.Y,
				r.Dx(), r.Dy(), r.Dx())
		}
		pathsByColor[c] = paths
	}

	return pathsByColor, nil
}

// findRect finds the largest possible rectangle of a given color starting at (x, y).
func findRect(data *ImageData, x, y int, c color.RGBA) image.Rectangle {
	bounds := data.Bounds()
	startX := x

	// Find width
	for x+1 < bounds.Dx() && !data.IsVisited(x+1, y) && data.At(x+1, y) == c {
		x++
	}
	width := x - startX + 1

	// Find height
	yEnd := y
	for yEnd+1 < bounds.Dy() {
		canExpand := true
		for i := range width {
			if data.IsVisited(startX+i, yEnd+1) || data.At(startX+i, yEnd+1) != c {
				canExpand = false
				break
			}
		}
		if !canExpand {
			break
		}
		yEnd++
	}
	height := yEnd - y + 1

	return image.Rect(startX, y, startX+width, y+height)
}
