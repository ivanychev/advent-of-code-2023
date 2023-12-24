package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"github.com/zyedidia/generic/stack"
	"log"
	"strconv"
	"strings"
)

const (
	Empty         = '.'
	Rock          = '#'
	ReplaceSlopes = false
)

type Field struct {
	startX, startY, endX, endY int
	tiles                      [][]rune
	directions                 common.Directions
}

func ParseField(rows []string, directions common.Directions) Field {
	if ReplaceSlopes {
		replacer := strings.NewReplacer(">", ".", "<", ".", "v", ".", "^", ".")
		rows = lo.Map(rows, func(s string, index int) string {
			return replacer.Replace(s)
		})
	}

	tiles := lo.Map(rows, func(row string, index int) []rune {
		return []rune(row)
	})
	startY := 0
	startX := lo.IndexOf(tiles[0], Empty)
	endY := len(tiles) - 1
	endX := lo.IndexOf(tiles[len(tiles)-1], Empty)

	return Field{
		startY:     startY,
		startX:     startX,
		endX:       endX,
		endY:       endY,
		tiles:      tiles,
		directions: directions,
	}
}

func (f Field) CanGoTo(at common.Coord, visited map[common.Coord]int) []common.Coord {
	runeAt := f.tiles[at.Y][at.X]
	var canGoTo []common.Coord
	if f.directions.IsSlope(runeAt) {
		direction := f.directions.SlopeToDirection(runeAt)
		canGoTo = []common.Coord{
			{at.X + direction.DeltaX, at.Y + direction.DeltaY},
		}
	} else {
		canGoTo = common.AdjacentCoords(at, f.tiles)
	}

	return lo.Filter(canGoTo, func(item common.Coord, index int) bool {
		if !(common.IsValidCoord(item, f.tiles) && f.tiles[item.Y][item.X] != Rock) {
			return false
		}
		_, exists := visited[item]

		return !exists
	})
}

type Step struct {
	common.Coord
	goTo           []common.Coord
	currentlyInIdx int
}

func (s Step) AllVisited() bool {
	return s.currentlyInIdx == len(s.goTo)-1
}

func (f Field) DebugPrint(path map[common.Coord]int) string {
	var sb strings.Builder
	for y := 0; y < len(f.tiles); y++ {
		for x := 0; x < len(f.tiles[0]); x++ {
			if x == f.startX && y == f.startY {
				sb.WriteRune('S')
			} else if val, exists := path[common.Coord{X: x, Y: y}]; exists {
				sb.WriteString(strconv.Itoa(val % 10))
			} else {
				sb.WriteRune(f.tiles[y][x])
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (f Field) DebugPrintStack(s stack.Stack[*Step]) {
	coords := make([]*Step, 0)
	visited := make(map[common.Coord]int)
	for s.Size() > 0 {
		coords = append(coords, s.Pop())
	}
	for i := len(coords) - 1; i >= 0; i-- {
		visitIdx := len(coords) - 1 - i
		visited[coords[i].Coord] = visitIdx
		fmt.Printf("Step: %d\n%s\n\n", visitIdx, f.DebugPrint(visited))
	}
}

func (f Field) ExploreLongestPaths() int {
	maxLength := 0

	start := common.Coord{f.startX, f.startY}
	s := stack.New[*Step]()
	visited := make(map[common.Coord]int)
	visited[start] = 0
	s.Push(&Step{Coord: start, goTo: f.CanGoTo(start, visited), currentlyInIdx: -1})

	skipLastMulti := false

	for s.Size() > 0 {
		curr := s.Peek()
		if curr.X == f.endX && curr.Y == f.endY {
			if s.Size() > maxLength {
				maxLength = s.Size()
				fmt.Printf("Encountered length: %d, guess is %d\n", maxLength, maxLength-1)
			}
			delete(visited, curr.Coord)
			s.Pop()
			skipLastMulti = true
			continue
			//if maxLength == 167 {
			//	f.DebugPrintStack(*s.Copy())
			//}
			//fmt.Printf("Size: %d\n%s\n\n", maxLength, f.DebugPrint(visited))
		}
		//fmt.Printf("%s\n\n", f.DebugPrint(visited))

		if skipLastMulti || curr.AllVisited() {
			delete(visited, curr.Coord)
			s.Pop()
			if len(curr.goTo) > 1 {
				skipLastMulti = false
			}
			continue
		}
		curr.currentlyInIdx++

		visited[curr.goTo[curr.currentlyInIdx]] = s.Size()
		s.Push(&Step{
			Coord:          curr.goTo[curr.currentlyInIdx],
			currentlyInIdx: -1,
			goTo:           f.CanGoTo(curr.goTo[curr.currentlyInIdx], visited),
		})
	}
	return maxLength
}

func main() {
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t23-long-walk/1.txt")
	if err != nil {
		log.Fatalf("%w", err)
	}

	directions := common.NewDirections()
	field := ParseField(rows, directions)
	fmt.Printf("Max len: %d\n", field.ExploreLongestPaths()-1)
}
