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
	Accept             = "A"
	Reject             = "R"
	Input              = "in"
	MIN_RANGE          = 1
	MAX_RANGE_EXCLUDED = 4000 + 1
)

var TriggerRegex = regexp2.MustCompile("^(?<prop>\\w{1})(?<op>[<>]{1})(?<threshold>\\d+):(?<outcome>\\w+)$", regexp2.IgnoreCase)
var TerminationTriggerRegex = regexp2.MustCompile("^\\w+$", regexp2.IgnoreCase)
var DetailRe = regexp2.MustCompile("^{x=(?<x>\\d+),m=(?<m>\\d+),a=(?<a>\\d+),s=(?<s>\\d+)}$", regexp2.IgnoreCase)

func AlwaysTrueCond(d Detail) bool {
	return true
}

func AlwaysTrueCondRange(d DetailRange) (DetailRange, DetailRange) {
	return d, DetailRange{}
}

type WorkflowTrigger struct {
	Cond           func(Detail) bool
	CondRange      func(detailRange DetailRange) (DetailRange, DetailRange)
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

func (w Workflow) ProcessRange(d DetailRange) []lo.Tuple2[DetailRange, string] {
	returned := make([]lo.Tuple2[DetailRange, string], 0)
	for _, trigger := range w.triggers {
		trueRange, falseRange := trigger.CondRange(d)
		if !trueRange.IsEmpty() {
			returned = append(returned, lo.T2(trueRange, trigger.OutputWorkflow))
		}
		d = falseRange
	}
	return returned
}

type Detail struct {
	x, m, a, s int
}

func (d Detail) ToRange() DetailRange {
	return DetailRange{
		xMin: d.x,
		xMax: d.x + 1,
		mMin: d.m,
		mMax: d.m + 1,
		aMin: d.a,
		aMax: d.a + 1,
		sMin: d.s,
		sMax: d.s + 1,
	}
}

type DetailRange struct {
	xMin, xMax, mMin, mMax, aMin, aMax, sMin, sMax int
}

func (r DetailRange) ToString() string {
	return fmt.Sprintf("x[%d,%d) m[%d,%d) a[%d,%d) s[%d,%d)", r.xMin, r.xMax, r.mMin, r.mMax, r.aMin, r.aMax, r.sMin, r.sMax)
}

func (r DetailRange) SplitAtX(x int) (DetailRange, DetailRange) {
	if x < r.xMin {
		return DetailRange{}, r
	} else if x > r.xMax {
		return r, DetailRange{}
	} else {
		left := r
		right := r
		left.xMax = x
		right.xMin = x
		return left, right
	}
}

func (r DetailRange) SplitAtM(m int) (DetailRange, DetailRange) {
	if m < r.mMin {
		return DetailRange{}, r
	} else if m > r.mMax {
		return r, DetailRange{}
	} else {
		left := r
		right := r
		left.mMax = m
		right.mMin = m
		return left, right
	}
}

func (r DetailRange) SplitAtA(a int) (DetailRange, DetailRange) {
	if a < r.aMin {
		return DetailRange{}, r
	} else if a > r.aMax {
		return r, DetailRange{}
	} else {
		left := r
		right := r
		left.aMax = a
		right.aMin = a
		return left, right
	}
}

func (r DetailRange) SplitAtS(s int) (DetailRange, DetailRange) {
	if s < r.sMin {
		return DetailRange{}, r
	} else if s > r.sMax {
		return r, DetailRange{}
	} else {
		left := r
		right := r
		left.sMax = s
		right.sMin = s
		return left, right
	}
}

func (r DetailRange) Size() int {
	return (r.xMax - r.xMin) * (r.mMax - r.mMin) * (r.aMax - r.aMin) * (r.sMax - r.sMin)
}

func (r DetailRange) DetailTotal() int {
	return r.xMin + r.mMin + r.aMin + r.sMin
}

func (r DetailRange) IsEmpty() bool {
	return r.Size() == 0
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
			CondRange:      AlwaysTrueCondRange,
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
	var splitter func(detailRange DetailRange, at int) (DetailRange, DetailRange)
	var predicate func(Detail) bool
	var rangePredicate func(detailRange DetailRange) (DetailRange, DetailRange)

	switch testedProerty {
	case "x":
		testedProertyAccessor = func(d Detail) int {
			return d.x
		}
		splitter = func(detailRange DetailRange, at int) (DetailRange, DetailRange) {
			return detailRange.SplitAtX(at)
		}
	case "m":
		testedProertyAccessor = func(d Detail) int {
			return d.m
		}
		splitter = func(detailRange DetailRange, at int) (DetailRange, DetailRange) {
			return detailRange.SplitAtM(at)
		}
	case "a":
		testedProertyAccessor = func(d Detail) int {
			return d.a
		}
		splitter = func(detailRange DetailRange, at int) (DetailRange, DetailRange) {
			return detailRange.SplitAtA(at)
		}
	case "s":
		testedProertyAccessor = func(d Detail) int {
			return d.s
		}
		splitter = func(detailRange DetailRange, at int) (DetailRange, DetailRange) {
			return detailRange.SplitAtS(at)
		}
	default:
		log.Fatalf("Unknown op: %s", testedProerty)
	}

	switch operation {
	case "<":
		predicate = func(detail Detail) bool {
			return testedProertyAccessor(detail) < threshold
		}
		rangePredicate = func(detailRange DetailRange) (DetailRange, DetailRange) {
			beforeThreshold, hereAndAfter := splitter(detailRange, threshold)
			return beforeThreshold, hereAndAfter
		}
	case ">":
		predicate = func(detail Detail) bool {
			return testedProertyAccessor(detail) > threshold
		}
		rangePredicate = func(detailRange DetailRange) (DetailRange, DetailRange) {
			lowerOrEqual, greaterThan := splitter(detailRange, threshold+1)
			return greaterThan, lowerOrEqual
		}
	default:
		log.Fatalf("Unknown op: %s", operation)
	}

	outcomeWorkflow := m.GroupByName("outcome").String()
	return WorkflowTrigger{
		OutputWorkflow: outcomeWorkflow,
		Cond:           predicate,
		CondRange:      rangePredicate,
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

func DebugPrintActiveDetails(m []lo.Tuple2[DetailRange, string]) {
	nameToRanges := lo.GroupBy(m, func(item lo.Tuple2[DetailRange, string]) string {
		return item.B
	})
	for name, ranges := range nameToRanges {
		fmt.Printf("%s\n", name)
		for _, r := range ranges {
			fmt.Printf("  %s\n", r.A.ToString())
		}
	}
	fmt.Printf("---\n")
}

func ProcessDetailRange(detailRanges []DetailRange, workflows []Workflow) []DetailRange {
	nameToWorkflow := lo.MapValues(lo.GroupBy(workflows, func(w Workflow) string {
		return w.name
	}), func(w []Workflow, key string) Workflow {
		return w[0]
	})
	activeDetails := lo.Map(detailRanges, func(d DetailRange, i int) lo.Tuple2[DetailRange, string] {
		return lo.T2(d, Input)
	})
	acceptedDetails := make([]DetailRange, 0)
	for len(activeDetails) > 0 {
		newActiveDetails := make([]lo.Tuple2[DetailRange, string], 0)
		for _, detail := range activeDetails {
			newRangesAndWorkflows := nameToWorkflow[detail.B].ProcessRange(detail.A)
			for _, pair := range newRangesAndWorkflows {
				rng := pair.A
				workflow := pair.B
				switch workflow {
				case Reject:
				case Accept:
					acceptedDetails = append(acceptedDetails, rng)
				default:
					newActiveDetails = append(newActiveDetails, lo.T2(rng, workflow))
				}
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
	workflows, _ := ParseWorkflowsAndDetails(string(file))

	//detailRanges := lo.Map(details, func(item Detail, index int) DetailRange {
	//	return item.ToRange()
	//})

	detailRanges := []DetailRange{{
		xMin: MIN_RANGE,
		xMax: MAX_RANGE_EXCLUDED,
		mMin: MIN_RANGE,
		mMax: MAX_RANGE_EXCLUDED,
		aMin: MIN_RANGE,
		aMax: MAX_RANGE_EXCLUDED,
		sMin: MIN_RANGE,
		sMax: MAX_RANGE_EXCLUDED,
	}}

	acceptedDetails := ProcessDetailRange(detailRanges, workflows)
	fmt.Printf("Sum: %d\n", lo.SumBy(acceptedDetails, func(d DetailRange) int {
		return d.Size()
	}))
	//fmt.Printf("Sum: %d\n", lo.SumBy(acceptedDetails, func(d DetailRange) int {
	//	return d.DetailTotal()
	//}))
}

// 167409079868000
// 172703741240000
