package main

import (
	"fmt"
	"github.com/samber/lo"
	"io"
	"log"
	"os"
	"strings"
)

func Gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func Lcm(a, b int64) int64 {
	return a * b / Gcd(a, b)
}

const LAST = 2

type Edge struct {
	name  string
	left  *Edge
	right *Edge
}

func ensureEdgeExists(name string, nameToEdge *map[string]*Edge) *Edge {
	edge, ok := (*nameToEdge)[name]
	if !ok {
		edge = &Edge{
			name:  name,
			right: nil,
			left:  nil,
		}
		(*nameToEdge)[name] = edge
	}
	return edge
}

func EdgeFromString(s string, nameToEdge *map[string]*Edge, replacer *strings.Replacer) *Edge {
	components := strings.Fields(replacer.Replace(s))
	name := components[0]
	left := components[1]
	right := components[2]

	currentEdge := ensureEdgeExists(name, nameToEdge)
	leftEdge := ensureEdgeExists(left, nameToEdge)
	rightEdge := ensureEdgeExists(right, nameToEdge)
	currentEdge.left = leftEdge
	currentEdge.right = rightEdge
	return currentEdge
}

func createEdgesAndCommands(path string) (map[string]*Edge, []byte) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open file %w", err)
	}
	contents, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file %w", err)
	}
	components := strings.Split(string(contents), "\n\n")
	commands := []byte(strings.TrimSpace(components[0]))
	rawEdges := strings.Split(components[1], "\n")
	replacer := strings.NewReplacer("=", " ", "(", " ", ")", " ", ",", " ")
	nameToEdge := make(map[string]*Edge)
	for _, edge := range rawEdges {
		EdgeFromString(edge, &nameToEdge, replacer)
	}
	return nameToEdge, commands
}

//func main() {
//	nameToEdge, commands := createEdgesAndCommands("/Users/iv/Code/advent-of-code-2023/t8-haunted/1.txt")
//
//	currentStep := -1
//	totalSteps := 0
//	currentEdge := nameToEdge["AAA"]
//	for {
//		currentStep = (currentStep + 1) % len(commands)
//		totalSteps += 1
//		switch commands[currentStep] {
//		case 'L':
//			currentEdge = currentEdge.left
//		case 'R':
//			currentEdge = currentEdge.right
//		default:
//			fmt.Errorf("Invalid command %c", commands[currentStep])
//		}
//		if currentEdge.name == "ZZZ" {
//			break
//		}
//	}
//	fmt.Printf("%d", totalSteps)
//}

func allEnding(edges []*Edge) bool {
	for _, e := range edges {
		if e.name[LAST] != 'Z' {
			return false
		}
	}
	return true
}

func stepsToReach(currentEdge *Edge, commands []byte) int {
	currentStep := -1
	totalSteps := 0
	for {
		currentStep = (currentStep + 1) % len(commands)
		totalSteps += 1
		switch commands[currentStep] {
		case 'L':
			currentEdge = currentEdge.left
		case 'R':
			currentEdge = currentEdge.right
		default:
			fmt.Errorf("Invalid command %c", commands[currentStep])
		}
		if currentEdge.name == "ZZZ" {
			break
		}
	}
	return totalSteps
}

func stepsToReachAll(currentEdges []*Edge, commands []byte) int {
	currentStep := -1
	totalSteps := 0
	for {
		currentStep = (currentStep + 1) % len(commands)
		totalSteps += 1
		switch commands[currentStep] {
		case 'L':
			currentEdges = lo.Map(currentEdges, func(e *Edge, index int) *Edge {
				return e.left
			})
		case 'R':
			currentEdges = lo.Map(currentEdges, func(e *Edge, index int) *Edge {
				return e.right
			})
		default:
			fmt.Errorf("Invalid command %c", commands[currentStep])
		}
		if allEnding(currentEdges) {
			break
		}
		if totalSteps%30000000 == 0 {
			fmt.Printf("Done %d iterations\n", totalSteps)
		}
	}
	return totalSteps
}

type EdgeStep struct {
	e              *Edge
	commandPointer int64
}

type CycleInfo struct {
	firstEncounterSteps  int64
	secondEncounterSteps int64
	commandPointer       int64
	endsEncounteredAt    []int64
	edge                 *Edge
}

type CyclePosition struct {
	info                       CycleInfo
	length                     int64
	coord                      int64
	nextEndIndex               int
	endsEncounteredInCycle     []int64
	endsEncounteredBeforeCycle []int64
}

func (p CyclePosition) UntilNext() int64 {
	if len(p.endsEncounteredBeforeCycle) > 0 {
		fmt.Errorf("Not implemented")
	}
	if p.coord < 0 {
		return p.endsEncounteredInCycle[p.nextEndIndex] - p.coord
	}
	steps := (p.endsEncounteredInCycle[p.nextEndIndex] + p.length - p.coord) % p.length
	if steps == 0 {
		steps = p.length
	}
	return steps
}

