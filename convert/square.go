package convert

import (
	"fmt"
	"image"
	"strings"
)

func Square(img image.Image, scale uint) (string, error) {
	if scale == 0 {
		return "", fmt.Errorf("scale can't be 0")
	}

	bounds := img.Bounds()
	pixelsMap := make(map[uint32][]Point)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			color := RGBAToColor(img.At(x, y).RGBA())
			pixelsMap[color] = append(pixelsMap[color], Point{x, y})
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

	for color, pixels := range pixelsMap {
		svg.WriteString(`<path fill="` + ColorToHex(color) + `" d="`)
		anchorPos := Point{}
		for _, pixel := range pixels {
			svg.WriteString(
				fmt.Sprintf(
					`m%d,%dh%dv%dh-%d`,
					(anchorPos.x+pixel.x)*int(scale),
					(anchorPos.y+pixel.y)*int(scale),
					scale,
					scale,
					scale,
				),
			)
			anchorPos.x = -pixel.x
			anchorPos.y = -pixel.y - 1
		}
		svg.WriteString(`"/>`)

	}

	svg.WriteString("</svg>")

	return svg.String(), nil
}
