package main

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/rawbits2010/AoC24/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()

	leftList, rightList, err := readInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	slices.Sort(leftList)
	slices.Sort(rightList)

	var resultP1 int
	for idx := 0; idx < len(leftList); idx++ {
		temp := leftList[idx] - rightList[idx]
		if temp < 0 {
			temp = -temp
		}
		resultP1 += temp
	}

	numOccurrence := make(map[int]int)
	for _, num := range rightList {

		count, ok := numOccurrence[num]
		if !ok {
			numOccurrence[num] = 1
		} else {
			numOccurrence[num] = count + 1
		}
	}

	var resultP2 int
	for _, num := range leftList {

		count, ok := numOccurrence[num]
		if !ok {
			count = 0
		}

		if count == 0 {
			continue
		}

		temp, err := mulInt(num, count)
		if err != nil {
			log.Fatal(err)
		}
		resultP2, err = addInt(resultP2, temp)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Result - Part1: %d Part2: %d", resultP1, resultP2)
}

func readInput(lines []string) ([]int, []int, error) {

	left := make([]int, len(lines))
	right := make([]int, len(lines))

	for lineIdx, line := range lines {

		temp := strings.Fields(line)
		if len(temp) != 2 {
			return nil, nil, fmt.Errorf("invalid line - expected 2 ints separated by spaces '%s'", line)
		}

		leftVal, err := strconv.Atoi(temp[0])
		if err != nil {
			return nil, nil, fmt.Errorf("invalid left value in '%s': %w", line, err)
		}
		rightVal, err := strconv.Atoi(temp[1])
		if err != nil {
			return nil, nil, fmt.Errorf("invalid left value in '%s': %w", line, err)
		}

		left[lineIdx] = leftVal
		right[lineIdx] = rightVal
	}

	return left, right, nil
}

func mulInt(val1, val2 int) (int, error) {
	res := val1 * val2
	if (res < 0) == ((val1 < 0) != (val2 < 0)) {
		if res/val2 == val1 {
			return res, nil
		}
	}
	return 0, fmt.Errorf("multiplication overflow")
}

func addInt(val1, val2 int) (int, error) {
	res := val1 + val2
	if (res > val1) == (val2 > 0) {
		return res, nil
	}
	return 0, fmt.Errorf("addition overflow")
}
