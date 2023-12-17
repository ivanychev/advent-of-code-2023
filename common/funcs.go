package common

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type DirectionName string

const (
	UP    DirectionName = "UP"
	DOWN  DirectionName = "DOWN"
	LEFT  DirectionName = "LEFT"
	RIGHT DirectionName = "RIGHT"
)

type DirectionDesc struct {
	Name           DirectionName
	Char           rune
	DeltaX, DeltaY int
	*Directions
}

func (d DirectionDesc) Turns() [2]DirectionDesc {
	return d.Directions.Turns(d)
}

type Directions struct {
	Up, Down, Left, Right DirectionDesc
}

func NewDirections() Directions {
	dirs := Directions{
		Up: DirectionDesc{
			UP, '^', 0, -1, nil,
		},
		Down: DirectionDesc{
			DOWN, 'v', 0, 1, nil,
		},
		Left: DirectionDesc{
			LEFT, '<', -1, 0, nil,
		},
		Right: DirectionDesc{
			RIGHT, '>', 1, 0, nil,
		},
	}
	dirs.Left.Directions = &dirs
	dirs.Right.Directions = &dirs
	dirs.Down.Directions = &dirs
	dirs.Up.Directions = &dirs

	return dirs
}

func (d Directions) Turns(desc DirectionDesc) [2]DirectionDesc {
	switch desc.Char {
	case '^', 'v':
		return [2]DirectionDesc{d.Left, d.Right}
	case '>', '<':
		return [2]DirectionDesc{d.Up, d.Down}
	}
	log.Fatalf("Unreacheable")
	return [2]DirectionDesc{d.Up, d.Down}
}

type Number interface {
	int
}

func Sum[T Number](items []T) T {
	result := *new(T)
	for _, item := range items {
		result += item
	}
	return result
}

func MaxPair[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func MinPair[T Number](a, b T) T {
	if a > b {
		return b
	}
	return a
}

func NoIndex[T, R any](f func(T) R) func(T, int) R {
	return func(t T, _ int) R {
		return f(t)
	}
}

func FileToRows(path string) ([]string, error) {
	readFile, err := os.Open(path)

	if err != nil {
		return []string{}, fmt.Errorf("Failed to read the file: %w", err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	return lines, nil
}

func StringOfNumbersToInts(s string) []int {
	parts := strings.Fields(s)
	var returned = make([]int, 0, len(parts))
	for _, part := range parts {
		val, err := strconv.Atoi(part)
		if err != nil {
			log.Fatalf("Failed to parse int: %s", part)
		}
		returned = append(returned, val)
	}
	return returned
}
