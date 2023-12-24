package main

import (
	"advent_of_code/common"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"log"
	"os"
	"strings"
)

var ErrNoIntersection = errors.New("no intersection")

const (
	//Xmin = 7
	//Xmax = 27
	Xmin = 200000000000000
	Xmax = 400000000000000
	Ymin = Xmin
	Ymax = Xmax
)

type Coord struct {
	x, y, z float64
}

func (c Coord) WithinBoundsXY() bool {
	return Xmin <= c.x && c.x <= Xmax && Ymin <= c.y && c.y <= Ymax
}

type Stone struct {
	x, y, z, vx, vy, vz int64
}

func Sign(v float64) int {
	if v > 0 {
		return 1
	} else if v < 0 {
		return -1
	}
	return 0
}

func (s Stone) IntersectsWithXY(o Stone) (Coord, error) {

	div := float64(o.vx)*float64(s.vy) - float64(o.vy)*float64(s.vx)
	kx := float64(o.vx)*float64(o.y)*float64(s.vx) - float64(o.vx)*float64(s.vx)*float64(s.y) + float64(o.vx)*float64(s.vy)*float64(s.x) - float64(o.vy)*float64(o.x)*float64(s.vx)
	ky := float64(o.vx)*float64(o.y)*float64(s.vy) - float64(o.vy)*float64(o.x)*float64(s.vy) - float64(o.vy)*float64(s.vx)*float64(s.y) + float64(o.vy)*float64(s.vy)*float64(s.x)

	if div == 0 {
		return Coord{}, ErrNoIntersection
	}
	xInt := float64(kx) / float64(div)
	yInt := float64(ky) / float64(div)
	return Coord{xInt, yInt, float64(0)}, nil
}

func (s Stone) InFuture(c Coord) bool {
	return Sign(c.x-float64(s.x)) == Sign(float64(s.vx))
}

var Replacer = strings.NewReplacer(",", "", "@", "")

func ParseStone(s string) Stone {
	// 343240821178976, 142303638369464, 376763854620819 @ -104, 127, -12
	s = Replacer.Replace(s)
	fields := strings.Fields(s)
	stone := Stone{
		x:  int64(common.MustAtoi(fields[0])),
		y:  int64(common.MustAtoi(fields[1])),
		z:  int64(common.MustAtoi(fields[2])),
		vx: int64(common.MustAtoi(fields[3])),
		vy: int64(common.MustAtoi(fields[4])),
		vz: int64(common.MustAtoi(fields[5])),
	}
	return stone
}

func main() {
	rows, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t24-never-tell-me-the-odds/1.txt")

	if err != nil {
		log.Fatalf("%w", err)
	}
	correctRaw, err := os.ReadFile("/Users/iv/Code/advent-of-code-2023/t24-never-tell-me-the-odds/correct.json")
	if err != nil {
		log.Fatalf("%w", err)
	}

	var correct [][]int
	json.Unmarshal(correctRaw, &correct)
	correctSet := make(map[lo.Tuple2[int, int]]struct{})

	for _, row := range correct {
		correctSet[lo.T2(row[0], row[1])] = struct{}{}
	}
	//actualSet := make(map[lo.Tuple2[int, int]]struct{})
	stones := lo.Map(rows, func(s string, index int) Stone {
		return ParseStone(s)
	})
	total := 0

	//{A:49 B:242}, false

	//a := ParseStone("345505256784794, 468834640747538, 167655858528405 @ -80, -656, 66")
	//b := ParseStone("222301148235939, 318207219159249, 245399299867216 @ 70, -68, 94")
	//coord, err := a.IntersectsWithXY(b)
	//fmt.Printf("%+v %w", coord, err)
	//fmt.Printf("Total rows %d\n", len(rows))

	for i := 0; i < len(stones); i++ {
		for j := i + 1; j < len(stones); j++ {
			coord, err := stones[i].IntersectsWithXY(stones[j])
			if errors.Is(err, ErrNoIntersection) {
				continue
			}
			if coord.WithinBoundsXY() && stones[i].InFuture(coord) && stones[j].InFuture(coord) {
				//actualSet[lo.T2(i, j)] = struct{}{}
				total += 1
			}
		}
	}

	fmt.Printf("Total: %d\n", total)
	//
	//for k, _ := range correctSet {
	//	_, exists := actualSet[k]
	//	fmt.Printf("%+v, %v\n", k, exists)
	//}
}

// 19978 too low
