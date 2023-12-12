package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"gonum.org/v1/gonum/stat/combin"
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
	field             []rune
	unknownIndices    []int
	notWorksSequences []int
}

func (s Springs) BrokenSpringsCount() int {
	return lo.Sum(s.notWorksSequences)
}

func (s Springs) ObservedBrokenSpringsCount() int {
	return lo.Count(s.field, BROKEN)
}

func (s Springs) UnknownSpringsCount() int {
	return len(s.unknownIndices)
}

func testField(field []rune, notWorksSequences []int) bool {
	currentSequenceSize := 0
	expectedSequenceSizeIdx := 0
	for idx := 0; idx <= len(field); idx++ {
		// Ending always with "."
		cell := OPERATIONAL
		if idx < len(field) {
			cell = field[idx]
		}

		switch cell {
		case OPERATIONAL:
			if currentSequenceSize == 0 {
				continue
			}

			if expectedSequenceSizeIdx >= len(notWorksSequences) || notWorksSequences[expectedSequenceSizeIdx] != currentSequenceSize {
				return false
			}
			expectedSequenceSizeIdx += 1
			currentSequenceSize = 0
		case BROKEN:
			currentSequenceSize += 1
		default:
			log.Fatalf("Invalid char: %v", cell)
		}
	}
	return true
}

func (s Springs) CountCombinations() int {
	n := s.UnknownSpringsCount()
	k := s.BrokenSpringsCount() - s.ObservedBrokenSpringsCount()
	indices := make([]int, k)
	gen := combin.NewCombinationGenerator(n, k)

	attemptField := make([]rune, len(s.field))
	copy(attemptField, s.field)

	counter := 0
	for gen.Next() {
		gen.Combination(indices)

		for _, i := range s.unknownIndices {
			attemptField[i] = OPERATIONAL
		}
		for _, i := range indices {
			attemptField[s.unknownIndices[i]] = BROKEN
		}
		if testField(attemptField, s.notWorksSequences) {
			counter++
		}
	}
	return counter
}

func SprintsFromString(raw string) Springs {
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
	field := components[0]
	fieldChars := []rune(field)
	unknownIndicesTuples := lo.Filter(lo.Map(fieldChars, func(c rune, index int) lo.Tuple2[rune, int] {
		return lo.T2(c, index)
	}), func(item lo.Tuple2[rune, int], index int) bool {
		return item.A == UNKNOWN
	})
	unknownIndices := lo.Map(unknownIndicesTuples, func(item lo.Tuple2[rune, int], index int) int {
		return item.B
	})
	return Springs{
		field:             fieldChars,
		unknownIndices:    unknownIndices,
		notWorksSequences: sequences,
	}
}

func main() {
	rawSprings, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t12-hot-springs/1.txt")
	if err != nil {
		log.Fatalf("Failed to read file: %w", err)
	}

	springs := lo.Map(rawSprings, common.NoIndex(SprintsFromString))
	fmt.Printf("%d", lo.Sum(lo.Map(springs, func(s Springs, index int) int {
		return s.CountCombinations()
	})))
}
