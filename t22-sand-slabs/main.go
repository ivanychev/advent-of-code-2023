package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"github.com/zyedidia/generic/queue"
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

func EnsureInitSet(k *Rect, m *map[*Rect]map[*Rect]struct{}) {
	if _, exists := (*m)[k]; !exists {
		(*m)[k] = make(map[*Rect]struct{})
	}
}

func AddSliceToSet(key *Rect, values []*Rect, m *map[*Rect]map[*Rect]struct{}) {
	for _, v := range values {
		(*m)[key][v] = struct{}{}
	}
}

func (f *Field) FallAll() {
	for len(f.canFall) > 0 {
		falling := f.canFall[len(f.canFall)-1]
		f.canFall = f.canFall[:len(f.canFall)-1]
		candidates := f.AboveRect(falling)
		f.Fall(falling)
		for _, c := range candidates {
			if f.CanFall(c) {
				f.canFall = append(f.canFall, c)
			}
		}
	}
}

func (f *Field) BuildSupportIndexes() (map[*Rect]map[*Rect]struct{}, map[*Rect]map[*Rect]struct{}, map[*Rect]struct{}) {
	supportsMap := make(map[*Rect]map[*Rect]struct{})
	supportedByMap := make(map[*Rect]map[*Rect]struct{})
	supportedBySingle := make(map[*Rect]struct{})

	for _, rect := range f.RectPts() {
		unders := f.UnderRect(rect)
		EnsureInitSet(rect, &supportedByMap)
		AddSliceToSet(rect, unders, &supportedByMap)

		for _, support := range unders {
			EnsureInitSet(support, &supportsMap)
			supportsMap[support][rect] = struct{}{}
		}
	}

	for k, v := range supportedByMap {
		if len(v) == 1 {
			supportedBySingle[k] = struct{}{}
		}
	}

	return supportsMap, supportedByMap, supportedBySingle
}

func AnyKey[K comparable, V any](m map[K]V) K {
	for k, _ := range m {
		return k
	}
	log.Fatalf("Unreacheable")
	var k K
	return k
}

func (f Field) ComputeFallenBricks(rect *Rect) int {
	supportsMap, supportedByMap, _ := f.BuildSupportIndexes()
	q := queue.New[*Rect]()
	q.Enqueue(rect)
	total := -1

	for !q.Empty() {
		curr := q.Dequeue()
		total += 1
		currSupports := supportsMap[curr]
		for supports, _ := range currSupports {
			delete(supportedByMap[supports], curr)
			if len(supportedByMap[supports]) == 0 {
				q.Enqueue(supports)
			}
		}
		delete(supportsMap, curr)
	}
	return total
}

func main() {
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t22-sand-slabs/1.txt")
	if err != nil {
		log.Fatalf("failed to read file: %w", err)
	}
	rects := lo.Map(rows, common.NoIndex(ParseRect))
	field := BuildField(rects)
	field.FallAll()

	supportsMap, _, supportedBySingle := field.BuildSupportIndexes()

	total := 0
	for _, rect := range field.RectPts() {
		rectSupports := supportsMap[rect]
		willFall := make([]*Rect, 0)
		for supports, _ := range rectSupports {
			if _, exists := supportedBySingle[supports]; exists {
				willFall = append(willFall, supports)
			}
		}
		if len(willFall) == 0 {
			continue
		}
		total += field.ComputeFallenBricks(rect)
	}

	fmt.Printf("Will fall: %d", total)
}
