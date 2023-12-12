package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"log"
	"strconv"
	"strings"
)

const (
	UNKNOWN     = '?'
	OPERATIONAL = '.'
	BROKEN      = '#'
)

type Springs struct {
	Field           []rune
	UnknownIndices  []int
	BrokenSequences []int
}

func unfoldSlice[T any](s []T, times int) []T {
	newSlice := make([]T, 0, len(s)*times)
	for i := 0; i < times; i++ {
		newSlice = append(newSlice, s...)
	}
	return newSlice
}

func unfoldTiles(tiles []rune, times int) []rune {
	newSlice := make([]rune, 0, len(tiles)*times+(times-1))
	for i := 0; i < times-1; i++ {
		newSlice = append(newSlice, tiles...)
		newSlice = append(newSlice, UNKNOWN)
	}
	newSlice = append(newSlice, tiles...)
	return newSlice
}

func SprintsFromString(raw string, unfolds int) Springs {
	// ????.######..#####. 1,6,5
	components := strings.Fields(raw)
	rawSequences := strings.Split(components[1], ",")
	sequences := lo.Map(rawSequences, func(item string, index int) int {
		num, err := strconv.Atoi(item)
		if err != nil {
			log.Fatalf("Failed to parse int: %s", item)
		}
		return num
	})
	sequences = unfoldSlice(sequences, unfolds)
	field := components[0]
	fieldChars := []rune(field)
	fieldChars = unfoldTiles(fieldChars, unfolds)
	unknownIndicesTuples := lo.Filter(lo.Map(fieldChars, func(c rune, index int) lo.Tuple2[rune, int] {
		return lo.T2(c, index)
	}), func(item lo.Tuple2[rune, int], index int) bool {
		return item.A == UNKNOWN
	})
	unknownIndices := lo.Map(unknownIndicesTuples, func(item lo.Tuple2[rune, int], index int) int {
		return item.B
	})
	return Springs{
		Field:           fieldChars,
		UnknownIndices:  unknownIndices,
		BrokenSequences: sequences,
	}
}

func ExistsBrokenRangeAt(tiles []rune, idx int, length int) bool {
	if idx+length > len(tiles) {
		return false
	}
	for i := 0; i < length; i++ {
		if tiles[idx+i] == OPERATIONAL {
			return false
		}
	}
	return idx+length == len(tiles) || tiles[idx+length] != BROKEN
}

func CountCombinations(tiles []rune, brokenSeqs []int, tilesIdx int, brokenSeqsIdx int, cache *map[lo.Tuple2[int, int]]int) int {

	solution, exists := (*cache)[lo.T2(tilesIdx, brokenSeqsIdx)]
	if exists {
		return solution
	}

	if tilesIdx >= len(tiles) {
		// Ended scan.
		if brokenSeqsIdx == len(brokenSeqs) {
			(*cache)[lo.T2(tilesIdx, brokenSeqsIdx)] = 1
			return 1
		} else {
			(*cache)[lo.T2(tilesIdx, brokenSeqsIdx)] = 0
			return 0
		}
	}

	switch tiles[tilesIdx] {
	case OPERATIONAL:
		count := CountCombinations(tiles, brokenSeqs, tilesIdx+1, brokenSeqsIdx, cache)
		(*cache)[lo.T2(tilesIdx, brokenSeqsIdx)] = count
		return count
	case BROKEN:
		count := computeBroken(tiles, brokenSeqs, tilesIdx, brokenSeqsIdx, cache)
		(*cache)[lo.T2(tilesIdx, brokenSeqsIdx)] = count
		return count
	case UNKNOWN:
		ifBrokenCount := computeBroken(tiles, brokenSeqs, tilesIdx, brokenSeqsIdx, cache)
		ifOperationalCount := CountCombinations(tiles, brokenSeqs, tilesIdx+1, brokenSeqsIdx, cache)
		(*cache)[lo.T2(tilesIdx, brokenSeqsIdx)] = ifOperationalCount + ifBrokenCount
		return ifOperationalCount + ifBrokenCount
	}
	log.Fatalf("Unreacheable")
	return 0
}

func computeBroken(tiles []rune, brokenSeqs []int, tilesIdx int, brokenSeqsIdx int, cache *map[lo.Tuple2[int, int]]int) int {
	if brokenSeqsIdx >= len(brokenSeqs) {
		return 0
	}

	if ExistsBrokenRangeAt(tiles, tilesIdx, brokenSeqs[brokenSeqsIdx]) {
		return CountCombinations(tiles, brokenSeqs, tilesIdx+brokenSeqs[brokenSeqsIdx]+1, brokenSeqsIdx+1, cache)
	}
	return 0
}

func main() {
	//rawSprings, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t12-hot-springs/test.txt")
	rawSprings, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t12-hot-springs/1.txt")
	if err != nil {
		log.Fatalf("Failed to read file: %w", err)
	}

	springs := lo.Map(rawSprings, common.NoIndex(func(s string) Springs {
		return SprintsFromString(s, 5)
	}))
	fmt.Printf("%d", lo.Sum(lo.Map(springs, func(s Springs, index int) int {
		cache := make(map[lo.Tuple2[int, int]]int)
		return CountCombinations(s.Field, s.BrokenSequences, 0, 0, &cache)
	})))
	//cache := make(map[lo.Tuple2[int, int]]int)
	//fmt.Printf("%d", CountCombinations(
	//	[]rune{UNKNOWN, UNKNOWN, UNKNOWN, OPERATIONAL, BROKEN, BROKEN, BROKEN},
	//	[]int{1, 1, 3},
	//	0, 0, &cache))
	//fmt.Printf("%d", CountCombinations(
	//	[]rune{UNKNOWN, UNKNOWN, UNKNOWN, OPERATIONAL},
	//	[]int{1, 1},
	//	0, 0, &cache))
}
