package convert

import (
	"fmt"
	"image"
	"sort"
	"strings"
)

type Rect struct {
	pos    Point
	width  int
	height int
	color  uint32
}

func findChunk(img image.Image, x, y int, visited *[][]bool) (int, int) {
	color := RGBAToColor(img.At(x, y).RGBA())
	width := 0
	for i := 1; x+i <= img.Bounds().Max.X; i++ {
		nextColor := RGBAToColor(img.At(x+i, y).RGBA())
		// TODO make occupation check optional.
		// Disabling occupation check will lead to overlapping
		// of some chunks but will produce smaller file size.
		(*visited)[y][x+i-1] = true
		if color != nextColor || (*visited)[y][x+i] {
			width += i
			break
		}
	}

	height := 0

line:
	for i := 1; y+i <= img.Bounds().Max.Y; i++ {
		for j := 0; j < width; j++ {
			nextColor := RGBAToColor(img.At(x+j, y+i).RGBA())

			if color != nextColor {
				height += i
				break line
			}
		}
		for j := 0; j < width; j++ {
			(*visited)[y+i][x+j] = true
		}
	}

	return width, height
}

func Rectangle(img image.Image, scale uint) (string, error) {
	if scale == 0 {
		return "", fmt.Errorf("scale can't be 0")
	}

	bounds := img.Bounds()
	rectangles := []Rect{}

	visited := make([][]bool, bounds.Max.Y)
	for i := range visited {
		visited[i] = make([]bool, bounds.Max.X)
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if visited[y][x] {
				continue
			}

			color := RGBAToColor(img.At(x, y).RGBA())
			width, height := findChunk(img, x, y, &visited)
			rectangles = append(rectangles, Rect{Point{x, y}, width, height, color})
		}
	}

	sort.Slice(rectangles, func(i, j int) bool {
		a := rectangles[i]
		b := rectangles[j]

		if a.color != b.color {
			return a.color < b.color
		}

		if a.height != b.height {
			return a.height < b.height
		}

		return a.width < b.width
	})

	var svg strings.Builder

	svg.WriteString(
		fmt.Sprintf(
			`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg" shape-rendering="crispEdges">`,
			bounds.Max.X*int(scale),
			bounds.Max.Y*int(scale),
		),
	)

	prevColor := rectangles[0].color
	svg.WriteString(`<g fill="` + ColorToHex(prevColor) + `">`)

	for _, rect := range rectangles {
		if rect.color != prevColor {
			prevColor = rect.color
			svg.WriteString("</g>")
			svg.WriteString(`<g fill="` + ColorToHex(rect.color) + `">`)
		}
		svg.WriteString(
			fmt.Sprintf(
				`<path d="m%d,%dh%dv%dh-%d"/>`,
				rect.pos.x*int(scale),
				rect.pos.y*int(scale),
				rect.width*int(scale),
				rect.height*int(scale),
				rect.width*int(scale),
			),
		)
	}

	svg.WriteString("</g>")
	svg.WriteString("</svg>")

	return svg.String(), nil
}
