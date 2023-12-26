package main

import (
	"advent_of_code/common"
	"log"
)

func main() {
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t25-snowverload/test.txt")
	if err != nil {
		log.Fatalf("error: %w", err)
	}

}
