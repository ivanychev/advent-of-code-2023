package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"log"
)

func AllZeros(values []int) bool {
	return lo.EveryBy(values, func(item int) bool {
		return item == 0
	})
}

func Differentiate(values []int) []int {
	diff := make([]int, 0, len(values)-1)
	for i := 0; i < len(values)-1; i++ {
		diff = append(diff, values[i+1]-values[i])
	}
	return diff
}

func (h History) ExtrapolateRight() int {
	diffs := [][]int{h.values}
	for !AllZeros(diffs[len(diffs)-1]) {
		diffs = append(diffs, Differentiate(diffs[len(diffs)-1]))
	}

	extrapolatingIndex := len(diffs) - 1
	extrapolatedValue := 0
	for extrapolatingIndex > 0 {
		newExtrapolatedValue := extrapolatedValue + diffs[extrapolatingIndex-1][len(diffs[extrapolatingIndex-1])-1]
		extrapolatedValue = newExtrapolatedValue
		extrapolatingIndex--
	}
	return extrapolatedValue
}

func (h History) ExtrapolateLeft() int {
	diffs := [][]int{h.values}
	for !AllZeros(diffs[len(diffs)-1]) {
		diffs = append(diffs, Differentiate(diffs[len(diffs)-1]))
	}

	extrapolatingIndex := len(diffs) - 1
	extrapolatedValue := 0
	for extrapolatingIndex > 0 {
		newExtrapolatedValue := diffs[extrapolatingIndex-1][0] - extrapolatedValue
		extrapolatedValue = newExtrapolatedValue
		extrapolatingIndex--
	}
	return extrapolatedValue
}

type History struct {
	values []int
}

func HistoryFromString(s string) History {
	return History{
		common.StringOfNumbersToInts(s),
	}
}

// Part 1
//func main() {
//	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t9-migrate/1.txt")
//	if err != nil {
//		log.Fatalf("Failed to open file %w", err)
//	}
//	histories := lo.Map(rows, common.NoIndex(HistoryFromString))
//	extrapolated := lo.Map(histories, func(h History, index int) int {
//		return h.ExtrapolateRight()
//	})
//	fmt.Printf("Sum: %d\n", lo.Sum(extrapolated))
//}

// Part 2
func main() {
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t9-migrate/1.txt")
	if err != nil {
		log.Fatalf("Failed to open file %w", err)
	}
	histories := lo.Map(rows, common.NoIndex(HistoryFromString))
	extrapolated := lo.Map(histories, func(h History, index int) int {
		return h.ExtrapolateLeft()
	})
	fmt.Printf("Sum: %d\n", lo.Sum(extrapolated))
}
