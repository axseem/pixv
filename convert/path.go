package convert

import (
	"fmt"
	"image"
	"strconv"
	"strings"
)

type line struct {
	start     Point
	len       int
	direction Direction
}

func (l line) end() Point {
	return l.start.movedN(l.direction, l.len)
}

func findPath(img image.Image, ant Point, antDirection Direction, vertices *[][]bool) []line {
	path := []line{{ant, 1, antDirection}}

	// color check left-hand pixel to avoid one pixel gap bug.
	// TODO find better solution
	for path[len(path)-1].end() != path[0].start || ColorAt(img, ant) == ColorAt(img, ant.moved(antDirection.rotated(LEFT))) {
		leftHandEqual := ColorAt(img, ant) == ColorAt(img, ant.moved(antDirection.rotated(LEFT)))
		frontEqual := ColorAt(img, ant) == ColorAt(img, ant.moved(antDirection))
		last := &path[len(path)-1]

		if leftHandEqual {
			last.len--
			antDirection.rotate(LEFT)
			ant.move(antDirection)
			path = append(path, line{last.end(), 1, antDirection})
			(*vertices)[last.start.y][last.start.x] = true
			continue
		}

		if frontEqual {
			last.len++
			ant.move(antDirection)

		} else {
			antDirection.rotate(RIGHT)
			path = append(path, line{last.end(), 1, antDirection})
			(*vertices)[last.start.y][last.start.x] = true
		}
	}

	// no need to draw the last line since inner area will stay the same
	return path
}

func findShape(img image.Image, entry Point, visited *[][]bool) []line {
	vertices := make([][]bool, img.Bounds().Max.Y+1)
	for i := range vertices {
		vertices[i] = make([]bool, img.Bounds().Max.X+1)
	}

	shapePaths := findPath(img, entry, RIGHT, &vertices)
	shapeColor := RGBAToColor(img.At(entry.xy()).RGBA())
	stack := []Point{entry}
	for len(stack) != 0 {
		pixel := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		(*visited)[pixel.y][pixel.x] = true

		for i := range 4 {
			direction := Direction(i)
			neighbor := pixel.moved(direction)

			if ColorAt(img, neighbor) == shapeColor && !(*visited)[neighbor.y][neighbor.x] {
				stack = append(stack, neighbor)
				continue
			}

			if direction == RIGHT &&
				!vertices[neighbor.moved(DOWN).y][neighbor.moved(DOWN).x] &&
				ColorAt(img, neighbor) != shapeColor &&
				ColorAt(img, pixel.moved(DOWN)) == shapeColor &&
				ColorAt(img, neighbor.moved(DOWN)) == shapeColor {

				innerPath := findPath(img, neighbor.moved(DOWN), RIGHT, &vertices)
				shapePaths = append(shapePaths, innerPath...)
			}
		}
	}

	return shapePaths
}

func Path(img image.Image, scale uint) (string, error) {
	if scale == 0 {
		return "", fmt.Errorf("scale can't be 0")
	}

	bounds := img.Bounds()
	pathsMap := make(map[uint32][]line)

	visited := make([][]bool, bounds.Max.Y)
	for i := range visited {
		visited[i] = make([]bool, bounds.Max.X)
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if !visited[y][x] {
				color := RGBAToColor(img.At(x, y).RGBA())
				shape := findShape(img, Point{x, y}, &visited)
				pathsMap[color] = append(pathsMap[color], shape...)
			}
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

	for color, lines := range pathsMap {
		svg.WriteString(`<path fill="` + ColorToHex(color) + `" d="`)
		prevLine := lines[0]
		svg.WriteString(fmt.Sprintf(`m%d,%d`, prevLine.start.x*int(scale), prevLine.start.y*int(scale)))
		prevLine.start = prevLine.end()
		prevLine.direction.rotate(DOWN)

		for _, line := range lines {
			if prevLine.end() != line.start {
				svg.WriteString(
					fmt.Sprintf(
						`m%d,%d`,
						(-prevLine.end().x+line.start.x)*int(scale),
						(-prevLine.end().y+line.start.y)*int(scale),
					),
				)
			}
			prevLine = line
			switch line.direction {
			case UP:
				svg.WriteString("v-")
			case RIGHT:
				svg.WriteString("h")
			case DOWN:
				svg.WriteString("v")
			case LEFT:
				svg.WriteString("h-")
			}
			svg.WriteString(strconv.Itoa(line.len * int(scale)))
		}
		svg.WriteString(`"/>`)
	}

	svg.WriteString("</svg>")

	return svg.String(), nil
}
