package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/emirpasic/gods/utils"
	"github.com/samber/lo"
	"log"
	"slices"
	"strconv"
	"strings"
)

type Card int8

const (
	C_A Card = 'M'
	C_K Card = 'L'
	C_Q Card = 'K'
	C_J Card = 'J'
	C_T Card = 'I'
	C_9 Card = 'H'
	C_8 Card = 'G'
	C_7 Card = 'F'
	C_6 Card = 'E'
	C_5 Card = 'D'
	C_4 Card = 'C'
	C_3 Card = 'B'
	C_2 Card = 'A'
)

func RuneToCard(r rune) Card {
	var c Card
	switch r {
	case 'A':
		c = C_A
	case 'K':
		c = C_K
	case 'Q':
		c = C_Q
	case 'J':
		c = C_J
	case 'T':
		c = C_T
	case '9':
		c = C_9
	case '8':
		c = C_8
	case '7':
		c = C_7
	case '6':
		c = C_6
	case '5':
		c = C_5
	case '4':
		c = C_4
	case '3':
		c = C_3
	case '2':
		c = C_2
	}
	return c
}

type Hand struct {
	RawCards       string
	Cards          []Card
	CardsCombined  string
	CardCounts     map[Card]int
	CountsCombined string
	Bid            int
}

func ParseHand(rawHand string) Hand {
	rawComponents := strings.Fields(rawHand)
	rawCards := []rune(rawComponents[0])
	rawBid := rawComponents[1]
	bid, _ := strconv.Atoi(rawBid)
	cards := lo.Map(rawCards, common.NoIndex(RuneToCard))
	cardsCombined := strings.Join(lo.Map(cards, func(c Card, index int) string {
		return fmt.Sprintf("%c", c)
	}), "")

	cardCounts := lo.CountValues(cards)
	counts := lo.Values(cardCounts)
	slices.SortFunc(counts, func(a, b int) int {
		return b - a
	})
	countsCombined := strings.Join(lo.Map(counts, func(c int, index int) string {
		return fmt.Sprintf("%d", c)
	}), "")
	return Hand{RawCards: rawComponents[0], Cards: cards, CardCounts: cardCounts, Bid: bid,
		CountsCombined: countsCombined, CardsCombined: cardsCombined}
}

func CompareHands(a, b Hand) int {
	compareHandCounts := utils.StringComparator(a.CountsCombined, b.CountsCombined)
	if compareHandCounts != 0 {
		return compareHandCounts
	}
	return utils.StringComparator(a.CardsCombined, b.CardsCombined)
}

func main() {
	rawBids, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t7-camel/1.txt")
	if err != nil {
		log.Fatalf("Failed to read file %w", err)
	}
	hands := lo.Map(rawBids, common.NoIndex(ParseHand))
	slices.SortFunc(hands, CompareHands)
	winnings := lo.Map(hands, func(h Hand, index int) int64 {
		return int64(index+1) * int64(h.Bid)
	})

	fmt.Printf("%+v", lo.Sum(winnings))
}
