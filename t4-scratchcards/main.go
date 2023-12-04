package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/emirpasic/gods/utils"
	"github.com/samber/lo"
	"log"
	"slices"
	"strings"
)

type Card struct {
	Index          int
	WinningNumbers []int
	Numbers        []int
}

type CardWithCount struct {
	Card
	Count int
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

func (c Card) MatchingNumbers() int {
	total := 0
	for _, num := range c.Numbers {
		if lo.Contains(c.WinningNumbers, num) {
			total++
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

// Part 1

//func main() {
//	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t4-scratchcards/1.txt")
//	if err != nil {
//		log.Fatalf("Failed to read file: %w", err)
//	}
//	cards := lo.Map(rows, common.NoIndex(CardFromRow))
//	points := lo.Map(cards, common.NoIndex(Card.WorthPoints))
//	fmt.Printf("Total: %d", lo.Sum(points))
//}

func main() {
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t4-scratchcards/1.txt")
	if err != nil {
		log.Fatalf("Failed to read file: %w", err)
	}
	cards := lo.Map(rows, common.NoIndex(CardFromRow))

	//var indexToCard = make(map[int]Card)
	//for _, card := range cards {
	//	indexToCard[card.Index] = card
	//}

	slices.SortFunc(cards, func(a, b Card) int {
		return a.Index - b.Index
	})
	q := priorityqueue.NewWith(func(a, b interface{}) int {
		cardA := a.(Card)
		cardB := b.(Card)
		return utils.IntComparator(cardA.Index, cardB.Index)
	})

	indexToCardWithCount := make(map[int]*CardWithCount, len(cards))
	for _, c := range cards {
		indexToCardWithCount[c.Index] = &(CardWithCount{c, 1})
		q.Enqueue(c)
	}

	processed := 0

	for !q.Empty() {
		processedCardRaw, _ := q.Dequeue()
		processedCard := processedCardRaw.(Card)
		cardWithCount := indexToCardWithCount[processedCard.Index]
		delete(indexToCardWithCount, processedCard.Index)
		fmt.Printf("Processed card %d times %d\n", cardWithCount.Card.Index, cardWithCount.Count)
		processed += cardWithCount.Count

		matchingNumbersCount := processedCard.MatchingNumbers()
		for i := 1; i <= matchingNumbersCount; i++ {
			newIndex := processedCard.Index + i
			newCardWithIndex, ok := indexToCardWithCount[newIndex]
			if ok {
				newCardWithIndex.Count += cardWithCount.Count
			}
		}
	}
	fmt.Printf("%d\n", processed)

}
