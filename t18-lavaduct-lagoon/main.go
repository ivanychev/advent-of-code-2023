package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/samber/lo"
	"github.com/zyedidia/generic/queue"
	"log"
	"strconv"
	"strings"
)

type TaskContext struct {
	directions common.Directions
}

type Tile struct {
	startX, endX, startY, endY int
}

func (t Tile) Area() int64 {
	return int64(t.endX-t.startX) * int64(t.endY-t.startY)
}

func (t Tile) IsAtEdge(f Field) bool {
	return t.startX == f.minX || t.startY == f.minY ||
		t.startX == f.maxX || t.startY == f.maxY
}

func (t Tile) Adjacent(f Field) []Tile {
	foundNextEndX, nextEndX := f.distinctXset.Ceiling(t.endX + 1)
	if nextEndX == nil {
		if t.startX == f.maxX {
			nextEndX = nil
		} else {
			nextEndX = struct{}{}
			foundNextEndX = f.maxX
		}
	}
	foundNextEndY, nextEndY := f.distinctYset.Ceiling(t.endY + 1)
	if nextEndY == nil {
		if t.startY == f.maxY {
			nextEndY = nil
		} else {
			nextEndY = struct{}{}
			foundNextEndY = f.maxY
		}
	}

	foundPrevStartX, prevStartX := f.distinctXset.Floor(t.startX - 1)
	foundPrevStartY, prevStartY := f.distinctYset.Floor(t.startY - 1)
	tiles := make([]Tile, 0)
	if nextEndX != nil {
		tiles = append(tiles, Tile{startX: t.endX, endX: foundNextEndX.(int), startY: t.startY, endY: t.endY})
	}
	if prevStartX != nil {
		tiles = append(tiles, Tile{startX: foundPrevStartX.(int), endX: t.startX, startY: t.startY, endY: t.endY})
	}
	if nextEndY != nil {
		tiles = append(tiles, Tile{
			startX: t.startX,
			endX:   t.endX,
			startY: t.endY,
			endY:   foundNextEndY.(int),
		})
	}
	if prevStartY != nil {
		tiles = append(tiles, Tile{
			startX: t.startX,
			endX:   t.endX,
			startY: foundPrevStartY.(int),
			endY:   t.startY,
		})
	}
	return tiles
}

func (t Tile) Includes(other Tile) bool {
	return t.startX <= other.startX && t.endX >= other.endX &&
		t.startY <= other.startY && t.endY >= other.endY
}

func (t Tile) IsBar() bool {
	return t.endX-t.startX == 1 || t.endY-t.startY == 1
}

func (t Tile) IsVerticalBar() bool {
	return t.IsBar() && t.endX-t.startX == 1
}

func (t Tile) IsHorizontalBar() bool {
	return t.IsBar() && t.endY-t.startY == 1
}

func TileFromInclusiveCoords(x0, x1, y0, y1 int) Tile {
	return Tile{
		startX: min(x0, x1),
		endX:   max(x0, x1) + 1,
		startY: min(y0, y1),
		endY:   max(y0, y1) + 1,
	}
}

func TileComparator(a, b interface{}) int {
	aTile := a.(Tile)
	bTile := b.(Tile)

	if aTile.startX != bTile.startX {
		return aTile.startX - bTile.startX
	}
	if aTile.startY != bTile.startY {
		return aTile.startY - bTile.startY
	}
	return 0
}

type Field struct {
	minX, maxX, minY, maxY int
	distinctXset           *treemap.Map // int -> struct{}
	distinctYset           *treemap.Map // int -> struct{}
	xToVerticalBorders     *treemap.Map // int -> *treemap.Map[Tile -> struct{}]
	yToHorizontalBorders   *treemap.Map // int -> *treemap.Map[Tile -> struct{}]
}

func (f Field) IsBarHorizontalBorder(bar Tile) bool {
	if bar.IsHorizontalBar() {
		row, found := f.yToHorizontalBorders.Get(bar.startY)
		if !found {
			return false
		}
		foundKey, foundValue := row.(*treemap.Map).Floor(bar)
		if foundValue == nil {
			return false
		}
		foundTile := foundKey.(Tile)
		return foundTile.Includes(bar)
	}
	return false
}

func (f Field) IsBarVerticalBorder(bar Tile) bool {
	if bar.IsVerticalBar() {
		column, found := f.xToVerticalBorders.Get(bar.startX)
		if !found {
			return false
		}
		foundKey, foundValue := column.(*treemap.Map).Floor(bar)
		if foundValue == nil {
			return false
		}
		foundTile := foundKey.(Tile)
		return foundTile.Includes(bar)
	}
	return false
}

func (f Field) IsBarBorder(bar Tile) bool {
	if !bar.IsBar() {
		return false
	}
	return f.IsBarHorizontalBorder(bar) || f.IsBarVerticalBorder(bar)
}

