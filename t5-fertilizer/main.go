package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"log"
	"strings"
)

type Range struct {
	Value       int
	Range       int
	Description string
}

func (r Range) WithinRange(value int) bool {
	return value >= r.Value && value < r.Value+r.Range
}

func parseSeeds(rawSeeds string) []int {
	rawNumberAsStrings, _ := strings.CutPrefix(strings.TrimSpace(rawSeeds), "seeds: ")
	return common.StringOfNumbersToInts(rawNumberAsStrings)
}

func parseSeedRanges(rawSeeds string) []Range {
	rawNumberAsStrings, _ := strings.CutPrefix(strings.TrimSpace(rawSeeds), "seeds: ")
	components := common.StringOfNumbersToInts(rawNumberAsStrings)
	ranges := make([]Range, 0, len(components)/2)
	for i := 0; i < len(components)/2; i++ {
		start := components[2*i]
		end := components[2*i+1]
		ranges = append(ranges, Range{start, end - start, "seeds"})
	}
	return ranges
}

func parseRangeMap(rawRangeMap string, sourceMapName string) *treemap.Map {
	treeMap := treemap.NewWith(func(a, b interface{}) int {
		aRange := a.(Range)
		bRange := b.(Range)
		return utils.IntComparator(aRange.Value, bRange.Value)
	})
	//fromDesc, toDesc
	fromDesc, afterDesc, _ := strings.Cut(sourceMapName, "-to-")
	mapNamePrefix := fmt.Sprintf("%s map:\n", sourceMapName)
	rawRangeMap, found := strings.CutPrefix(rawRangeMap, mapNamePrefix)
	if !found {
		log.Fatalf("Failed to cut prefix %s", mapNamePrefix)
	}
	rawRangesCombined := strings.Split(strings.TrimSpace(rawRangeMap), "\n")
	for _, rawRange := range rawRangesCombined {
		components := common.StringOfNumbersToInts(rawRange)
		sourceRange := Range{components[1], components[2], fromDesc}
		destRange := Range{components[0], components[2], afterDesc}
		treeMap.Put(sourceRange, destRange)
	}

	return treeMap
}

func pipelineSearch(value int, maps ...*treemap.Map) int {
	result := value
	latestDesc := "seed"
	for _, m := range maps {
		searchRange := Range{result, 1, ""}
		foundKeyRaw, foundValueRaw := m.Floor(searchRange)
		if foundKeyRaw == nil {
			// Not found, that means that value equals key.
			result = result
			fmt.Printf("For value %d (desc: %s) nothing found, so result is %d\n",
				result, latestDesc, result)
			continue
		}
		foundKey := foundKeyRaw.(Range)
		foundValue := foundValueRaw.(Range)
		if !foundKey.WithinRange(result) {
			// Not found, that means that value equals key.
			result = result
			fmt.Printf("For value %d (desc: %s) nothing found, so result is %d\n",
				result, latestDesc, result)
			continue
		}
		delta := result - foundKey.Value
		newResult := foundValue.Value + delta
		fmt.Printf("For value %d (desc: %s) found k: %+v, v: %+v, so result is %d\n",
			result, latestDesc, foundKey, foundValue, newResult)
		latestDesc = foundValue.Description
		result = newResult
	}
	return result
}

// Part 1
//
//func main() {
//	const path = "/Users/iv/Code/advent-of-code-2023/t5-fertilizer/1.txt"
//	file, err := os.Open(path)
//	if err != nil {
//		log.Fatalf("Failed to open file %s", file)
//	}
//	defer file.Close()
//
//	rawContents, err := io.ReadAll(file)
//	if err != nil {
//		log.Fatalf("Failed to read file %s", file)
//	}
//	contents := string(rawContents)
//	components := strings.Split(contents, "\n\n")
//	seeds := parseSeeds(components[0])
//	seedToSoil := parseRangeMap(components[1], "seed-to-soil")
//	soilToFertilizer := parseRangeMap(components[2], "soil-to-fertilizer")
//	fertilizerToWater := parseRangeMap(components[3], "fertilizer-to-water")
//	waterToLight := parseRangeMap(components[4], "water-to-light")
//	lightToTemperature := parseRangeMap(components[5], "light-to-temperature")
//	temperatureToHumidity := parseRangeMap(components[6], "temperature-to-humidity")
//	humidityToLocation := parseRangeMap(components[7], "humidity-to-location")
//
//	locations := lo.Map(seeds, func(seed int, index int) int {
//		return pipelineSearch(seed, seedToSoil, soilToFertilizer, fertilizerToWater, waterToLight, lightToTemperature, temperatureToHumidity, humidityToLocation)
//	})
//	fmt.Printf("Locations: %d", lo.Min(locations))
//}
