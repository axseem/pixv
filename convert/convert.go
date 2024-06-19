package convert

import (
	"image"
	"log"
	"os"
	"strconv"

	_ "image/jpeg"
	_ "image/png"
)

type Direction uint8

// for global direction FORWARD can be used as UP and BACKWARD as DOWN
const (
	FORWARD Direction = iota
	RIGHT
	BACKWARD
	LEFT
)

func (d Direction) rotated(r Direction) Direction {
	return (d + r) % 4
}

func (d Direction) String() string {
	switch d {
	case FORWARD:
		return "Forward"
	case RIGHT:
		return "Right"
	case BACKWARD:
		return "Backward"
	case LEFT:
		return "Left"
	}
	panic("unreachable")
}

type Point struct {
	x, y int
}

func (p Point) xy() (int, int) {
	return p.x, p.y
}

func (p Point) moved(d Direction) Point {
	switch d {
	case FORWARD:
		return Point{p.x, p.y - 1}
	case RIGHT:
		return Point{p.x + 1, p.y}
	case BACKWARD:
		return Point{p.x, p.y + 1}
	case LEFT:
		return Point{p.x - 1, p.y}
	}
	panic("unreachable")
}

func RGBAToColor(r, g, b, a uint32) uint32 {
	return (r/257)<<24 + (g/257)<<16 + (b/257)<<8 + (a / 257)
}

func ColorToHex(color uint32) string {
	hex := strconv.FormatInt(int64(color), 16)
	for len(hex) < 8 {
		hex = "0" + hex
	}
	if hex[6:] == "ff" {
		hex = hex[:6]
	}
	return "#" + hex
}

func OpenImage(file string) (image.Image, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}
