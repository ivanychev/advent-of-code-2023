package common

import (
	"bufio"
	"fmt"
	"os"
)

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
