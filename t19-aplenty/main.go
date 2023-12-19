package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/samber/lo"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	Accept = "A"
	Reject = "R"
	Input  = "in"
)

var TriggerRegex = regexp2.MustCompile("^(?<prop>\\w{1})(?<op>[<>]{1})(?<threshold>\\d+):(?<outcome>\\w+)$", regexp2.IgnoreCase)
var TerminationTriggerRegex = regexp2.MustCompile("^\\w+$", regexp2.IgnoreCase)
var DetailRe = regexp2.MustCompile("^{x=(?<x>\\d+),m=(?<m>\\d+),a=(?<a>\\d+),s=(?<s>\\d+)}$", regexp2.IgnoreCase)

func AlwaysTrueCond(d Detail) bool {
	return true
}

type WorkflowTrigger struct {
	Cond           func(Detail) bool
	OutputWorkflow string
}

type Workflow struct {
	name     string
	triggers []WorkflowTrigger
}

func (w Workflow) Process(d Detail) string {
	for _, trigger := range w.triggers {
		if trigger.Cond(d) {
			return trigger.OutputWorkflow
		}
	}
	log.Fatalf("Unreacheable")
	return ""
}

type Detail struct {
	x, m, a, s int
}

func (d Detail) Total() int {
	return d.x + d.m + d.a + d.s
}

func ParseTrigger(rawTrigger string) WorkflowTrigger {
	m, err := TriggerRegex.FindStringMatch(rawTrigger)
	if err != nil {
		log.Fatalf("Failed to parse %s", rawTrigger)
	}
	if m == nil {
		m, err = TerminationTriggerRegex.FindStringMatch(rawTrigger)
		if m == nil || err != nil {
			log.Fatalf("Invalid trigger: %s", rawTrigger)
		}
		return WorkflowTrigger{
			Cond:           AlwaysTrueCond,
			OutputWorkflow: rawTrigger,
		}
	}

	testedProerty := m.GroupByName("prop").String()
	operation := m.GroupByName("op").String()
	rawThreshold := m.GroupByName("threshold").String()
	threshold, err := strconv.Atoi(rawThreshold)
	if err != nil {
		log.Fatalf("Failed to parse int %s", rawThreshold)
	}

	var testedProertyAccessor func(Detail) int
	var predicate func(Detail) bool

	switch testedProerty {
	case "x":
		testedProertyAccessor = func(d Detail) int {
			return d.x
		}
	case "m":
		testedProertyAccessor = func(d Detail) int {
			return d.m
		}
	case "a":
		testedProertyAccessor = func(d Detail) int {
			return d.a
		}
	case "s":
		testedProertyAccessor = func(d Detail) int {
			return d.s
		}
	default:
		log.Fatalf("Unknown op: %s", testedProerty)
	}

	switch operation {
	case "<":
		predicate = func(detail Detail) bool {
			return testedProertyAccessor(detail) < threshold
		}
	case ">":
		predicate = func(detail Detail) bool {
			return testedProertyAccessor(detail) > threshold
		}
	default:
		log.Fatalf("Unknown op: %s", operation)
	}

	outcomeWorkflow := m.GroupByName("outcome").String()
	return WorkflowTrigger{
		OutputWorkflow: outcomeWorkflow,
		Cond:           predicate,
	}
}

func ParseWorkflow(raw string) Workflow {
	components := strings.Split(raw, "{")
	name := components[0]
	allRawTriggers := strings.TrimRight(components[1], "}")
	rawTriggers := strings.Split(allRawTriggers, ",")
	triggers := lo.Map(rawTriggers, common.NoIndex(ParseTrigger))
	return Workflow{
		name:     name,
		triggers: triggers,
	}
}

func ParseDetail(rawDetail string) Detail {
	m, err := DetailRe.FindStringMatch(rawDetail)
	if err != nil {
		log.Fatalf("Failed to parse detail %s", rawDetail)
	}
	xRaw := m.GroupByName("x").String()
	mRaw := m.GroupByName("m").String()
	aRaw := m.GroupByName("a").String()
	sRaw := m.GroupByName("s").String()

	return Detail{
		x: common.MustAtoi(xRaw),
		m: common.MustAtoi(mRaw),
		a: common.MustAtoi(aRaw),
		s: common.MustAtoi(sRaw),
	}
}

func ParseWorkflowsAndDetails(contents string) ([]Workflow, []Detail) {
	components := strings.Split(contents, "\n\n")
	rawWorkflows := strings.Split(components[0], "\n")
	rawDetails := strings.Split(components[1], "\n")

	workflows := lo.Map(rawWorkflows, common.NoIndex(ParseWorkflow))
	details := lo.Map(rawDetails, common.NoIndex(ParseDetail))
	return workflows, details
}

func ProcessDetails(details []Detail, workflows []Workflow) []Detail {
	nameToWorkflow := lo.MapValues(lo.GroupBy(workflows, func(w Workflow) string {
		return w.name
	}), func(w []Workflow, key string) Workflow {
		return w[0]
	})
	activeDetails := lo.Map(details, func(d Detail, index int) lo.Tuple2[Detail, string] {
		return lo.T2(d, Input)
	})
	acceptedDetails := make([]Detail, 0)
	for len(activeDetails) > 0 {
		newActiveDetails := make([]lo.Tuple2[Detail, string], 0)
		for _, detail := range activeDetails {
			newWorkflow := nameToWorkflow[detail.B].Process(detail.A)
			switch newWorkflow {
			case Reject:
			case Accept:
				acceptedDetails = append(acceptedDetails, detail.A)
			default:
				newActiveDetails = append(newActiveDetails, lo.T2(detail.A, newWorkflow))
			}
		}
		activeDetails = newActiveDetails
	}
	return acceptedDetails
}

func main() {
	file, err := os.ReadFile("/Users/iv/Code/advent-of-code-2023/t19-aplenty/1.txt")
	if err != nil {
		log.Fatalf("Failed to read file: %w", err)
	}
	workflows, details := ParseWorkflowsAndDetails(string(file))
	acceptedDetails := ProcessDetails(details, workflows)
	fmt.Printf("Sum: %d\n", lo.SumBy(acceptedDetails, func(d Detail) int {
		return d.Total()
	}))
}