func (p *CyclePosition) StepBy(steps int64) {
	steps = steps % p.length
	endCoord := p.coord + steps
	if endCoord > 0 {
		endCoord = (p.coord + steps) % p.length
	}
	offset := int64(0)
	if p.endsEncounteredInCycle[p.nextEndIndex] <= p.coord {
		offset = p.length
	}
	newIndex := p.nextEndIndex
	for offset+p.endsEncounteredInCycle[newIndex] <= steps {
		if newIndex+1 == len(p.endsEncounteredInCycle) {
			offset += p.length
			newIndex = 0
		} else {
			newIndex++
		}
	}

	p.coord = endCoord
	p.nextEndIndex = newIndex
}

func prevIndex(index int64, length int64) int64 {
	prev := index - 1
	if prev == -1 {
		prev = length - 1
	}
	return prev
}

func nextIndex(index int64, length int64) int64 {
	next := index + 1
	if next == length {
		next = 0
	}
	return next
}

//func (p CyclePosition) AtEnd() bool {
//	return p.endsEncounteredBeforeCycle[p.nextEndIndex] == p.coord || p.endsEncounteredBeforeCycle[prevIndex(p.nextEndIndex, p.length)] == p.coord
//}

func findCycle(e *Edge, commands []byte) CycleInfo {
	commandPointer := int64(0)
	totalSteps := int64(0)
	endsEncounteredAt := make([]int64, 0)
	edgesSet := make(map[EdgeStep]int64)
	edgesSet[EdgeStep{e, totalSteps}] = int64(0)
	for {
		if e.name[LAST] == 'Z' {
			endsEncounteredAt = append(endsEncounteredAt, totalSteps)
		}
		if commands[commandPointer] == 'L' {
			e = e.left
		} else {
			e = e.right
		}
		commandPointer = (commandPointer + 1) % int64(len(commands))
		totalSteps++
		es := EdgeStep{e, commandPointer}
		if firstEncounterSteps, exists := edgesSet[es]; exists {
			return CycleInfo{
				firstEncounterSteps:  firstEncounterSteps,
				secondEncounterSteps: totalSteps,
				commandPointer:       commandPointer,
				endsEncounteredAt:    endsEncounteredAt,
				edge:                 e,
			}
		} else {
			edgesSet[es] = totalSteps
		}
	}
}

func main() {
	nameToEdge, commands := createEdgesAndCommands("/Users/iv/Code/advent-of-code-2023/t8-haunted/1.txt")
	//nameToEdge, commands := createEdgesAndCommands("/Users/iv/Code/advent-of-code-2023/t8-haunted/1.txt")

	edges := lo.Filter(lo.Values(nameToEdge), func(item *Edge, index int) bool {
		return item.name[LAST] == 'A'
	})
	cycleInfos := lo.Map(edges, func(e *Edge, index int) CycleInfo {
		return findCycle(e, commands)
	})
	cyclePositions := lo.Map(cycleInfos, func(c CycleInfo, index int) CyclePosition {
		return CyclePosition{
			info:         c,
			coord:        -c.firstEncounterSteps,
			length:       c.secondEncounterSteps - c.firstEncounterSteps,
			nextEndIndex: 0,
			endsEncounteredBeforeCycle: lo.Filter(lo.Map(c.endsEncounteredAt, func(item int64, index int) int64 {
				return item - c.firstEncounterSteps
			}), func(item int64, index int) bool {
				return item < 0
			}),
			endsEncounteredInCycle: lo.Filter(lo.Map(c.endsEncounteredAt, func(item int64, index int) int64 {
				return item - c.firstEncounterSteps
			}), func(item int64, index int) bool {
				return item >= 0
			}),
		}
	})

	res := lo.Reduce(cyclePositions, func(agg int64, p CyclePosition, index int) int64 {
		var length = p.info.secondEncounterSteps - p.info.firstEncounterSteps
		return Lcm(agg, length)
	}, 1)
	fmt.Printf("%d\n", res)

	//for _, c := range cyclePositions {
	//	fmt.Printf("%+v\n", c)
	//}

	//p := CyclePosition{
	//	info: CycleInfo{
	//		firstEncounterSteps:  0,
	//		secondEncounterSteps: 5,
	//		commandPointer:       0,
	//		endsEncounteredAt:    []int64{1, 4},
	//		edge:                 &Edge{"fdf", nil, nil},
	//	},
	//	length:                     5,
	//	coord:                      0,
	//	nextEndIndex:               0,
	//	endsEncounteredInCycle:     []int64{1, 4},
	//	endsEncounteredBeforeCycle: []int64{},
	//}
	//
	//fmt.Printf("%+v\n", p.UntilNext())
	//p.StepBy(1)
	//fmt.Printf("coord: %d, %+v\n", p.coord, p.UntilNext())
	//p.StepBy(5)
	//fmt.Printf("coord: %d, %+v\n", p.coord, p.UntilNext())
}
