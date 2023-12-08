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
	edge                 *Edge
}

func findCycle(e *Edge, commands []byte) CycleInfo {
	commandPointer := int64(0)
	totalSteps := int64(0)
	edgesSet := make(map[EdgeStep]int64)
	edgesSet[EdgeStep{e, totalSteps}] = int64(0)
	for {
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
				edge:                 e,
			}
		} else {
			edgesSet[es] = totalSteps
		}
	}
}

func main() {
	nameToEdge, commands := createEdgesAndCommands("/Users/iv/Code/advent-of-code-2023/t8-haunted/1.txt")

	edges := lo.Filter(lo.Values(nameToEdge), func(item *Edge, index int) bool {
		return item.name[LAST] == 'A'
	})
	fmt.Printf("Found %d edges\n", len(edges))
	for _, edge := range edges {
		fmt.Printf("%+v\n", findCycle(edge, commands))
	}
}