func (f Field) Tiles() []Tile {
	tiles := make([]Tile, 0, (f.distinctXset.Size()-1)*(f.distinctYset.Size()-1))

	xCoords := f.distinctXset.Keys()
	yCoords := f.distinctYset.Keys()

	for i := 0; i < len(xCoords); i++ {
		for j := 0; j < len(yCoords); j++ {
			startX := xCoords[i].(int)
			startY := yCoords[j].(int)

			var endX, endY int
			if i != len(xCoords)-1 {
				endX = xCoords[i+1].(int)
			} else {
				endX = startX + 1
			}

			if j != len(yCoords)-1 {
				endY = yCoords[j+1].(int)
			} else {
				endY = startY + 1
			}
			tile := Tile{
				startX: startX,
				startY: startY,
				endX:   endX,
				endY:   endY,
			}
			tiles = append(tiles, tile)
		}
	}
	return tiles
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

func ParseDigStepSimple(ctx TaskContext, s string) DigStep {
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

func ParseDigStep(ctx TaskContext, s string) DigStep {
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
	var x, y int

	distinctX := treemap.NewWithIntComparator()
	distinctY := treemap.NewWithIntComparator()
	xToVerticalBorders := treemap.NewWithIntComparator()
	yToHorizontalBorders := treemap.NewWithIntComparator()

	for _, step := range digSteps {
		newX := x + step.direction.DeltaX*step.size
		newY := y + step.direction.DeltaY*step.size

		if newX == x {
			insertTileToMap(xToVerticalBorders, x, newX, y, newY, x)
		} else {
			insertTileToMap(yToHorizontalBorders, x, newX, y, newY, y)
		}

		x = newX
		y = newY

		distinctX.Put(x, struct{}{})
		distinctX.Put(x-1, struct{}{})
		distinctX.Put(x+1, struct{}{})
		distinctY.Put(y, struct{}{})
		distinctY.Put(y-1, struct{}{})
		distinctY.Put(y+1, struct{}{})
	}

	minX, _ := distinctX.Min()
	maxX, _ := distinctX.Max()
	minY, _ := distinctY.Min()
	maxY, _ := distinctY.Max()
	return Field{
		minX:                 minX.(int),
		minY:                 minY.(int),
		maxX:                 maxX.(int),
		maxY:                 maxY.(int),
		distinctXset:         distinctX,
		distinctYset:         distinctY,
		yToHorizontalBorders: yToHorizontalBorders,
		xToVerticalBorders:   xToVerticalBorders,
	}
}

func insertTileToMap(m *treemap.Map, x int, newX int, y int, newY int, mapKey int) {
	row, found := m.Get(mapKey)
	if !found {
		row = treemap.NewWith(TileComparator)
		m.Put(mapKey, row)
	}
	row.(*treemap.Map).Put(TileFromInclusiveCoords(x, newX, y, newY), struct{}{})
}

func DebugRenderTiles(allTiles []Tile, borderTiles []Tile, internalTiles []Tile) string {
	xs := lo.FlatMap(allTiles, func(item Tile, index int) []int {
		return []int{item.startX, item.endX}
	})
	ys := lo.FlatMap(allTiles, func(item Tile, index int) []int {
		return []int{item.startY, item.endY}
	})
	xMax, xMin := lo.Max(xs), lo.Min(xs)
	yMax, yMin := lo.Max(ys), lo.Min(ys)
	m := common.CreateRuneMatrix(xMax-xMin+1, yMax-yMin+1, '.')
	for _, tile := range borderTiles {
		for x := tile.startX; x < tile.endX; x++ {
			for y := tile.startY; y < tile.endY; y++ {
				m[y-yMin][x-xMin] = '#'
			}
		}
	}
	for _, tile := range internalTiles {
		for x := tile.startX; x < tile.endX; x++ {
			for y := tile.startY; y < tile.endY; y++ {
				m[y-yMin][x-xMin] = '!'
			}
		}
	}
	return common.RuneMatrixToString(m)
}

func (f Field) Bfs(start Tile, visited *map[Tile]struct{}, visitor func(Tile)) {
	q := queue.New[Tile]()
	q.Enqueue(start)
	for !q.Empty() {
		curr := q.Dequeue()
		if _, already := (*visited)[curr]; already {
			continue
		}
		(*visited)[curr] = struct{}{}
		visitor(curr)
		adjacent := lo.Filter(curr.Adjacent(f), func(item Tile, index int) bool {
			return !f.IsBarBorder(item)
		})
		for _, tile := range adjacent {
			q.Enqueue(tile)
		}
	}
}

func main() {
	//rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t18-lavaduct-lagoon/test.txt")
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t18-lavaduct-lagoon/1.txt")
	if err != nil {
		log.Fatalf("Failed to read file")
	}
	ctx := TaskContext{
		directions: common.NewDirections(),
	}
	digSteps := lo.Map(rows, common.NoIndex(func(row string) DigStep {
		return ParseDigStep(ctx, row)
	}))
	field := NewField(digSteps)
	tiles := field.Tiles()
	tileToVisited := make(map[Tile]struct{})
	internalTiles := make([]Tile, 0)
	for _, tile := range tiles {
		if _, already := tileToVisited[tile]; already {
			continue
		}
		if field.IsBarBorder(tile) {
			continue
		}

		encounteredTiles := make([]Tile, 0)
		isInternal := true
		field.Bfs(tile, &tileToVisited, func(t Tile) {
			encounteredTiles = append(encounteredTiles, t)
			isInternal = isInternal && !t.IsAtEdge(field)
		})
		if isInternal {
			internalTiles = append(internalTiles, encounteredTiles...)
		}
	}

	borders := lo.Filter(tiles, func(t Tile, index int) bool {
		return field.IsBarBorder(t)
	})
	dag := make([]Tile, 0)
	dag = append(dag, internalTiles...)
	dag = append(dag, borders...)

	fmt.Printf("Anser: %d", lo.Sum(lo.Map(dag, common.NoIndex(Tile.Area))))
}
