package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"log"
	"strings"
)

type Brick struct {
	X, Y, Z int
}

func (b Brick) StartCoord() Coord {
	return Coord{X: b.X, Y: b.Y, Z: b.Z}
}

type Coord struct {
	X, Y, Z int
}

type Rect struct {
	bricks []Brick
}

func (r Rect) Coords() []Coord {
	return lo.Map(r.bricks, common.NoIndex(Brick.StartCoord))
}

type Field struct {
	coordToRect map[Coord]*Rect
	rects       []Rect
	canFall     []*Rect
}

func (f Field) RectPts() []*Rect {
	return lo.Map(lo.Range(len(f.rects)), func(i int, index int) *Rect {
		return &f.rects[i]
	})
}

func (f Field) DeltaRect(r *Rect, deltaX, deltaY, deltaZ int) []*Rect {
	rects := make(map[*Rect]struct{}, 0)
	for _, coord := range r.Coords() {
		underCoord := Coord{coord.X + deltaX, coord.Y + deltaY, coord.Z + deltaZ}
		underR, exists := f.coordToRect[underCoord]
		if exists && r != underR {
			rects[underR] = struct{}{}
		}
	}
	return lo.Keys(rects)
}

func (f Field) UnderRect(r *Rect) []*Rect {
	return f.DeltaRect(r, 0, 0, -1)
}

func (f Field) AboveRect(r *Rect) []*Rect {
	return f.DeltaRect(r, 0, 0, 1)
}

func (r Rect) StandsOnGround() bool {
	return len(lo.Filter(r.Coords(), func(c Coord, index int) bool {
		return c.Z == 1
	})) != 0
}

func (f Field) CanFall(r *Rect) bool {
	return len(f.UnderRect(r)) == 0 && !r.StandsOnGround()
}

func (f Field) Fall(r *Rect) {
	for i := 0; i < len(r.bricks); i++ {
		delete(f.coordToRect, r.bricks[i].StartCoord())
	}
	for f.CanFall(r) {
		for i := 0; i < len(r.bricks); i++ {
			r.bricks[i].Z--
		}
	}
	for i := 0; i < len(r.bricks); i++ {
		f.coordToRect[r.bricks[i].StartCoord()] = r
	}
}

func ParseBrick(s string) Brick {
	// 2,1,6
	rawCoords := strings.Split(s, ",")
	coords := lo.Map(rawCoords, common.NoIndex(common.MustAtoi))
	return Brick{coords[0], coords[1], coords[2]}
}

func ParseRect(s string) Rect {
	// 0,1,6~2,1,6
	components := strings.Split(s, "~")
	startBricks := ParseBrick(components[0])
	endBricks := ParseBrick(components[1])

	bricks := make([]Brick, 0)
	for x := startBricks.X; x <= endBricks.X; x++ {
		for y := startBricks.Y; y <= endBricks.Y; y++ {
			for z := startBricks.Z; z <= endBricks.Z; z++ {
				bricks = append(bricks, Brick{X: x, Y: y, Z: z})
			}
		}
	}
	return Rect{bricks: bricks}
}

func BuildField(rects []Rect) Field {
	index := make(map[Coord]*Rect)
	for i := 0; i < len(rects); i++ {
		rectPtr := &rects[i]
		for _, brick := range (*rectPtr).bricks {
			index[brick.StartCoord()] = rectPtr
		}
	}

	f := Field{index, rects, nil}
	f.canFall = lo.FlatMap(lo.Range(len(f.rects)), func(i int, index int) []*Rect {
		if f.CanFall(&f.rects[i]) {
			return []*Rect{&f.rects[i]}
		} else {
			return nil
		}
	})

	return f
}

func main() {
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t22-sand-slabs/1.txt")
	if err != nil {
		log.Fatalf("failed to read file: %w", err)
	}
	rects := lo.Map(rows, common.NoIndex(ParseRect))
	field := BuildField(rects)

	for len(field.canFall) > 0 {
		falling := field.canFall[len(field.canFall)-1]
		field.canFall = field.canFall[:len(field.canFall)-1]
		candidates := field.AboveRect(falling)
		field.Fall(falling)
		for _, c := range candidates {
			if field.CanFall(c) {
				field.canFall = append(field.canFall, c)
			}
		}
	}

	belowSet := make(map[*Rect]int)

	for i, rect := range field.RectPts() {
		below := field.UnderRect(rect)
		if len(below) == 1 {
			belowSet[below[0]] = i
		}
	}
	//fmt.Printf("%+v\n", len(belowSet))
	fmt.Printf("Can be: %d", len(field.rects)-len(belowSet))
}
