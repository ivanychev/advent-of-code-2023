package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"log"
)

type Direction string

const (
	UP    Direction = "UP"
	DOWN  Direction = "DOWN"
	LEFT  Direction = "LEFT"
	RIGHT Direction = "RIGHT"

	UD   = '|'
	LR   = '-'
	UR   = 'L'
	UL   = 'J'
	DL   = '7'
	DR   = 'F'
	NOOP = '.'
)

type Coord struct {
	x, y int
}

func tileToDirections(tile rune) []Direction {
	switch tile {
	case UD:
		return []Direction{UP, DOWN}
	case LR:
		return []Direction{LEFT, RIGHT}
	case UR:
		return []Direction{UP, RIGHT}
	case UL:
		return []Direction{UP, LEFT}
	case DL:
		return []Direction{DOWN, LEFT}
	case DR:
		return []Direction{DOWN, RIGHT}
	case NOOP:
		return nil
	}
	return nil
}

func getOppositeDirection(d Direction) Direction {
	var opposite Direction
	switch d {
	case UP:
		opposite = DOWN
	case DOWN:
		opposite = UP
	case LEFT:
		opposite = RIGHT
	case RIGHT:
		opposite = LEFT
	default:
		log.Fatalf("Unexpected")
	}
	return opposite
}

func findStartingCoords(tiles [][]rune) Coord {
	for y := 0; y < len(tiles); y++ {
		for x := 0; x < len(tiles[0]); x++ {
			if tiles[y][x] == 'S' {
				return Coord{x, y}
			}
		}
	}
	log.Fatalf("Failed to find")
	return Coord{}
}

type Field struct {
	tiles         [][]rune
	startingCoord Coord
}

func (f Field) Move(from Coord, notTo *Direction) (Coord, *Direction, *Direction) {
	directions := lo.Filter(tileToDirections(f.tiles[from.y][from.x]), func(item Direction, index int) bool {
		return notTo == nil || item != *notTo
	})
	if notTo == nil && len(directions) != 2 {
		log.Fatalf("Unexpected %d", len(directions))
	}
	if notTo != nil && len(directions) != 1 {
		log.Fatalf("Unexpected %d", len(directions))
	}
	directionToMove := directions[0]
	switch directionToMove {
	case UP:
		from.y -= 1
	case DOWN:
		from.y += 1
	case RIGHT:
		from.x += 1
	case LEFT:
		from.x -= 1
	default:
		log.Fatalf("Unexpected %v", directionToMove)
	}
	opposite := getOppositeDirection(directionToMove)
	return from, &directionToMove, &opposite
}

func main() {
	//rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t10-pipes/test.txt")
	//startIsRune := DR
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t10-pipes/1.txt")
	startIsRune := UD
	if err != nil {
		log.Fatalf("Failed to open file: %w", err)
	}
	tiles := lo.Map(rows, common.NoIndex(func(s string) []rune {
		return []rune(s)
	}))
	startingCoord := findStartingCoords(tiles)
	tiles[startingCoord.y][startingCoord.x] = startIsRune
	field := Field{tiles: tiles, startingCoord: startingCoord}

	moved := false
	slowerPointer := startingCoord
	fasterPointer := startingCoord
	var lastSlowDirection, lastSlowOppositeDirection *Direction
	var lastFastDirection, lastFastOppositeDirection *Direction
	slowerMoved := 0
	for !moved || fasterPointer != startingCoord {
		slowerPointer, lastSlowDirection, lastSlowOppositeDirection = field.Move(slowerPointer, lastSlowOppositeDirection)
		fasterPointer, lastFastDirection, lastFastOppositeDirection = field.Move(fasterPointer, lastFastOppositeDirection)
		fasterPointer, lastFastDirection, lastFastOppositeDirection = field.Move(fasterPointer, lastFastOppositeDirection)
		slowerMoved++
		moved = true
	}
	fmt.Printf("%d\n", slowerMoved)
	fmt.Printf("%v %v\n", lastSlowDirection, lastFastDirection)
}
