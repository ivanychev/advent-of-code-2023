package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"log"
	"strings"
)

type GameRound struct {
	Blue  int
	Red   int
	Green int
}

type Game struct {
	GameIndex int
	Rounds    []GameRound
}

type Inventory struct {
	Blue  int
	Red   int
	Green int
}

func (i Inventory) Power() int64 {
	return int64(i.Red) * int64(i.Green) * int64(i.Blue)
}

func (g *Game) PossibleForInventory(i Inventory) bool {
	return lo.EveryBy(g.Rounds, func(r GameRound) bool {
		return r.Red <= i.Red && r.Blue <= i.Blue && r.Green <= i.Green
	})
}

func (g *Game) MinimumInventory() Inventory {
	return Inventory{
		Red: lo.Max(lo.Map(g.Rounds, func(r GameRound, index int) int {
			return r.Red
		})),
		Green: lo.Max(lo.Map(g.Rounds, func(r GameRound, index int) int {
			return r.Green
		})),
		Blue: lo.Max(lo.Map(g.Rounds, func(r GameRound, index int) int {
			return r.Blue
		})),
	}
}

func parseGameRound(rawRound string) GameRound {
	rawColorCounts := strings.Split(rawRound, ",")

	var blue, green, red, count int
	var color string
	for _, rawColorCount := range rawColorCounts {
		_, err := fmt.Sscanf(rawColorCount, "%d %s", &count, &color)
		if err != nil {
			log.Fatalf("Failed to parse color: %s, %w", rawColorCount, err)
		}

		switch color {
		case "red":
			red = count
		case "green":
			green = count
		case "blue":
			blue = count
		default:
			log.Fatalf("Invalid color parsed: %s", color)
		}
	}
	return GameRound{
		Blue:  blue,
		Red:   red,
		Green: green,
	}
}

func GameFromLine(line string) Game {
	// Example of input:
	// Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
	var gameIndex int
	_, err := fmt.Sscanf(line, "Game %d:", &gameIndex)
	if err != nil {
		log.Fatalf("Failed to parse game: %s %w", line, err)
	}
	rest := strings.TrimSpace(strings.Split(line, ":")[1])
	rawRounds := strings.Split(rest, ";")
	rounds := lo.Map(rawRounds, common.NoIndex(parseGameRound))
	return Game{
		GameIndex: gameIndex,
		Rounds:    rounds,
	}
}

// Part 1

//func main() {
//	//rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t2-cube/test.txt")
//	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t2-cube/1.txt")
//	inventory := Inventory{
//		Red:   12,
//		Green: 13,
//		Blue:  14,
//	}
//	if err != nil {
//		log.Fatalf("%w", err)
//	}
//	games := lo.Map(rows, common.NoIndex(GameFromLine))
//	possibleGames := lo.Filter(games, func(item Game, index int) bool {
//		return item.PossibleForInventory(inventory)
//	})
//	fmt.Printf("%d", lo.SumBy(possibleGames, func(g Game) int {
//		return g.GameIndex
//	}))
//}

// Part 2
func main() {
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t2-cube/1.txt")
	//rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t2-cube/test.txt")
	if err != nil {
		log.Fatalf("%w", err)
	}
	games := lo.Map(rows, common.NoIndex(GameFromLine))
	powers := lo.Map(games, func(g Game, index int) int64 {
		return g.MinimumInventory().Power()
	})
	println(lo.Reduce(powers, func(agg int64, item int64, index int) int64 {
		return agg + item
	}, int64(0)))
}
