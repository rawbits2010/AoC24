package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/rawbits2010/AoC24/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()

	reports, err := readInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	resultP1 := partOne(reports)
	resultP2 := partTwo(reports)

	fmt.Printf("Result - Part1: %d Part2: %d", resultP1, resultP2)
}

func readInput(lines []string) ([][]int, error) {

	reports := make([][]int, len(lines))

	for lineIdx, line := range lines {

		levels := strings.Fields(line)
		if len(levels) == 0 {
			return nil, fmt.Errorf("empty report")
		}

		levelsConv := make([]int, len(levels))
		for idx, val := range levels {
			tmp, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("invalid level in report '%s': %w", line, err)
			}
			levelsConv[idx] = tmp
		}

		reports[lineIdx] = levelsConv
	}

	return reports, nil
}

func partOne(reports [][]int) int {

	var safe int
	for _, levels := range reports {

		// text didn't specify what if there is only 1 level
		if len(levels) < 2 {
			continue
		}

		valid, _ := levelCheck(levels)

		if valid {
			safe++
		}
	}

	return safe
}

func partTwo(reports [][]int) int {

	var safe int
	for _, levels := range reports {

		// text didn't specify what if there is only 1 level
		if len(levels) < 2 {
			continue
		}

		valid, badIdx := levelCheck(levels)
		if !valid {

			// remove bad level
			temp1 := removeAt(levels, badIdx)
			temp2 := removeAt(levels, badIdx-1)

			valid1, _ := levelCheck(temp1)
			valid2, _ := levelCheck(temp2)

			valid = valid1 || valid2

			// edge case, tendency is decided with 1st 2 elements
			// 1 2 1 0
			//  + -
			if badIdx == 2 {

				temp3 := removeAt(levels, 0)
				valid3, _ := levelCheck(temp3)

				valid = valid || valid3
			}
		}

		if valid {
			safe++
		}
	}

	return safe
}

func levelCheck(levels []int) (bool, int) {

	valid := true
	tendency := levels[0] - levels[1]
	idx := 0
	for ; idx < len(levels)-1; idx++ {

		diff := levels[idx] - levels[idx+1]

		if diff == 0 {
			valid = false
			break
		}

		if tendency < 0 && diff > 0 {
			valid = false
			break
		}
		if tendency > 0 && diff < 0 {
			valid = false
			break
		}

		if diff < -3 || diff > 3 {
			valid = false
			break
		}
	}

	return valid, idx + 1
}

func removeAt(list []int, idx int) []int {

	if idx == 0 {
		return list[1:]
	} else if idx == len(list)-1 {
		return list[:idx]
	}

	temp := make([]int, 0, len(list)-1)
	temp = append(temp, list[:idx]...)
	temp = append(temp, list[idx+1:]...)
	return temp
}
