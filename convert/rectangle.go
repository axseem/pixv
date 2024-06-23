package convert

import (
	"fmt"
	"image"
	"strings"
)

type Rect struct {
	pos    Point
	width  int
	height int
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
	rectanglesMap := make(map[uint32][]Rect)

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
			rectanglesMap[color] = append(rectanglesMap[color], Rect{Point{x, y}, width, height})
		}
	}

	var svg strings.Builder

	svg.WriteString(
		fmt.Sprintf(
			`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg" shape-rendering="crispEdges">`,
			bounds.Max.X*int(scale),
			bounds.Max.Y*int(scale),
		),
	)

	for color, rectangles := range rectanglesMap {
		svg.WriteString(`<path fill="` + ColorToHex(color) + `" d="`)
		anchorPos := Point{}
		for _, rect := range rectangles {
			svg.WriteString(
				fmt.Sprintf(
					`m%d,%dh%dv%dh-%d`,
					(anchorPos.x+rect.pos.x)*int(scale),
					(anchorPos.y+rect.pos.y)*int(scale),
					rect.width*int(scale),
					rect.height*int(scale),
					rect.width*int(scale),
				),
			)
			anchorPos.x = -rect.pos.x
			anchorPos.y = -rect.pos.y - rect.height
		}
		svg.WriteString(`"/>`)
	}

	svg.WriteString("</svg>")

	return svg.String(), nil
}
