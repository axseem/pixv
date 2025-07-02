package pixv

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"runtime"
	"sync"
)

// vectorizeSquare converts every non-transparent pixel into a 1x1 square path.
func vectorizeSquare(ctx context.Context, data *ImageData) (map[color.RGBA][]string, error) {
	bounds := data.Bounds()
	pixelsByColor := make(map[color.RGBA][]image.Point)
	mu := &sync.Mutex{}

	jobs := make(chan int)
	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU()
	wg.Add(numWorkers)
	for range numWorkers {
		go func() {
			defer wg.Done()
			localPixels := make(map[color.RGBA][]image.Point)
			for y := range jobs {
				for x := range bounds.Dx() {
					c := data.At(x, y)
					if c.A == 0 {
						continue
					}
					localPixels[c] = append(localPixels[c], image.Pt(x, y))
				}
			}

			mu.Lock()
			defer mu.Unlock()
			for c, points := range localPixels {
				pixelsByColor[c] = append(pixelsByColor[c], points...)
			}
		}()
	}

	for y := range bounds.Dy() {
		select {
		case <-ctx.Done():
			close(jobs)
			return nil, ctx.Err()
		case jobs <- y:
		}
	}
	close(jobs)

	wg.Wait()

	pathsByColor := make(map[color.RGBA][]string)
	for c, points := range pixelsByColor {
		paths := make([]string, len(points))
		for i, p := range points {
			paths[i] = fmt.Sprintf("M%d,%dh1v1h-1", p.X, p.Y)
		}
		pathsByColor[c] = paths
	}

	return pathsByColor, nil
}
