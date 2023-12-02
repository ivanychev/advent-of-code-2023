package main

import (
	"advent_of_code/common"
	"bufio"
	"fmt"
	"github.com/dlclark/regexp2"
	"log"
	"os"
)

func ReadPuzzle(path string) []string {
	readFile, err := os.Open(path)

	if err != nil {
		log.Fatalf("Failed to read the file: %w", err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	return lines
}

func PuzzleToNumbers(puzzle []string, extractor DigitExtractor) [][]int {
	var numbers = make([][]int, 0, len(puzzle))

	for _, p := range puzzle {
		numbers = append(numbers, extractor.Extract(p))
	}
	return numbers
}

func GetCalibrationNumbers(allNumbers [][]int) []int {
	calNumbers := make([]int, 0, len(allNumbers))
	for _, numbersInRow := range allNumbers {
		calNumbers = append(calNumbers, numbersInRow[0]*10+numbersInRow[len(numbersInRow)-1])
	}
	return calNumbers
}

func WriteItemsToFile[T any](path string, items []T) {
	writeFile, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to open file for writing: %w", err)
	}
	for _, item := range items {
		fmt.Fprintf(writeFile, "%v\n", item)
	}
}

// Part 1

//func main() {
//	puzzle := ReadPuzzle("/Users/iv/Code/advent-of-code-2023/t1-trebuchet/1.txt")
//	numbers := PuzzleToNumbers(puzzle, &TrivialExtractor{re: regexp2.MustCompile("\\d", regexp2.IgnoreCase)})
//	calNumbers := GetCalibrationNumbers(numbers)
//	WriteItemsToFile("/Users/iv/Code/advent-of-code-2023/t1-trebuchet/1-output.txt", []int{sum(calNumbers)})
//}

// Part 2
func main() {
	puzzle := ReadPuzzle("/Users/iv/Code/advent-of-code-2023/t1-trebuchet/1.txt")
	numbers := PuzzleToNumbers(puzzle, &ComplexExtractor{
		re: regexp2.MustCompile("(?=(\\d|one|two|three|four|five|six|seven|eight|nine))", regexp2.IgnoreCase),
		strToDigit: map[string]int{
			"1":     1,
			"2":     2,
			"3":     3,
			"4":     4,
			"5":     5,
			"6":     6,
			"7":     7,
			"8":     8,
			"9":     9,
			"one":   1,
			"two":   2,
			"three": 3,
			"four":  4,
			"five":  5,
			"six":   6,
			"seven": 7,
			"eight": 8,
			"nine":  9,
		},
	})
	calNumbers := GetCalibrationNumbers(numbers)
	WriteItemsToFile("/Users/iv/Code/advent-of-code-2023/t1-trebuchet/2-output.txt", []int{common.Sum(calNumbers)})
}
