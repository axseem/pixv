package pixv

import (
	"fmt"
	"image/color"
	"sort"
	"strings"
)

// BuildSVG constructs the final SVG string from the generated path data.
func BuildSVG(width, height, scale int, pathsByColor map[color.RGBA][]string) string {
	var svg strings.Builder
	svg.WriteString(fmt.Sprintf(
		`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg" shape-rendering="crispEdges">`,
		width, height,
	))

	if scale > 1 {
		svg.WriteString(fmt.Sprintf(`<g transform="scale(%d)">`, scale))
	}

	// Sort colors to ensure deterministic SVG output.
	colors := make([]color.RGBA, 0, len(pathsByColor))
	for c := range pathsByColor {
		colors = append(colors, c)
	}
	sort.Slice(colors, func(i, j int) bool {
		ci, cj := colors[i], colors[j]
		if ci.R != cj.R {
			return ci.R < cj.R
		}
		if ci.G != cj.G {
			return ci.G < cj.G
		}
		if ci.B != cj.B {
			return ci.B < cj.B
		}
		return ci.A < cj.A
	})

	for _, c := range colors {
		paths := pathsByColor[c]
		if len(paths) == 0 {
			continue
		}
		svg.WriteString(fmt.Sprintf(`<path fill="%s"`, toHex(c)))
		if c.A < 255 {
			svg.WriteString(fmt.Sprintf(` fill-opacity="%.2f"`, float64(c.A)/255.0))
		}
		svg.WriteString(` d="`)
		svg.WriteString(strings.Join(paths, " "))
		svg.WriteString(`"/>`)
	}

	if scale > 1 {
		svg.WriteString("</g>")
	}

	svg.WriteString("</svg>")
	return svg.String()
}

// toHex converts a color.RGBA to a hex string like #RRGGBB.
func toHex(c color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}
