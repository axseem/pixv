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

const (
	UP Direction = iota
	RIGHT
	DOWN
	LEFT
)

func (d Direction) rotated(r Direction) Direction {
	return (d + r) % 4
}

func (d *Direction) rotate(r Direction) {
	*d = d.rotated(r)
}

type Point struct {
	x, y int
}

func (p Point) xy() (int, int) {
	return p.x, p.y
}

func (p Point) movedN(d Direction, n int) Point {
	switch d {
	case UP:
		return Point{p.x, p.y - n}
	case RIGHT:
		return Point{p.x + n, p.y}
	case DOWN:
		return Point{p.x, p.y + n}
	case LEFT:
		return Point{p.x - n, p.y}
	}
	panic("unreachable")
}

func (p Point) moved(d Direction) Point {
	return p.movedN(d, 1)
}

func (p *Point) move(d Direction) {
	*p = p.moved(d)
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

func ColorAt(img image.Image, p Point) uint32 {
	return RGBAToColor(img.At(p.xy()).RGBA())
}
