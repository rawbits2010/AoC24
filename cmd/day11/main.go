package main

import (
	"container/list"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/rawbits2010/AoC24/internal/inputhandler"
)

// this brute force will choke on Part 2
func main() {

	lines := inputhandler.ReadInput()
	stones := strings.Fields(lines[0])
	stonesList, err := makeLinkedList(stones)
	if err != nil {
		log.Fatal(err)
	}

	for blinkCount := 1; blinkCount <= 25; blinkCount++ {
		fmt.Println(blinkCount)

		err := blinkWithList(stonesList)
		if err != nil {
			log.Fatal(err)
		}
	}
	resultP1 := stonesList.Len()

	fmt.Printf("Result - Part 1: %d", resultP1)
}

func makeLinkedList(stones []string) (*list.List, error) {

	stoneList := list.New()
	for idx := 0; idx < len(stones); idx++ {

		val, err := strconv.Atoi(stones[idx])
		if err != nil {
			return nil, fmt.Errorf("error converting value '%s'", stones[idx])
		}

		stoneList.PushBack(val)
	}

	return stoneList, nil
}

func blinkWithList(stones *list.List) error {

	currElement := stones.Front()

	for {
		if currElement == nil {
			break
		}

		currValue := currElement.Value.(int)

		if currValue == 0 {
			currElement.Value = 1
			currElement = currElement.Next()
			continue
		}

		currValueStr := strconv.Itoa(currValue)
		if len(currValueStr)%2 == 0 {

			newLeftStr := currValueStr[:(len(currValueStr) / 2)]
			newLeft, err := strconv.Atoi(newLeftStr)
			if err != nil {
				return fmt.Errorf("error converting value '%s'", newLeftStr)
			}
			currElement.Value = newLeft

			newRightStr := currValueStr[(len(currValueStr) / 2):]
			newRight, err := strconv.Atoi(newRightStr)
			if err != nil {
				return fmt.Errorf("error converting value '%s'", newLeftStr)
			}
			currElement = stones.InsertAfter(newRight, currElement)

			currElement = currElement.Next()
			continue
		}

		val, err := mulInt(currValue, 2024)
		if err != nil {
			return fmt.Errorf("overflow happened when multiplying '%d'", currValue)
		}
		currElement.Value = val

		currElement = currElement.Next()
	}

	return nil
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
