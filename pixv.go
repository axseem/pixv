package pixv

import (
	"context"
	"fmt"
	"image"
	"image/color"
)

// Strategy defines the vectorization algorithm to use.
type Strategy string

const (
	StrategyPath      Strategy = "path"
	StrategyRectangle Strategy = "rectangle"
	StrategySquare    Strategy = "square"
)

// IsValid checks if the strategy is a recognized value.
func (s Strategy) IsValid() bool {
	switch s {
	case StrategyPath, StrategyRectangle, StrategySquare:
		return true
	default:
		return false
	}
}

// Options configures the vectorization process.
type Options struct {
	Method Strategy
	Scale  int
}

// Vectorize converts a raster image into an SVG string using the specified options.
func Vectorize(ctx context.Context, img image.Image, opts Options) (string, error) {
	if opts.Scale <= 0 {
		return "", fmt.Errorf("scale must be positive")
	}

	imgData := NewImageData(img)
	width := imgData.Bounds().Dx() * opts.Scale
	height := imgData.Bounds().Dy() * opts.Scale

	var pathsByColor map[color.RGBA][]string
	var err error

	switch opts.Method {
	case StrategyPath:
		pathsByColor, err = vectorizePath(imgData)
	case StrategyRectangle:
		pathsByColor, err = vectorizeRectangle(ctx, imgData)
	case StrategySquare:
		pathsByColor, err = vectorizeSquare(ctx, imgData)
	default:
		return "", fmt.Errorf("unknown vectorization method: %s", opts.Method)
	}

	if err != nil {
		return "", err
	}

	return BuildSVG(width, height, opts.Scale, pathsByColor), nil
}
