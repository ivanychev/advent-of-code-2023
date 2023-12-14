package main

import (
	"advent_of_code/common"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/samber/lo"
	"log"
	"strings"
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

func (f Field) Hash() string {
	hash := md5.Sum(lo.FlatMap(f.tiles, func(row []Tile, index int) []byte {
		return lo.Map(row, common.NoIndex(func(t Tile) byte {
			return byte(t)
		}))
	}))
	return hex.EncodeToString(hash[:])
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

func GetLoopCond(deltaCoord, max int) (int, func(int) bool, int) {
	step := -deltaCoord
	if step == 0 {
		step = 1
	}

	if step == 1 {
		return 0, func(i int) bool {
			return i < max
		}, step
	}
	return max - 1, func(i int) bool {
		return i >= 0
	}, step
}

func (f *Field) RollTo(deltaX, deltaY int) {
	x0, xCond, xStep := GetLoopCond(deltaX, len(f.tiles[0]))
	y0, yCond, yStep := GetLoopCond(deltaY, len(f.tiles))

	for y := y0; yCond(y); y += yStep {
		for x := x0; xCond(x); x += xStep {
			if f.tiles[y][x] == STONE {
				f.RollStone(x, y, deltaX, deltaY)
			}
		}
	}
}

func (f Field) ToString() string {
	b := strings.Builder{}
	for y := 0; y < len(f.tiles); y++ {
		for x := 0; x < len(f.tiles[0]); x++ {
			b.Write([]byte{byte(f.tiles[y][x])})
		}
		b.WriteString("\n")
	}
	return b.String()
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
	//rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t14-reflector/test.txt")
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t14-reflector/1.txt")
	if err != nil {
		log.Fatalf("Couldn't read file %w", err)
	}
	field := Field{
		tiles: lo.Map(rows, common.NoIndex(func(s string) []Tile {
			return []Tile(s)
		})),
	}
	hashes := make(map[string]int, 0)
	hashes[field.Hash()] = 0
	rotations := 1000000000
	//var cycleLength int
	for i := 0; i < rotations; i++ {
		field.RollTo(0, -1)
		//fmt.Printf("%s\n", field.ToString())
		field.RollTo(-1, 0)
		//fmt.Printf("%s\n", field.ToString())
		field.RollTo(0, 1)
		//fmt.Printf("%s\n", field.ToString())
		field.RollTo(1, 0)
		//fmt.Printf("%s\n", field.ToString())

		currentHash := field.Hash()
		prevIndex, exists := hashes[currentHash]
		if exists {
			cycleLength := i - prevIndex
			for i+cycleLength < rotations {
				i += cycleLength
			}
		} else {
			hashes[currentHash] = i
		}
	}
	//field.RollTo(0, -1)

	fmt.Printf("Total: %d\n", field.CalculateScore())
}
