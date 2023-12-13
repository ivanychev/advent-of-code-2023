package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"io"
	"log"
	"os"
	"strings"
)

type Direction int

const (
	HORIZONTAL Direction = 1
	VERTICAL   Direction = 2
)

type Intersection struct {
	direction            Direction
	secondPartStartCoord int
}

func (i Intersection) ToValue() int {
	if i.direction == VERTICAL {
		return i.secondPartStartCoord
	}
	return i.secondPartStartCoord * 100
}

func CheckXMirror(m [][]rune, atX int, requiredErrors int) bool {
	leftColumn := atX
	rightColumn := atX + 1
	encounteredErrors := 0
	if rightColumn >= len(m[0]) {
		return false
	}
	for leftColumn >= 0 && rightColumn < len(m[0]) {
		for y := 0; y < len(m); y++ {
			if m[y][leftColumn] != m[y][rightColumn] {
				encounteredErrors += 1
			}
			if encounteredErrors > requiredErrors {
				return false
			}
		}
		leftColumn--
		rightColumn++
	}
	return requiredErrors == encounteredErrors
}

func CheckYMirror(m [][]rune, atY int, requiredErrors int) bool {
	upperRow := atY
	lowerRow := atY + 1
	encounteredErrors := 0
	if lowerRow >= len(m) {
		return false
	}
	for upperRow >= 0 && lowerRow < len(m) {
		for x := 0; x < len(m[0]); x++ {
			if m[upperRow][x] != m[lowerRow][x] {
				encounteredErrors += 1
			}
			if encounteredErrors > requiredErrors {
				return false
			}
		}
		upperRow--
		lowerRow++
	}
	return requiredErrors == encounteredErrors
}

func findIntersections(rawMap string, requiredErrors int) []Intersection {
	rowsAsStrings := strings.Split(rawMap, "\n")
	rows := lo.Map(rowsAsStrings, common.NoIndex(func(s string) []rune {
		return []rune(s)
	}))
	intersections := make([]Intersection, 0)
	for x := 0; x < len(rows[0])-1; x++ {
		if CheckXMirror(rows, x, requiredErrors) {
			intersections = append(intersections, Intersection{VERTICAL, x + 1})
		}
	}
	for y := 0; y < len(rows)-1; y++ {
		if CheckYMirror(rows, y, requiredErrors) {
			intersections = append(intersections, Intersection{HORIZONTAL, y + 1})
		}
	}
	return intersections
}

func main() {
	file, err := os.Open("/Users/iv/Code/advent-of-code-2023/t13-point-of-incidence/1.txt")
	//file, err := os.Open("/Users/iv/Code/advent-of-code-2023/t13-point-of-incidence/test.txt")
	if err != nil {
		log.Fatalf("Failed to open file")
	}
	allMaps, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file")
	}
	rawMaps := strings.Split(string(allMaps), "\n\n")
	intersections := lo.FlatMap(rawMaps, func(m string, index int) []Intersection {
		return findIntersections(m, 0)
	})
	fmt.Printf("%d\n", lo.SumBy(intersections, func(i Intersection) int {
		return i.ToValue()
	}))
}
