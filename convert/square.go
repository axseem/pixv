package convert

import (
	"fmt"
	"image"
	"sort"
	"strings"
)

type Pixel struct {
	pos   Point
	color uint32
}

func Square(img image.Image, scale uint) (string, error) {
	if scale == 0 {
		return "", fmt.Errorf("scale can't be 0")
	}

	bounds := img.Bounds()
	pixels := []Pixel{}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			color := RGBAToColor(img.At(x, y).RGBA())
			pixels = append(pixels, Pixel{Point{x, y}, color})
		}
	}

	sort.Slice(pixels, func(i, j int) bool {
		return pixels[i].color < pixels[j].color
	})

	var svg strings.Builder

	svg.WriteString(
		fmt.Sprintf(
			`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg" shape-rendering="crispEdges">`,
			bounds.Max.X*int(scale),
			bounds.Max.Y*int(scale),
		),
	)

	prevColor := pixels[0].color
	svg.WriteString(`<g fill="` + ColorToHex(prevColor) + `">`)

	for _, pixel := range pixels {
		if pixel.color != prevColor {
			prevColor = pixel.color
			svg.WriteString("</g>")
			svg.WriteString(`<g fill="` + ColorToHex(pixel.color) + `">`)
		}
		svg.WriteString(
			fmt.Sprintf(
				`<path d="m%d,%dh%dv%dh-%d"/>`,
				pixel.pos.x*int(scale),
				pixel.pos.y*int(scale),
				scale,
				scale,
				scale,
			),
		)
	}

	svg.WriteString("</g>")
	svg.WriteString("</svg>")

	return svg.String(), nil
}
