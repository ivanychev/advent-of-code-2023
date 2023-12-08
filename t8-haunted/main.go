package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

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

func main() {
	nameToEdge, commands := createEdgesAndCommands("/Users/iv/Code/advent-of-code-2023/t8-haunted/1.txt")

	currentStep := -1
	totalSteps := 0
	currentEdge := nameToEdge["AAA"]
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
	fmt.Printf("%d", totalSteps)
}
