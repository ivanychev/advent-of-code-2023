package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"log"
)

type Tile rune

const (
	FIXED Tile = '#'
	STONE Tile = 'O'
	SPACE Tile = '.'
)

type Field struct {
	tiles [][]Tile
}

func (f Field) WithinField(x, y int) bool {
	return 0 <= x && x < len(f.tiles[0]) && 0 <= y && y < len(f.tiles)
}

func (f Field) IsOccupied(x, y int) bool {
	return f.tiles[y][x] != SPACE
}

func (f *Field) RollStone(x, y, deltaX, deltaY int) {
	f.tiles[y][x] = SPACE
	for f.WithinField(x+deltaX, y+deltaY) && !f.IsOccupied(x+deltaX, y+deltaY) {
		x += deltaX
		y += deltaY
	}
	f.tiles[y][x] = STONE
}

func (f *Field) RollNorth() {
	for y := 0; y < len(f.tiles); y++ {
		for x := 0; x < len(f.tiles[0]); x++ {
			if f.tiles[y][x] == STONE {
				f.RollStone(x, y, 0, -1)
			}
		}
	}
}

func (f Field) CalculateScore() int {
	total := 0
	for y := 0; y < len(f.tiles); y++ {
		for x := 0; x < len(f.tiles[0]); x++ {
			if f.tiles[y][x] == STONE {
				total += len(f.tiles) - y
			}
		}
	}
	return total
}

func main() {
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t14-reflector/1.txt")
	if err != nil {
		log.Fatalf("Couldn't read file %w", err)
	}
	field := Field{
		tiles: lo.Map(rows, common.NoIndex(func(s string) []Tile {
			return []Tile(s)
		})),
	}
	field.RollNorth()
	fmt.Printf("Total: %d\n", field.CalculateScore())
}
