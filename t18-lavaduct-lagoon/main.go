package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"log"
	"strconv"
	"strings"
)

const (
	Empty rune = '.'
	Dug   rune = '#'
)

type Coord struct {
	x, y int
}

type TaskContext struct {
	directions common.Directions
}

type Field struct {
	tiles          [][]rune
	startX, startY int
}

func (c TaskContext) DirectionFromChar(dir string) common.DirectionDesc {
	switch dir {
	case "R":
		return c.directions.Right
	case "L":
		return c.directions.Left
	case "U":
		return c.directions.Up
	case "D":
		return c.directions.Down
	}
	log.Fatalf("Unreacheable")
	return c.directions.Left
}

type DigStep struct {
	direction common.DirectionDesc
	size      int
}

func ParseDigStep(ctx TaskContext, s string) DigStep {
	// L 5 (#7e2d02)
	components := strings.Fields(s)
	direction := ctx.DirectionFromChar(components[0])
	size, err := strconv.Atoi(components[1])
	if err != nil {
		log.Fatalf("Failed to parse int")
	}
	return DigStep{
		direction: direction,
		size:      size,
	}
}

func ParseDigStep2(ctx TaskContext, s string) DigStep {
	// L 5 (#7e2d02)
	components := strings.Fields(s)
	rawRgb := strings.Trim(components[2], "()")
	var size, rawDirection int
	_, err := fmt.Sscanf(rawRgb, "#%05x%01x", &size, &rawDirection)
	var direction common.DirectionDesc
	// 0 means R, 1 means D, 2 means L, and 3 means U.
	switch rawDirection {
	case 0:
		direction = ctx.directions.Right
	case 1:
		direction = ctx.directions.Down
	case 2:
		direction = ctx.directions.Left
	case 3:
		direction = ctx.directions.Up
	default:
		log.Fatalf("Unreacheable")
	}
	if err != nil {
		log.Fatalf("Failed to parse color")
	}
	return DigStep{
		direction: direction,
		size:      size,
	}
}

func NewField(digSteps []DigStep) Field {
	var minX, maxX, minY, maxY int
	var x, y int
	for _, step := range digSteps {
		x += step.direction.DeltaX * step.size
		y += step.direction.DeltaY * step.size

		minY = min(minY, y)
		maxY = max(maxY, y)
		minX = min(minX, x)
		maxX = max(maxX, x)
	}
	lenX := maxX - minX + 1
	lenY := maxY - minY + 1
	startX := -minX
	startY := -minY

	tiles := make([][]rune, 0, lenY)
	for i := 0; i < lenY; i++ {
		tiles = append(tiles, make([]rune, 0, lenX))
		for j := 0; j < lenX; j++ {
			tiles[len(tiles)-1] = append(tiles[len(tiles)-1], Empty)
		}
	}
	return Field{
		tiles:  tiles,
		startX: startX,
		startY: startY,
	}
}

func (f *Field) DigSteps(digSteps []DigStep) {
	x, y := f.startX, f.startY
	f.tiles[f.startY][f.startX] = Dug
	for _, step := range digSteps {
		for i := 0; i < step.size; i++ {
			x += step.direction.DeltaX
			y += step.direction.DeltaY
			f.tiles[y][x] = Dug
		}
	}
}

func (f Field) WithinField(x, y int) bool {
	return 0 <= x && x < len(f.tiles[0]) && 0 <= y && y < len(f.tiles)
}

func (f Field) IsDugAt(x, y int) bool {
	return f.WithinField(x, y) && f.tiles[y][x] == Dug
}

func (f Field) NonIntersectingCorner(x, y int) bool {
	return f.IsDugAt(x, y) && f.IsDugAt(x, y-1) && f.IsDugAt(x+1, y) ||
		f.IsDugAt(x, y) && f.IsDugAt(x, y+1) && f.IsDugAt(x-1, y)
}

func (f *Field) FillInner() {
	diagStarts := make([]Coord, 0, len(f.tiles)+len(f.tiles[0])-1)
	for x := 0; x < len(f.tiles[0]); x++ {
		diagStarts = append(diagStarts, Coord{x, 0})
	}
	for y := 0; y < len(f.tiles); y++ {
		diagStarts = append(diagStarts, Coord{0, y})
	}

	toFill := make([]Coord, 0)

	for _, coord := range diagStarts {
		inside := false
		for f.WithinField(coord.x, coord.y) {

			if f.tiles[coord.y][coord.x] == Dug && !f.NonIntersectingCorner(coord.x, coord.y) {
				inside = !inside
			} else if inside {
				toFill = append(toFill, coord)
			}

			coord.x += 1
			coord.y += 1
		}
	}
	for _, coord := range toFill {
		f.tiles[coord.y][coord.x] = Dug
	}
}

func (f Field) CountDug() int {
	return lo.Sum(lo.Map(
		f.tiles, func(tiles []rune, index int) int {
			return lo.Count(tiles, Dug)
		}))
}

func (f Field) ToString() string {
	var sb strings.Builder
	for _, row := range f.tiles {
		for _, c := range row {
			sb.WriteRune(c)
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func main() {
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t18-lavaduct-lagoon/test.txt")
	if err != nil {
		log.Fatalf("Failed to read file")
	}
	ctx := TaskContext{
		directions: common.NewDirections(),
	}
	digSteps := lo.Map(rows, common.NoIndex(func(row string) DigStep {
		return ParseDigStep2(ctx, row)
	}))
	field := NewField(digSteps)
	field.DigSteps(digSteps)
	//fmt.Printf("%s\n\n", field.ToString())
	field.FillInner()
	//fmt.Printf("%s\n\n", field.ToString())
	fmt.Printf("Dug: %d\n", field.CountDug())
}
