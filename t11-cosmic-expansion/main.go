package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"log"
)

type Tile rune

const (
	EMPTY_SPACE      Tile = '*'
	GALAXY           Tile = '#'
	EXPANSION_FACTOR      = 2
)

type Coord struct {
	x, y int
}

func LenX(tiles [][]Tile) int {
	return len(tiles[0])
}

func LenY(tiles [][]Tile) int {
	return len(tiles)
}

func XAllSpace(tiles [][]Tile, x int) bool {
	for y := 0; y < LenY(tiles); y++ {
		if tiles[y][x] == GALAXY {
			return false
		}
	}
	return true
}

func YAllSpace(tiles [][]Tile, y int) bool {
	for x := 0; x < LenX(tiles); x++ {
		if tiles[y][x] == GALAXY {
			return false
		}
	}
	return true
}

func GalaxyCoords(tiles [][]Tile) []Coord {
	galaxies := make([]Coord, 0)
	for y := 0; y < LenY(tiles); y++ {
		for x := 0; x < LenX(tiles); x++ {
			if tiles[y][x] == GALAXY {
				galaxies = append(galaxies, Coord{x, y})
			}
		}
	}
	return galaxies
}

func absInt(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}

func manhattanDistance(first, second Coord, xCumDistances, yCumDistances []int) int {
	return absInt(xCumDistances[first.x]-xCumDistances[second.x]) + absInt(yCumDistances[first.y]-yCumDistances[second.y])
}

func main() {
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t11-cosmic-expansion/1.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %w", err)
	}

	tiles := lo.Map(rows, common.NoIndex(func(row string) []Tile {
		return []Tile(row)
	}))
	galaxyCoords := GalaxyCoords(tiles)
	xCumDistances := lo.Reduce(lo.Range(LenX(tiles)), func(agg []int, x int, index int) []int {
		var prev int
		if len(agg) > 0 {
			prev = agg[len(agg)-1]
		}
		if XAllSpace(tiles, x) {
			agg = append(agg, prev+EXPANSION_FACTOR)
		} else {
			agg = append(agg, prev+1)
		}
		return agg
	}, make([]int, 0, LenX(tiles)))
	yCumDistances := lo.Reduce(lo.Range(LenY(tiles)), func(agg []int, y int, index int) []int {
		var prev int
		if len(agg) > 0 {
			prev = agg[len(agg)-1]
		}
		if YAllSpace(tiles, y) {
			agg = append(agg, prev+EXPANSION_FACTOR)
		} else {
			agg = append(agg, prev+1)
		}
		return agg
	}, make([]int, 0, LenY(tiles)))

	total := 0
	for i := 0; i < len(galaxyCoords)-1; i++ {
		for j := i + 1; j < len(galaxyCoords); j++ {
			distance := manhattanDistance(galaxyCoords[i], galaxyCoords[j], xCumDistances, yCumDistances)
			total += distance
		}
	}
	fmt.Printf("Distance: %d\n", total)
}
