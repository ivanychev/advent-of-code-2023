package main

import (
	"advent_of_code/common"
	"encoding/binary"
	"fmt"
	"github.com/samber/lo"
	"hash/fnv"
	"log"
	"slices"
	"strings"
)

type Direction struct {
	deltaX, deltaY int
	name           string
}
type Coord struct {
	x, y int
}

type DirectionName string

const (
	UP    DirectionName = "UP"
	DOWN  DirectionName = "DOWN"
	LEFT  DirectionName = "LEFT"
	RIGHT DirectionName = "RIGHT"
)

func DirectionNameToChar(d DirectionName) rune {
	switch d {
	case UP:
		return '^'
	case LEFT:
		return '<'
	case DOWN:
		return 'v'
	case RIGHT:
		return '>'
	}
	log.Fatalf("dfsdfsd")
	return '0'
}

func DirectionToDelta(d DirectionName) (int, int) {
	switch d {
	case UP:
		return 0, -1
	case DOWN:
		return 0, 1
	case LEFT:
		return -1, 0
	case RIGHT:
		return 1, 0
	}
	log.Fatalf("Unreachible")
	return 0, 0
}

type BeamCoord struct {
	direction DirectionName
	x, y      int
}

func CompareBeamCoord(a, b BeamCoord) int {
	if a.x != b.x {
		return a.x - b.x
	}
	if a.y != b.y {
		return a.y - b.y
	}
	aDeltaX, aDeltaY := DirectionToDelta(a.direction)
	bDeltaX, bDeltaY := DirectionToDelta(b.direction)
	if aDeltaX != bDeltaX {
		return aDeltaX - bDeltaX
	}
	return aDeltaY - bDeltaY
}

func (c BeamCoord) ToUp() BeamCoord {
	return BeamCoord{
		UP, c.x, c.y - 1,
	}
}

func (c BeamCoord) ToDown() BeamCoord {
	return BeamCoord{
		DOWN, c.x, c.y + 1,
	}
}

func (c BeamCoord) ToRight() BeamCoord {
	return BeamCoord{
		RIGHT, c.x + 1, c.y,
	}
}

func (c BeamCoord) ToLeft() BeamCoord {
	return BeamCoord{
		LEFT, c.x - 1, c.y,
	}
}

type BeamField struct {
	tiles [][]rune
}

func (f BeamField) WithinField(x, y int) bool {
	return 0 <= x && x < len(f.tiles[0]) && 0 <= y && y < len(f.tiles)
}

func (c BeamCoord) FlyForward(f BeamField) []BeamCoord {
	deltaX, deltaY := DirectionToDelta(c.direction)
	newX, newY := c.x+deltaX, c.y+deltaY
	if !f.WithinField(newX, newY) {
		return nil
	} else {
		return []BeamCoord{{
			c.direction,
			newX, newY,
		}}
	}
}

func (c BeamCoord) WithinField(f BeamField) bool {
	return f.WithinField(c.x, c.y)
}

func FilterCorrectCoords(coords []BeamCoord, f BeamField) []BeamCoord {
	return lo.Filter(coords, common.NoIndex(func(c BeamCoord) bool {
		return c.WithinField(f)
	}))
}

func (c BeamCoord) Step(f BeamField) []BeamCoord {
	tile := f.tiles[c.y][c.x]

	switch tile {
	case '.':
		return c.FlyForward(f)
	case '|':
		switch c.direction {
		case UP, DOWN:
			return c.FlyForward(f)
		case LEFT, RIGHT:
			return FilterCorrectCoords(
				[]BeamCoord{
					c.ToUp(),
					c.ToDown()}, f)
		}
	case '-':
		switch c.direction {
		case LEFT, RIGHT:
			return c.FlyForward(f)
		case UP, DOWN:
			return FilterCorrectCoords([]BeamCoord{
				c.ToLeft(),
				c.ToRight()}, f)
		}
	case '\\':
		switch c.direction {
		case UP:
			return FilterCorrectCoords([]BeamCoord{c.ToLeft()}, f)
		case RIGHT:
			return FilterCorrectCoords([]BeamCoord{c.ToDown()}, f)
		case DOWN:
			return FilterCorrectCoords([]BeamCoord{c.ToRight()}, f)
		case LEFT:
			return FilterCorrectCoords([]BeamCoord{c.ToUp()}, f)
		}
	case '/':
		switch c.direction {
		case UP:
			return FilterCorrectCoords([]BeamCoord{c.ToRight()}, f)
		case RIGHT:
			return FilterCorrectCoords([]BeamCoord{c.ToUp()}, f)
		case DOWN:
			return FilterCorrectCoords([]BeamCoord{c.ToLeft()}, f)
		case LEFT:
			return FilterCorrectCoords([]BeamCoord{c.ToDown()}, f)
		}
	}
	log.Fatalf("Unreachible")
	return nil
}

func HashCoords(coords []BeamCoord) uint64 {
	hash := fnv.New64()
	b := make([]byte, 8)

	for _, coord := range coords {
		binary.LittleEndian.PutUint64(b, uint64(coord.x))
		hash.Write(b)
		binary.LittleEndian.PutUint64(b, uint64(coord.y))
		hash.Write(b)
		hash.Write([]byte(coord.direction))
	}
	return hash.Sum64()
}

func Draw(f BeamField, beams []BeamCoord) string {
	chars := make([][]rune, 0)
	for y := 0; y < len(f.tiles); y++ {
		chars = append(chars, make([]rune, 0))
		for x := 0; x < len(f.tiles[0]); x++ {
			chars[y] = append(chars[y], f.tiles[y][x])
		}
	}
	for _, b := range beams {
		chars[b.y][b.x] = DirectionNameToChar(b.direction)
	}

	var sb strings.Builder
	for _, bts := range chars {
		sb.WriteString(string(bts))
		sb.WriteRune('\n')
	}
	return sb.String()
}

func main() {
	//rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t16-beams/test.txt")
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t16-beams/1.txt")
	field := BeamField{tiles: lo.Map(rows, func(s string, index int) []rune {
		return []rune(s)
	})}
	if err != nil {
		log.Fatalf("Failed to open file: %w", err)
	}
	encounteredHashes := make(map[uint64][]BeamCoord)
	encounteredCoords := make(map[Coord]int)
	coords := []BeamCoord{
		{RIGHT, 0, 0},
	}
	hash := HashCoords(coords)
	if err != nil {
		log.Fatalf("Failed %w", err)
	}
	exists := false
	total := 0

	//fmt.Printf("%s\n\n", Draw(field, coords))
	//fmt.Printf("")
	for !exists {
		for _, c := range coords {
			encounteredCoords[Coord{c.x, c.y}]++
		}
		encounteredHashes[hash] = coords
		if total == 3 {
			fmt.Printf("Here")
		}
		coords = lo.FlatMap(coords, func(c BeamCoord, index int) []BeamCoord {
			return c.Step(field)
		})
		coords = lo.Uniq(coords)
		slices.SortFunc(coords, CompareBeamCoord)
		//fmt.Printf("%s\n\n", Draw(field, coords))
		//fmt.Printf("")

		hash = HashCoords(coords)
		if err != nil {
			log.Fatalf("Failed %w", err)
		}
		_, exists = encounteredHashes[hash]
		if exists {
			fmt.Printf("Here")
		}
		total += 1
	}
	fmt.Printf("Energized: %d\n", len(encounteredCoords))
}
