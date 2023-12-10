package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"github.com/zyedidia/generic/queue"
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

func (c Coord) To(d Direction) Coord {
	switch d {
	case UP:
		return Coord{c.x, c.y - 1}
	case DOWN:
		return Coord{c.x, c.y + 1}
	case LEFT:
		return Coord{c.x - 1, c.y}
	case RIGHT:
		return Coord{c.x + 1, c.y}
	}
	log.Fatalf("")
	return Coord{}
}

func (c Coord) Neighbours(lenX, lenY int) []Coord {
	coords := make([]Coord, 0)
	if 0 <= c.x-1 && c.x-1 < lenX {
		coords = append(coords, Coord{c.x - 1, c.y})
	}
	if 0 <= c.x+1 && c.x+1 < lenX {
		coords = append(coords, Coord{c.x + 1, c.y})
	}
	if 0 <= c.y-1 && c.y-1 < lenY {
		coords = append(coords, Coord{c.x, c.y - 1})
	}
	if 0 <= c.y+1 && c.y+1 < lenY {
		coords = append(coords, Coord{c.x, c.y + 1})
	}
	return coords
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

func tileAndMoveDirectionToInnerCandidates(tile rune, d Direction) []Direction {
	switch tile {
	case UD:
		if d == UP {
			return []Direction{RIGHT}
		} else if d == DOWN {
			return []Direction{LEFT}
		} else {
			log.Fatalf("Failed switch at %v", UD)
		}
	case LR:
		if d == RIGHT {
			return []Direction{DOWN}
		} else if d == LEFT {
			return []Direction{UP}
		} else {
			log.Fatalf("Failed switch at %v", LR)
		}
	case DR:
		if d == UP {
			return []Direction{}
		} else if d == LEFT {
			return []Direction{UP, LEFT}
		} else {
			log.Fatalf("Failed switch at %v", DR)
		}
	case DL:
		if d == UP {
			return []Direction{UP, RIGHT}
		} else if d == RIGHT {
			return []Direction{}
		} else {
			log.Fatalf("Failed switch at %v", DL)
		}
	case UR:
		if d == DOWN {
			return []Direction{LEFT, DOWN}
		} else if d == LEFT {
			return []Direction{}
		} else {
			log.Fatalf("Failed switch at %v", UR)
		}
	case UL:
		if d == DOWN {
			return []Direction{}
		} else if d == RIGHT {
			return []Direction{RIGHT, DOWN}
		} else {
			log.Fatalf("Failed switch at %v", UL)
		}
	}
	log.Fatalf("Broken switch")
	return []Direction{}
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

func (f Field) LenX() int {
	return len(f.tiles[0])
}
func (f Field) LenY() int {
	return len(f.tiles)
}

func (f Field) Move(from Coord, notTo *Direction, startMoveTo Direction) (Coord, *Direction, *Direction) {
	directions := lo.Filter(tileToDirections(f.tiles[from.y][from.x]), func(item Direction, index int) bool {
		return notTo == nil || item != *notTo
	})
	if notTo == nil && len(directions) != 2 {
		log.Fatalf("Unexpected while nil %d", len(directions))
	}
	if notTo != nil && len(directions) != 1 {
		log.Fatalf("Unexpected while not nil %d", len(directions))
	}
	var directionToMove Direction
	if notTo == nil {
		directionToMove = startMoveTo
	} else {
		directionToMove = directions[0]
	}
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
		log.Fatalf("Unexpected direction %v", directionToMove)
	}
	opposite := getOppositeDirection(directionToMove)
	return from, &directionToMove, &opposite
}

func scanCycle(startingCoord Coord, field Field, startMoveTo Direction, callback func(c Coord, movedTo *Direction)) {
	moved := false
	slowerPointer := startingCoord
	fasterPointer := startingCoord
	var lastFastOppositeDirection, movedTo, lastSlowOppositeDirection *Direction
	slowerMoved := 0
	callback(fasterPointer, nil)
	for !moved || fasterPointer != startingCoord {
		slowerPointer, _, lastSlowOppositeDirection = field.Move(slowerPointer, lastSlowOppositeDirection, startMoveTo)
		fasterPointer, movedTo, lastFastOppositeDirection = field.Move(fasterPointer, lastFastOppositeDirection, startMoveTo)
		callback(fasterPointer, movedTo)
		fasterPointer, movedTo, lastFastOppositeDirection = field.Move(fasterPointer, lastFastOppositeDirection, startMoveTo)
		callback(fasterPointer, movedTo)
		slowerMoved++
		moved = true
	}
}

func scanCycleTiles(startingCoord Coord, field Field, startMoveTo Direction) map[Coord]struct{} {
	cycleTiles := make(map[Coord]struct{})
	scanCycle(startingCoord, field, startMoveTo, func(c Coord, movedTo *Direction) {
		cycleTiles[c] = struct{}{}
	})
	return cycleTiles
}

func findInnerTiles(startingCoord Coord, field Field, cycleTileSet map[Coord]struct{}, startMoveTo Direction) map[Coord]struct{} {
	innerTiles := make(map[Coord]struct{})
	scanCycle(startingCoord, field, startMoveTo, func(c Coord, movedTo *Direction) {
		if movedTo == nil {
			return
		}
		candidates := lo.Map(
			tileAndMoveDirectionToInnerCandidates(field.tiles[c.y][c.x], *movedTo),
			func(d Direction, index int) Coord {
				return c.To(d)
			})

		for _, candidate := range candidates {
			_, partOfCycle := cycleTileSet[candidate]
			if _, exists := cycleTileSet[candidate]; !exists && !partOfCycle {
				innerTiles[candidate] = struct{}{}
			}
		}
	})
	return innerTiles
}

func bfsInner(field Field, starting Coord, visited *map[Coord]struct{}, predicate func(c Coord) bool, visitor func(c Coord)) {
	q := queue.New[Coord]()
	q.Enqueue(starting)
	for !q.Empty() {
		current := q.Dequeue()
		if _, already := (*visited)[current]; already {
			continue
		}
		visitor(current)
		(*visited)[current] = struct{}{}
		neighbours := lo.Filter(current.Neighbours(field.LenX(), field.LenY()), common.NoIndex(predicate))
		for _, n := range neighbours {
			q.Enqueue(n)
		}
	}
}

func bfs(field Field, startingCoords []Coord, predicate func(c Coord) bool, visitor func(c Coord)) {
	visited := make(map[Coord]struct{})
	for _, startingCoord := range startingCoords {
		if _, exists := visited[startingCoord]; exists {
			continue
		}
		bfsInner(field, startingCoord, &visited, predicate, visitor)
	}
}

func main() {
	//rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t10-pipes/test.txt")
	//startIsRune := DR
	//rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t10-pipes/2test1.txt")
	//startIsRune := DR
	//moveTo := RIGHT
	//rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t10-pipes/2test2.txt")
	//startIsRune := DR
	//moveTo := DOWN
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t10-pipes/1.txt")
	startIsRune := UD
	moveTo := UP
	if err != nil {
		log.Fatalf("Failed to open file: %w", err)
	}
	tiles := lo.Map(rows, common.NoIndex(func(s string) []rune {
		return []rune(s)
	}))
	startingCoord := findStartingCoords(tiles)
	tiles[startingCoord.y][startingCoord.x] = startIsRune
	field := Field{tiles: tiles, startingCoord: startingCoord}

	tileSet := scanCycleTiles(startingCoord, field, moveTo)
	startingCoords := make([]Coord, 0)
	for x := 0; x < field.LenX(); x++ {
		startingCoords = append(startingCoords, Coord{x, 0})
	}
	for y := 1; y < field.LenX(); y++ {
		startingCoords = append(startingCoords, Coord{0, y})
	}
	inLoop := 0
	for _, start := range startingCoords {
		delta := 0
		intersections := 0
		for start.x+delta < field.LenX() && start.y+delta < field.LenY() {
			coord := Coord{start.x + delta, start.y + delta}
			if _, partOfCycle := tileSet[coord]; partOfCycle {
				tile := field.tiles[coord.y][coord.x]
				if tile != UR && tile != DL {
					intersections++
				}
			} else if intersections%2 == 1 {
				//fmt.Printf("%+v\n", coord)
				inLoop += 1
			}
			delta++
		}
	}
	fmt.Printf("%d\n", inLoop)
}
