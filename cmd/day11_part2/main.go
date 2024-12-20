package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/rawbits2010/AoC24/internal/inputhandler"
)

// solution inspired by u/Dymatizeee
func main() {

	lines := inputhandler.ReadInput()
	stones := strings.Fields(lines[0])

	stonesMap := make(map[int]int)
	for _, valStr := range stones {

		val, err := strconv.Atoi(valStr)
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := stonesMap[val]; !ok {
			stonesMap[val] = 1
		} else {
			stonesMap[val]++
		}
	}

	var err error
	var resultP1 int
	for blinkCount := 1; blinkCount <= 75; blinkCount++ {

		stonesMap, err = blink(stonesMap)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d - %d\n", blinkCount, sumResult(stonesMap))

		if blinkCount == 25 {
			resultP1 = sumResult(stonesMap)
		}
	}

	resultP2 := sumResult(stonesMap)

	fmt.Printf("Result - Part 1: %d, Part 2: %d", resultP1, resultP2)
}

func sumResult(stonesMap map[int]int) int {
	var result int
	for _, count := range stonesMap {
		result += count
	}
	return result
}

func blink(stonesMap map[int]int) (map[int]int, error) {

	tempStonesMap := make(map[int]int)

	for key, val := range stonesMap {

		if key == 0 {
			if _, ok := tempStonesMap[1]; !ok {
				tempStonesMap[1] = val
			} else {
				tempStonesMap[1] += val
			}
			continue
		}

		keyStr := strconv.Itoa(key)
		if len(keyStr)%2 == 0 {

			newLeftStr := keyStr[:len(keyStr)/2]
			newLeft, err := strconv.Atoi(newLeftStr)
			if err != nil {
				return nil, fmt.Errorf("error converting left value '%s': %w", newLeftStr, err)
			}
			if _, ok := tempStonesMap[newLeft]; !ok {
				tempStonesMap[newLeft] = val
			} else {
				tempStonesMap[newLeft] += val
			}

			newRightStr := keyStr[len(keyStr)/2:]
			newRight, err := strconv.Atoi(newRightStr)
			if err != nil {
				return nil, fmt.Errorf("error converting left value '%s': %w", newRightStr, err)
			}
			if _, ok := tempStonesMap[newRight]; !ok {
				tempStonesMap[newRight] = val
			} else {
				tempStonesMap[newRight] += val
			}

			continue
		}

		newKey, err := mulInt(key, 2024)
		if err != nil {
			return nil, fmt.Errorf("error multiplying '%d': %w", key, err)
		}
		if _, ok := tempStonesMap[newKey]; !ok {
			tempStonesMap[newKey] = val
		} else {
			tempStonesMap[newKey] += val
		}
	}

	return tempStonesMap, nil
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
