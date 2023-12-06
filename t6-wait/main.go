package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"log"
	"strings"
)

type LapRecord struct {
	Time     int64
	Distance int64
}

type Run struct {
	WaitFor int64
	RunFor  int64
}

func (r LapRecord) RecordBeats() []Run {
	runs := make([]Run, 0)
	for wait := int64(1); wait <= r.Time-1; wait++ {
		run := r.Time - wait
		if run*wait > r.Distance {
			runs = append(runs, Run{wait, run})
		}
	}
	return runs
}

func (r LapRecord) RecordBeatCount() int64 {
	beats := int64(0)
	for wait := int64(1); wait <= r.Time-1; wait++ {
		run := r.Time - wait
		if run*wait > r.Distance {
			beats++
		}
	}
	return beats
}

func ReadRecords(path string) []LapRecord {
	rows, err := common.FileToRows(path)
	if err != nil {
		log.Fatalf("Failed to read file: %s", path)
	}

	rawTimes, _ := strings.CutPrefix(rows[0], "Time:")
	rawDistances, _ := strings.CutPrefix(rows[1], "Distance:")
	times := common.StringOfNumbersToInts(rawTimes)
	distances := common.StringOfNumbersToInts(rawDistances)
	runs := make([]LapRecord, 0, len(times))
	for i, t := range times {
		runs = append(runs, LapRecord{Time: int64(t), Distance: int64(distances[i])})
	}
	return runs
}

func main() {
	runs := ReadRecords("/Users/iv/Code/advent-of-code-2023/t6-wait/2.txt")
	recordBeats := lo.Map(runs, func(r LapRecord, index int) int64 {
		return r.RecordBeatCount()
	})

	multiplied := lo.Reduce(recordBeats, func(agg int64, r int64, index int) int64 {
		return agg * r
	}, int64(1))

	fmt.Printf("Multiplied: %+v\n", multiplied)
}
