package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"io"
	"os"
	"strings"
)

func HashString(s string) int {
	return HashRunes([]rune(s))
}

func HashRunes(s []rune) int {
	return lo.Reduce(s, func(agg int, item rune, index int) int {
		return 17 * (agg + int(item)) % 256
	}, 0)
}

func main() {
	file, _ := os.Open("/Users/iv/Code/advent-of-code-2023/t15-lens-library/1.txt")
	rawChars, _ := io.ReadAll(file)
	steps := strings.Split(strings.TrimSpace(string(rawChars)), ",")
	stepHashes := lo.Map(steps, common.NoIndex(HashString))

	fmt.Printf("Total: %d\n", lo.Sum(stepHashes))
}
