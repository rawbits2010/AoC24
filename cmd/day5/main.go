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

	ruleset, updates, err := parseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	var resultP1 int
	var resultP2 int
	for updateIdx, update := range updates {
		if verifyUpdate(ruleset, update) {

			if len(update)%2 == 0 {
				log.Fatalf("error in update at index %d - page count must be odd", updateIdx)
			}

			middlePage := update[(len(update) / 2)]
			resultP1 += middlePage
		} else {

			fixedUpdate := fixUpdate(ruleset, update)

			middlePage := fixedUpdate[(len(fixedUpdate) / 2)]
			resultP2 += middlePage
		}
	}

	fmt.Printf("Result - Part1: %d Part2: %d", resultP1, resultP2)
}

func parseInput(lines []string) (map[int][]int, [][]int, error) {

	rulset := make(map[int][]int)
	updateList := make([][]int, 0, 250)

	parsingRules := true
	for _, line := range lines {

		if len(line) == 0 {
			parsingRules = false
			continue
		}

		if parsingRules {

			temp := strings.Split(line, "|")
			if len(temp) != 2 {
				return nil, nil, fmt.Errorf("error parsing rule '%s'", line)
			}

			x, err := strconv.Atoi(temp[0])
			if err != nil {
				return nil, nil, fmt.Errorf("error converting X in rule '%s'", line)
			}
			y, err := strconv.Atoi(temp[1])
			if err != nil {
				return nil, nil, fmt.Errorf("error converting Y in rule '%s'", line)
			}

			_, ok := rulset[x]
			if !ok {
				rulset[x] = make([]int, 0, 30)
			}

			rulset[x] = append(rulset[x], y)

		} else {

			temp := strings.Split(line, ",")
			if len(temp) == 0 {
				return nil, nil, fmt.Errorf("found empty update line")
			}

			updateVals := make([]int, len(temp))

			for valIdx, valStr := range temp {

				val, err := strconv.Atoi(valStr)
				if err != nil {
					return nil, nil, fmt.Errorf("error converting page number '%s' in update '%s'", valStr, line)
				}

				updateVals[valIdx] = val
			}

			updateList = append(updateList, updateVals)
		}
	}

	return rulset, updateList, nil
}

func verifyUpdate(ruleset map[int][]int, update []int) bool {

	for currPageIdx, page := range update {

		rules, ok := ruleset[page]
		if !ok {
			continue
		}

		for idx := 0; idx < currPageIdx; idx++ {
			if contains(rules, update[idx]) {
				return false
			}
		}
	}

	return true
}

func contains(list []int, val int) bool {
	for _, currVal := range list {
		if val == currVal {
			return true
		}
	}
	return false
}

func fixUpdate(ruleset map[int][]int, update []int) []int {

fixUpdateStart:
	for pageIdx := len(update) - 1; pageIdx >= 0; {
		currPage := update[pageIdx]

		rules, ok := ruleset[currPage]
		if ok {

			for idx := 0; idx < pageIdx; idx++ {
				if contains(rules, update[idx]) {

					temp := make([]int, 0)
					temp = append(temp, update[:idx]...)
					temp = append(temp, currPage)
					temp = append(temp, update[idx:pageIdx]...)
					temp = append(temp, update[pageIdx+1:]...)
					update = temp

					continue fixUpdateStart
				}
			}
		}

		pageIdx--
	}

	return update
}
