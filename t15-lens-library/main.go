package main

import (
	"fmt"
	"github.com/samber/lo"
	"io"
	"log"
	"os"
	"strconv"
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

type Item struct {
	key   string
	value int
}

type HashMap struct {
	lenses [256][]Item
}

func NewHashMap() HashMap {
	var lenses [256][]Item
	for idx, _ := range lenses {
		lenses[idx] = make([]Item, 0)
	}
	return HashMap{lenses}
}

func IndexOfWord(items []Item, word string) (int, bool) {
	for i := 0; i < len(items); i++ {
		if items[i].key == word {
			return i, true
		}
	}
	return 0, false
}

func (m *HashMap) ProcessCommand(cmd string) {
	// rn=1,cm-
	if strings.Contains(cmd, "=") {
		components := strings.Split(cmd, "=")
		word := components[0]
		focalLength, _ := strconv.Atoi(components[1])
		hash := HashString(word)
		index, found := IndexOfWord(m.lenses[hash], word)
		if found {
			m.lenses[hash][index].value = focalLength
		} else {
			m.lenses[hash] = append(m.lenses[hash], Item{word, focalLength})
		}
	} else if strings.Contains(cmd, "-") {
		components := strings.Split(cmd, "-")
		word := components[0]
		hash := HashString(word)
		index, found := IndexOfWord(m.lenses[hash], word)
		if found {
			m.lenses[hash] = append(m.lenses[hash][:index], m.lenses[hash][index+1:]...)
		}
	} else {
		log.Fatalf("Invalid command")
	}
}

func (m HashMap) TotalPower() int {
	total := 0
	for idx, cell := range m.lenses {
		for boxIdx, item := range cell {
			total += (idx + 1) * (boxIdx + 1) * item.value
		}
	}
	return total
}

func main() {
	file, _ := os.Open("/Users/iv/Code/advent-of-code-2023/t15-lens-library/1.txt")
	//file, _ := os.Open("/Users/iv/Code/advent-of-code-2023/t15-lens-library/test.txt")
	rawChars, _ := io.ReadAll(file)
	steps := strings.Split(strings.TrimSpace(string(rawChars)), ",")
	hm := NewHashMap()
	for _, step := range steps {
		hm.ProcessCommand(step)
	}
	fmt.Printf("Total: %d\n", hm.TotalPower())
}
