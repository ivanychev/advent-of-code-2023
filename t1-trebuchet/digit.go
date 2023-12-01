package main

import (
	"github.com/dlclark/regexp2"
	"log"
	"strconv"
)

type DigitExtractor interface {
	Extract(s string) []int
}

type TrivialExtractor struct {
	re *regexp2.Regexp
}

func (t *TrivialExtractor) Extract(s string) []int {
	var matches []int
	m, _ := t.re.FindStringMatch(s)
	for m != nil {
		value, err := strconv.Atoi(m.String())
		if err != nil {
			log.Fatalf("Failed to parse int from %s: %w", m.String(), err)
		}
		matches = append(matches, value)
		m, _ = t.re.FindNextMatch(m)
	}
	return matches
}

type ComplexExtractor struct {
	re         *regexp2.Regexp
	strToDigit map[string]int
}

func (c *ComplexExtractor) Extract(s string) []int {
	var matches []int
	m, _ := c.re.FindStringMatch(s)
	for m != nil {
		value, ok := c.strToDigit[m.Groups()[1].String()]
		if !ok {
			log.Fatalf("Failed to get int from %s", m.String())
		}
		matches = append(matches, value)
		m, _ = c.re.FindNextMatch(m)
	}
	return matches
}
