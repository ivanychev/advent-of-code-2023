package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"log"
)

const (
	Empty = '.'
	Rock  = '#'
	Start = 'S'
	Steps = 26501365
)

type FlipCount struct {
	Vertical   int
	Horizontal int
}

type CoordState struct {
	x, y       int
	flipCounts map[FlipCount]int
}

type Field struct {
	startX, startY int
	tiles          [][]rune
}

func ParseField(rows []string) Field {
	tiles := common.CreateRuneMatrix(len(rows[0]), len(rows), Empty)
	var startX, startY int
	for x := 0; x < len(rows[0]); x++ {
		for y := 0; y < len(rows); y++ {
			switch rows[y][x] {
			case Empty:
			case Start:
				startX = x
				startY = y
			case Rock:
				tiles[y][x] = Rock
			}
		}
	}
	return Field{startY: startY, startX: startX, tiles: tiles}
}

func RoundCoord(c common.Coord, tiles [][]rune) common.Coord {
	for c.X < 0 {
		c.X += len(tiles[0])
	}
	for c.X >= len(tiles[0]) {
		c.X -= len(tiles[0])
	}
	for c.Y < 0 {
		c.Y += len(tiles)
	}
	for c.Y >= len(tiles) {
		c.Y -= len(tiles)
	}
	return c
}

func ChangeVerticals(m map[FlipCount]int, delta int) map[FlipCount]int {
	return lo.MapKeys(m, func(value int, key FlipCount) FlipCount {
		return FlipCount{Vertical: key.Vertical + delta, Horizontal: key.Horizontal}
	})
}

func ChangeHorizontals(m map[FlipCount]int, delta int) map[FlipCount]int {
	return lo.MapKeys(m, func(value int, key FlipCount) FlipCount {
		return FlipCount{Vertical: key.Vertical, Horizontal: key.Horizontal + delta}
	})
}

func AdjacentCoords(c CoordState, tiles [][]rune) []CoordState {
	states := make([]CoordState, 0, 4)
	if c.x+1 == len(tiles[0]) {
		states = append(states, CoordState{x: 0, y: c.y, flipCounts: ChangeHorizontals(c.flipCounts, 1)})
	} else {
		states = append(states, CoordState{x: c.x + 1, y: c.y, flipCounts: c.flipCounts})
	}

	if c.x-1 == -1 {
		states = append(states, CoordState{x: len(tiles[0]) - 1, y: c.y, flipCounts: ChangeHorizontals(c.flipCounts, -1)})
	} else {
		states = append(states, CoordState{x: c.x - 1, y: c.y, flipCounts: c.flipCounts})
	}

	if c.y+1 == len(tiles) {
		states = append(states, CoordState{x: c.x, y: 0, flipCounts: ChangeVerticals(c.flipCounts, 1)})
	} else {
		states = append(states, CoordState{x: c.x, y: c.y + 1, flipCounts: c.flipCounts})
	}

	if c.y-1 == -1 {
		states = append(states, CoordState{x: c.x, y: len(tiles) - 1, flipCounts: ChangeVerticals(c.flipCounts, -1)})
	} else {
		states = append(states, CoordState{x: c.x, y: c.y - 1, flipCounts: c.flipCounts})
	}
	return states
}

func SumTwoStates(a, b CoordState) CoordState {
	if a.x != b.x || a.y != b.y {
		log.Fatalf("Failed to sum")
	}

	flipCounts := a.flipCounts
	for k, v := range b.flipCounts {
		_, exists := flipCounts[k]
		if !exists {
			flipCounts[k] = v
		} else {
			flipCounts[k] += v
		}
	}

	return CoordState{x: a.x, y: a.y, flipCounts: flipCounts}
}

func SumStates(s []CoordState) CoordState {
	return lo.Reduce(s, func(agg CoordState, item CoordState, index int) CoordState {
		return SumTwoStates(agg, item)
	}, CoordState{s[0].x, s[0].y, make(map[FlipCount]int)})
}

func (f Field) FrontStep(states map[common.Coord]CoordState) map[common.Coord]CoordState {
	newStates := lo.FlatMap(lo.Values(states), func(item CoordState, index int) []CoordState {
		return lo.Filter(AdjacentCoords(item, f.tiles), func(s CoordState, index int) bool {
			return f.tiles[s.y][s.x] != Rock
		})
	})
	coordToStates := lo.GroupBy(newStates, func(s CoordState) common.Coord {
		return common.Coord{X: s.x, Y: s.y}
	})
	coordToFinalState := lo.MapValues(coordToStates, func(states []CoordState, key common.Coord) CoordState {
		return SumStates(states)
	})
	return coordToFinalState
}

func ComputeTotal(m map[common.Coord]CoordState) int64 {
	var total int64
	for _, state := range m {
		total += int64(len(state.flipCounts))
	}
	return total
}

func main() {
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t21-step-counter/1.txt")
	if err != nil {
		log.Fatalf("Failed to read file")
	}
	field := ParseField(rows)
	set := make(map[common.Coord]CoordState)
	set[common.Coord{X: field.startX, Y: field.startY}] = CoordState{x: field.startX, y: field.startY,
		flipCounts: map[FlipCount]int{FlipCount{Vertical: 0, Horizontal: 0}: 1}}
	fmt.Printf("%d %d\n", len(field.tiles[0]), len(field.tiles))
	for i := 0; i < Steps; i++ {
		set = field.FrontStep(set)
		if (i+1)%131 == 65 {
			fmt.Printf("%d cycles, steps %d: %d\n", (i+1)/131, i, ComputeTotal(set))
		}
	}
	//fmt.Printf("Total: %d\n", ComputeTotal(set))
}

// Observe and run some python:

//a = 121888
//b = 91702
//k=3
//
//i = 6
//total = 644659
//CYCLE_SIZE = 131
//MAX_STEPS = 26501365
//steps = (131*i + 65) - 1
//
//print(f"{steps=} {i=} {total}")
//while True:
//i += 2
//steps = (131*i + 65) - 1
//total += k*a + b
//k += 1
//print(f"{steps=} {i=} {total}")
//if steps >= MAX_STEPS:
//break
