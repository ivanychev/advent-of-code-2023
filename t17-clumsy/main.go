package main

import (
	"advent_of_code/common"
	"fmt"
	"github.com/samber/lo"
	"log"
	"strconv"
	"strings"
)

const MAX_STRAIGHT_STEPS = 10
const MIN_STRAIGHT_STEPS = 4

type MoveState struct {
	x, y          int
	direction     common.DirectionDesc
	sinceLastTurn int
}

type MoveStateList struct {
	MoveState
	prev *MoveStateList
}

type MoveStateListAndLoss struct {
	lst  *MoveStateList
	loss int
}

func (l MoveStateList) HasCoords(coords []lo.Tuple2[int, int]) bool {
	node := &l
	for idx := len(coords) - 1; idx >= 0; idx-- {
		if node == nil || node.x != coords[idx].A || node.y != coords[idx].B {
			return false
		}
		node = node.prev
	}
	return node == nil
}

func (s MoveState) IsCorrect(ctx TaskContext) bool {
	return s.sinceLastTurn <= MAX_STRAIGHT_STEPS && ctx.Field.WithinField(s.x, s.y)
}

type TaskContext struct {
	directions         common.Directions
	minLossByState     *map[MoveState]MoveStateListAndLoss
	minEncounteredLoss *int
	Field
}

func (c TaskContext) PrintMinLosses() string {
	losses := make([][]int, len(c.Field.tiles))
	for idx, _ := range losses {
		losses[idx] = make([]int, 0, len(c.Field.tiles[0]))
		for i := 0; i < len(c.Field.tiles[0]); i++ {
			losses[idx] = append(losses[idx], -1)
		}
	}

	for state, value := range *c.minLossByState {
		if losses[state.y][state.x] == -1 {
			losses[state.y][state.x] = value.loss
		} else {
			losses[state.y][state.x] = min(losses[state.y][state.x], value.loss)
		}
	}

	var sb strings.Builder
	for _, row := range losses {
		for _, value := range row {
			sb.WriteString(fmt.Sprintf("%d ", value))
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (s MoveStateList) NextSteps(ctx TaskContext) []MoveStateList {
	newStates := make([]MoveStateList, 0)
	currentLoss, exists := (*ctx.minLossByState)[s.MoveState]

	if !exists {
		log.Fatalf("Not state not exists")
	}

	if s.sinceLastTurn < MAX_STRAIGHT_STEPS {
		newStates = append(newStates, MoveStateList{
			MoveState: MoveState{
				direction:     s.direction,
				sinceLastTurn: s.sinceLastTurn + 1,
				x:             s.x + s.direction.DeltaX,
				y:             s.y + s.direction.DeltaY,
			},
			prev: &s,
		})
	}
	if s.sinceLastTurn >= MIN_STRAIGHT_STEPS {
		turnDirections := s.direction.Turns()
		for _, direction := range turnDirections {
			newStates = append(newStates, MoveStateList{
				MoveState: MoveState{
					direction:     direction,
					sinceLastTurn: 1,
					x:             s.x + direction.DeltaX,
					y:             s.y + direction.DeltaY,
				},
				prev: &s,
			})
		}
	}
	statesWithinField := lo.Filter(newStates, func(item MoveStateList, index int) bool {
		return item.IsCorrect(ctx)
	})
	statesWithLosses := lo.Map(statesWithinField, func(item MoveStateList, index int) lo.Tuple2[MoveStateList, int] {
		return lo.T2(item, currentLoss.loss+ctx.Field.tiles[item.y][item.x])
	})
	statesWithAppropriateLosses := lo.Filter(statesWithLosses, func(item lo.Tuple2[MoveStateList, int], index int) bool {
		loss := item.B
		state := item.A
		currentStateLoss, exists := (*ctx.minLossByState)[state.MoveState]
		return (*ctx.minEncounteredLoss < 0 || loss < *ctx.minEncounteredLoss) && (!exists || currentStateLoss.loss > loss)
	})
	returnedStates := lo.Map(statesWithAppropriateLosses, func(item lo.Tuple2[MoveStateList, int], index int) MoveStateList {
		(*ctx.minLossByState)[item.A.MoveState] = MoveStateListAndLoss{
			loss: item.B,
			lst:  &item.A,
		}
		if item.A.x == len(ctx.Field.tiles[0])-1 && item.A.y == len(ctx.Field.tiles)-1 && (item.B < *ctx.minEncounteredLoss || *ctx.minEncounteredLoss < 0) {
			*ctx.minEncounteredLoss = item.B
		}
		return item.A
	})
	return returnedStates
}

type Field struct {
	tiles [][]int
}

func (f Field) WithinField(x, y int) bool {
	return 0 <= x && x < len(f.tiles[0]) && 0 <= y && y < len(f.tiles)
}

func NewField(rawField []string) Field {
	tiles := make([][]int, 0, len(rawField))
	for idx, s := range rawField {
		tiles = append(tiles, make([]int, 0, len(rawField[0])))
		for _, c := range s {
			value, err := strconv.Atoi(string(c))
			if err != nil {
				log.Fatalf("Failed to parse %v", c)
			}
			tiles[idx] = append(tiles[idx], value)
		}
	}
	return Field{tiles: tiles}
}

func main() {
	values, err := common.FileToRows("/Users/iv/Code/advent-of-code-2023/t17-clumsy/1.txt")
	if err != nil {
		log.Fatalf("%w", err)
	}
	field := NewField(values)
	directions := common.NewDirections()
	states := []MoveStateList{
		{MoveState: MoveState{x: 0, y: 0, direction: directions.Right, sinceLastTurn: 0}, prev: nil},
		{MoveState: MoveState{x: 0, y: 0, direction: directions.Down, sinceLastTurn: 0}, prev: nil},
	}
	minLossByState := map[MoveState]MoveStateListAndLoss{
		states[0].MoveState: {&states[0], 0},
		states[1].MoveState: {&states[1], 0},
	}
	minEncounteredLoss := -1
	ctx := TaskContext{
		directions:         directions,
		minEncounteredLoss: &minEncounteredLoss,
		minLossByState:     &minLossByState,
		Field:              field,
	}
	iterations := 0
	for len(states) > 0 {
		fmt.Printf("Iterations: %d, states: %d\n", iterations, len(states))
		states = lo.FlatMap(states, func(item MoveStateList, index int) []MoveStateList {
			return item.NextSteps(ctx)
		})
	}
	fmt.Printf("Min loss: %d\n", *ctx.minEncounteredLoss)
}
