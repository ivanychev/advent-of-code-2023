package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"github.com/zyedidia/generic/queue"
	"log"
	"strings"
)

const (
	BroadcasterToken = "broadcaster"
	Button           = "button"
	FlipFlop         = "%"
	Conjunction      = "&"
	High             = 1
	Low              = 0
	PulseCount       = 1000
)

type Pulse struct {
	strength    int
	source      string
	destination string
}

type Module interface {
	Name() string
	GetDestinations() []string
	ApplyPulse(p Pulse) []Pulse
}

type BaseModule struct {
	name         string
	destinations []string
}

func (b BaseModule) GetDestinations() []string {
	return b.destinations
}

func (b BaseModule) Name() string {
	return b.name
}

type NoopModule struct {
	BaseModule
}

func (b *NoopModule) ApplyPulse(p Pulse) []Pulse {
	return lo.Map(b.destinations, func(item string, index int) Pulse {
		return Pulse{strength: p.strength, destination: item, source: b.name}
	})
}

type FlipFlopModule struct {
	BaseModule
	isOn bool
}

func (f *FlipFlopModule) ApplyPulse(p Pulse) []Pulse {
	if p.strength == High {
		return nil
	}
	f.isOn = !f.isOn
	var sent int
	if f.isOn {
		sent = High
	} else {
		sent = Low
	}
	return lo.Map(f.destinations, func(item string, index int) Pulse {
		return Pulse{strength: sent, destination: item, source: f.name}
	})
}

type ConjunctionModule struct {
	BaseModule
	sourceToPulse   map[string]int
	highInputsCount int
}

func (c *ConjunctionModule) ApplyPulse(p Pulse) []Pulse {
	prev := c.sourceToPulse[p.source]
	c.sourceToPulse[p.source] = p.strength
	c.highInputsCount += p.strength - prev
	var sent int
	if c.highInputsCount == len(c.sourceToPulse) {
		sent = Low
	} else {
		sent = High
	}
	return lo.Map(c.destinations, func(item string, index int) Pulse {
		return Pulse{strength: sent, destination: item, source: c.name}
	})
}

func RawModuleName(name string) string {
	return strings.TrimLeft(name, "%&")
}

func ParseConnections(rows []string) map[string]Module {
	moduleAndDestinations := lo.Map(rows, func(row string, index int) lo.Tuple2[string, []string] {
		components := strings.Split(row, " -> ")
		module, rawDests := components[0], components[1]
		destinations := strings.Split(rawDests, ", ")
		return lo.T2(module, destinations)
	})
	encounteredModules := lo.Uniq(lo.FlatMap(moduleAndDestinations, func(item lo.Tuple2[string, []string], index int) []string {
		modules := make([]string, 1+len(item.B))
		modules = append(modules, item.A)
		modules = append(modules, item.B...)
		return modules
	}))

	nameToSources := make(map[string][]string)
	for _, pair := range moduleAndDestinations {
		source := pair.A
		for _, dest := range pair.B {
			if _, exists := nameToSources[dest]; !exists {
				nameToSources[dest] = make([]string, 0)
			}
			nameToSources[dest] = append(nameToSources[dest], RawModuleName(source))
		}
	}

	modules := lo.Map(moduleAndDestinations, func(pair lo.Tuple2[string, []string], index int) Module {
		module := pair.A
		dests := pair.B

		if module == BroadcasterToken {
			return &NoopModule{BaseModule{name: module, destinations: dests}}
		} else if strings.HasPrefix(module, FlipFlop) {
			return &FlipFlopModule{
				BaseModule: BaseModule{name: RawModuleName(module), destinations: dests}, isOn: false}
		} else if strings.HasPrefix(module, Conjunction) {
			return &ConjunctionModule{
				BaseModule: BaseModule{name: RawModuleName(module), destinations: dests},
				sourceToPulse: lo.SliceToMap(nameToSources[RawModuleName(module)], func(source string) (string, int) {
					return source, Low
				})}
		} else {
			log.Fatalf("Unknown module: %s", module)
		}
		return nil
	})

	moduleMap := lo.SliceToMap(modules, func(m Module) (string, Module) {
		return m.Name(), m
	})

	for _, module := range encounteredModules {
		if _, exists := moduleMap[module]; !exists {
			m := NoopModule{BaseModule{name: module, destinations: nil}}
			moduleMap[module] = &m
		}
	}

	return moduleMap
}

func Push(nameToModule *map[string]Module) (int64, int64) {
	q := queue.New[Pulse]()
	var lows, highs int64
	q.Enqueue(Pulse{strength: Low, source: Button, destination: BroadcasterToken})
	for !q.Empty() {
		pulse := q.Dequeue()

		if pulse.strength == Low {
			lows += 1
		} else {
			highs += 1
		}

		destination, exists := (*nameToModule)[pulse.destination]
		if !exists {
			log.Fatalf("%s doesn't exist", pulse.destination)
		}
		for _, newPulse := range destination.ApplyPulse(pulse) {
			q.Enqueue(newPulse)
		}
	}
	return lows, highs
}

func main() {
	contents, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t20-pulse/1.txt")
	if err != nil {
		log.Fatalf("failed to read file: %w", err)
	}

	nameToModule := ParseConnections(contents)
	var totalLows, totalHighs int64
	for i := 0; i < PulseCount; i++ {
		lows, highs := Push(&nameToModule)
		totalLows += lows
		totalHighs += highs
	}
	fmt.Printf("Total: %d\n", totalLows*totalHighs)

}
