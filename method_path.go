package pixv

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

type direction uint8

const (
	up direction = iota
	right
	down
	left
)

// rotated calculates a new direction by rotating the current one.
func (d direction) rotated(r direction) direction {
	return (d + r) % 4
}

// rotate changes the direction d by rotating it.
func (d *direction) rotate(r direction) {
	*d = d.rotated(r)
}

type point struct{ x, y int }

// movedN calculates a new point by moving n steps in a direction.
func (p point) movedN(d direction, n int) point {
	switch d {
	case up:
		return point{p.x, p.y - n}
	case right:
		return point{p.x + n, p.y}
	case down:
		return point{p.x, p.y + n}
	case left:
		return point{p.x - n, p.y}
	}
	panic("unreachable")
}

// moved calculates a new point by moving 1 step in a direction.
func (p point) moved(d direction) point {
	return p.movedN(d, 1)
}

// move moves the point 1 step in a direction.
func (p *point) move(d direction) {
	*p = p.moved(d)
}

type line struct {
	start     point
	len       int
	direction direction
}

// end returns the end point of the line.
func (l line) end() point {
	return l.start.movedN(l.direction, l.len)
}

// vectorizePath converts shapes of same-colored pixels into a single path.
func vectorizePath(data *ImageData) (map[color.RGBA][]string, error) {
	bounds := data.Bounds()
	linesByColor := make(map[color.RGBA][]line)

	for y := range bounds.Dy() {
		for x := range bounds.Dx() {
			if data.IsVisited(x, y) {
				continue
			}
			c := data.At(x, y)
			if c.A == 0 {
				continue
			}

			shapeLines := findShape(data, point{x, y})
			linesByColor[c] = append(linesByColor[c], shapeLines...)
		}
	}

	return buildPathSVG(linesByColor), nil
}

// findShape traces outer and inner edges of a shape, building a path.
func findShape(data *ImageData, entry point) []line {
	bounds := data.Bounds()
	// vertices is sized to handle vertex coordinates, which are one larger than pixel coordinates.
	vertices := make([][]bool, bounds.Dy()+1)
	for i := range vertices {
		vertices[i] = make([]bool, bounds.Dx()+1)
	}

	shapeColor := data.At(entry.x, entry.y)
	shapePaths := findPath(data, entry, right, &vertices)

	// Flood fill to find inner holes and mark the whole shape as visited.
	current := []point{entry}
	for len(current) > 0 {
		pixel := current[len(current)-1]
		current = current[:len(current)-1]
		if !data.Visit(pixel.x, pixel.y) {
			continue
		}

		for i := range 4 {
			direction := direction(i)
			neighbor := pixel.moved(direction)

			if !neighbor.in(bounds) {
				continue
			}

			if data.At(neighbor.x, neighbor.y) == shapeColor {
				if !data.IsVisited(neighbor.x, neighbor.y) {
					current = append(current, neighbor)
				}
				continue
			}

			// Detecting inner holes.
			pDown := pixel.moved(down)
			npDown := neighbor.moved(down)
			if direction == right &&
				pDown.in(bounds) && npDown.in(bounds) &&
				!vertices[npDown.y][npDown.x] &&
				data.At(pDown.x, pDown.y) == shapeColor &&
				data.At(npDown.x, npDown.y) == shapeColor {

				innerPath := findPath(data, npDown, right, &vertices)
				shapePaths = append(shapePaths, innerPath...)
			}
		}
	}

	return shapePaths
}

// findPath traces a continuous shape outline into a path.
func findPath(data *ImageData, ant point, antDirection direction, vertices *[][]bool) []line {
	path := []line{{ant, 1, antDirection}}
	startPoint := ant

	for {
		if path[len(path)-1].end() == startPoint && data.At(ant.x, ant.y) != data.At(ant.moved(antDirection.rotated(left)).x, ant.moved(antDirection.rotated(left)).y) {
			break
		}

		leftHandPixel := ant.moved(antDirection.rotated(left))
		frontPixel := ant.moved(antDirection)

		leftHandEqual := data.At(ant.x, ant.y) == data.At(leftHandPixel.x, leftHandPixel.y)
		frontEqual := data.At(ant.x, ant.y) == data.At(frontPixel.x, frontPixel.y)
		last := &path[len(path)-1]

		if leftHandEqual {
			last.len--
			antDirection.rotate(left)
			ant.move(antDirection)
			path = append(path, line{last.end(), 1, antDirection})
			(*vertices)[last.start.y][last.start.x] = true
			continue
		}

		if frontEqual {
			last.len++
			ant.move(antDirection)
		} else {
			antDirection.rotate(right)
			path = append(path, line{last.end(), 1, antDirection})
			(*vertices)[last.start.y][last.start.x] = true
		}
	}
	return path
}

// in checks if the point is within the rectangle's bounds.
func (p point) in(r image.Rectangle) bool {
	return p.x >= r.Min.X && p.x < r.Max.X && p.y >= r.Min.Y && p.y < r.Max.Y
}

func buildPathSVG(linesByColor map[color.RGBA][]line) map[color.RGBA][]string {
	pathsByColor := make(map[color.RGBA][]string)
	for c, lines := range linesByColor {
		if len(lines) == 0 {
			continue
		}

		var b strings.Builder
		prevEnd := point{-1, -1} // Guaranteed to be different from any valid point.

		for _, l := range lines {
			if l.start != prevEnd {
				b.WriteString(fmt.Sprintf("M%d,%d", l.start.x, l.start.y))
			}

			switch l.direction {
			case up:
				b.WriteString(fmt.Sprintf("v-%d", l.len))
			case right:
				b.WriteString(fmt.Sprintf("h%d", l.len))
			case down:
				b.WriteString(fmt.Sprintf("v%d", l.len))
			case left:
				b.WriteString(fmt.Sprintf("h-%d", l.len))
			}
			prevEnd = l.end()
		}
		pathsByColor[c] = []string{b.String()}
	}
	return pathsByColor
}
