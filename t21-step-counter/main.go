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
	Steps = 64
)

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

func (f Field) FrontStep(coords []common.Coord) []common.Coord {
	nextCoordSet := make(map[common.Coord]struct{})
	for _, c := range coords {
		for _, newCoord := range common.AdjacentCoords(c, f.tiles) {
			if f.tiles[newCoord.Y][newCoord.X] != Rock {
				nextCoordSet[newCoord] = struct{}{}
			}
		}
	}
	return lo.MapToSlice(nextCoordSet, func(key common.Coord, value struct{}) common.Coord {
		return key
	})
}

func main() {
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t21-step-counter/1.txt")
	if err != nil {
		log.Fatalf("Failed to read file")
	}
	field := ParseField(rows)
	coords := []common.Coord{{field.startX, field.startY}}
	for i := 0; i < Steps; i++ {
		coords = field.FrontStep(coords)
	}
	fmt.Printf("Total: %d\n", len(coords))
}
