package convert

import (
	"fmt"
	"image"
	"strconv"
	"strings"
)

type Edge struct {
	pos       Point
	direction Direction
	length    int
}

func (e Edge) endPos() Point {
	switch e.direction {
	case FORWARD:
		return Point{e.pos.x, e.pos.y - e.length}
	case RIGHT:
		return Point{e.pos.x + e.length, e.pos.y}
	case BACKWARD:
		return Point{e.pos.x, e.pos.y + e.length}
	case LEFT:
		return Point{e.pos.x - e.length, e.pos.y}
	}
	panic("unreachable")
}

type Shape struct {
	edges []Edge
	color uint32
}

func findShape(entry Point, img image.Image, visited *[][]bool) Shape {
	ant := Ant{entry, RIGHT}
	edges := []Edge{{ant.pos, RIGHT, 1}}
	color := img.At(ant.pos.xy())

	for edges[len(edges)-1].endPos() != edges[0].pos {
		if img.At(ant.moved(LEFT).xy()) == img.At(ant.pos.xy()) {
			edges[len(edges)-1].length--
			ant.pos = ant.moved(LEFT)
			ant.direction = ant.direction.rotated(BACKWARD)
		}

		if img.At(ant.moved(FORWARD).xy()) == img.At(ant.pos.xy()) {
			edges[len(edges)-1].length++
			ant.pos = ant.moved(FORWARD)
		} else {
			ant.direction = ant.direction.rotated(RIGHT)
			edges = append(edges, Edge{edges[len(edges)-1].endPos(), ant.direction, 1})
		}
	}

	queue := []Point{entry}

	for len(queue) > 0 {
		pos := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		(*visited)[pos.y][pos.x] = true

		for i := range 4 {
			neighbor := pos.moved(Direction(i))
			if img.At(neighbor.xy()) == color && !(*visited)[neighbor.y][neighbor.x] {
				queue = append(queue, neighbor)
			}
		}
	}

	return Shape{edges, RGBAToColor(color.RGBA())}
}

type Ant struct {
	pos       Point
	direction Direction
}

func (a Ant) moved(d Direction) Point {
	return a.pos.moved(a.direction.rotated(d))
}

func Path(img image.Image, scale uint) (string, error) {
	bounds := img.Bounds()

	shapes := []Shape{}

	visited := make([][]bool, bounds.Max.Y)
	for i := range visited {
		visited[i] = make([]bool, bounds.Max.X)
	}

	for y, row := range visited {
		for x, v := range row {
			if !v {
				shapes = append(shapes, findShape(Point{x, y}, img, &visited))
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

	for _, shape := range shapes {
		svg.WriteString(`<path fill="` + ColorToHex(shape.color) + `" d="`)
		svg.WriteString("m" + strconv.Itoa(shape.edges[0].pos.x*int(scale)) + "," + strconv.Itoa(shape.edges[0].pos.y*int(scale)))
		for _, edge := range shape.edges {
			switch edge.direction {
			case FORWARD:
				svg.WriteString(`v-`)
			case RIGHT:
				svg.WriteString(`h`)
			case BACKWARD:
				svg.WriteString(`v`)
			case LEFT:
				svg.WriteString(`h-`)
			}
			svg.WriteString(strconv.Itoa(edge.length * int(scale)))
		}
		svg.WriteString(`"/>`)
	}
	svg.WriteString("</svg>")

	return svg.String(), nil
}
