package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/samber/lo"
	"log"
	"strconv"
	"unicode"
)

type Schematic struct {
	Rows []string
}

type Number struct {
	X      int
	Y      int
	Length int
	Value  int
}

type Coord struct {
	X int
	Y int
}

func (s Schematic) LenX() int {
	return len(s.Rows[0])
}

func (s Schematic) LenY() int {
	return len(s.Rows)
}

func (s Schematic) IsValidCoord(x, y int) bool {
	return 0 <= x && x < s.LenY() && 0 <= y && y < s.LenY()
}

func ExtractNumbersFromRow(row string, y int, re *regexp2.Regexp) []Number {
	var numbers []Number
	m, _ := re.FindStringMatch(row)
	for m != nil {
		group := m.Groups()[0]
		value, err := strconv.Atoi(group.String())
		if err != nil {
			log.Fatalf("Failed to parse int from %s", group.String())
		}
		numbers = append(numbers, Number{
			X:      group.Index,
			Y:      y,
			Length: group.Length,
			Value:  value,
		})
		m, _ = re.FindNextMatch(m)
	}
	return numbers
}

func AdjacentCells(n Number, s Schematic) []Coord {
	var coords []Coord
	var x, y int

	x = n.X - 1
	y = n.Y
	if s.IsValidCoord(x, y) {
		coords = append(coords, Coord{x, y})
	}

	x = n.X + n.Length
	y = n.Y
	if s.IsValidCoord(x, y) {
		coords = append(coords, Coord{x, y})
	}

	for i := -1; i <= n.Length; i++ {
		x = n.X + i
		y = n.Y - 1
		if s.IsValidCoord(x, y) {
			coords = append(coords, Coord{x, y})
		}

		x = n.X + i
		y = n.Y + 1
		if s.IsValidCoord(x, y) {
			coords = append(coords, Coord{x, y})
		}
	}
	return coords
}

func IsPartNumber(n Number, s Schematic) bool {
	for _, coord := range AdjacentCells(n, s) {
		s := rune(s.Rows[coord.Y][coord.X])
		if s != '.' && !unicode.IsDigit(s) {
			return true
		}
	}
	return false
}

func main() {
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t3-gear/1.txt")
	if err != nil {
		log.Fatalf("Failed to read file: %w", err)
	}
	numbersRe := regexp2.MustCompile("\\d+", regexp2.IgnoreCase)
	schematic := Schematic{Rows: rows}
	numbers := lo.FlatMap(schematic.Rows, func(item string, index int) []Number {
		return ExtractNumbersFromRow(item, index, numbersRe)
	})
	partNumbers := lo.Filter(numbers, func(item Number, index int) bool {
		return IsPartNumber(item, schematic)
	})
	sum := lo.SumBy(partNumbers, func(item Number) int64 {
		return int64(item.Value)
	})
	fmt.Printf("Total sum: %v", sum)
}
