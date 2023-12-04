package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"log"
	"strings"
)

type Card struct {
	Index          int
	WinningNumbers []int
	Numbers        []int
}

func (c Card) WorthPoints() int {
	total := 0
	for _, num := range c.Numbers {
		if lo.Contains(c.WinningNumbers, num) {
			if total == 0 {
				total = 1
			} else {
				total *= 2
			}
		}
	}
	return total
}

func CardFromRow(row string) Card {
	rawCard, rawNumbers, _ := strings.Cut(row, ": ")
	var cardIndex int
	count, err := fmt.Sscanf(rawCard, "Card %d", &cardIndex)
	if err != nil || count == 0 {
		log.Fatalf("Failed to parse Card index: %s", rawCard)
	}
	rawWinningNumbers, rawNumbers, _ := strings.Cut(rawNumbers, "|")

	winningNumbers := common.StringOfNumbersToInts(strings.TrimSpace(rawWinningNumbers))
	numbers := common.StringOfNumbersToInts(strings.TrimSpace(rawNumbers))

	return Card{
		Index:          cardIndex,
		WinningNumbers: winningNumbers,
		Numbers:        numbers,
	}
}

func main() {
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t4-scratchcards/1.txt")
	if err != nil {
		log.Fatalf("Failed to read file: %w", err)
	}
	cards := lo.Map(rows, common.NoIndex(CardFromRow))
	points := lo.Map(cards, common.NoIndex(Card.WorthPoints))
	fmt.Printf("Total: %d", lo.Sum(points))
}
