package main

import (
	"fmt"
	"log"

	"github.com/rawbits2010/AoC24/internal/inputhandler"
)

const wordToFindP1 = "XMAS"
const wordToFindP2 = "MAS"

func main() {

	lines := inputhandler.ReadInput()

	resultP1 := searchForWord(wordToFindP1, lines)

	resultP2, err := searchPattern(wordToFindP2, lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Result - Part1: %d Part2: %d", resultP1, resultP2)
}

//
// Part 1

func searchForWord(word string, table []string) int {

	var matchCount int

	wordMaxIdx := len(word) - 1
	vIdxMax := len(table) - 1
	hIdxMax := len(table[0]) - 1
	for vIdx := 0; vIdx <= vIdxMax; vIdx++ {
		for hIdx := 0; hIdx <= hIdxMax; hIdx++ {

			if table[vIdx][hIdx] != word[0] {
				continue
			}

			hasSpaceUp := vIdx-wordMaxIdx >= 0
			hasSpaceDown := vIdx+wordMaxIdx <= vIdxMax
			hasSpaceRight := hIdx+wordMaxIdx <= hIdxMax
			hasSpaceLeft := hIdx-wordMaxIdx >= 0

			if hasSpaceRight {
				if checkHorRight(word, hIdx, vIdx, table) {
					//fmt.Printf(" R: %d, %d\n", hIdx, vIdx)
					matchCount++
				}
			}

			if hasSpaceLeft {
				if checkHorLeft(word, hIdx, vIdx, table) {
					//fmt.Printf(" L: %d, %d\n", hIdx, vIdx)
					matchCount++
				}
			}

			if hasSpaceUp {
				if checkVerUp(word, hIdx, vIdx, table) {
					//fmt.Printf(" U: %d, %d\n", hIdx, vIdx)
					matchCount++
				}
			}

			if hasSpaceDown {
				if checkVerDown(word, hIdx, vIdx, table) {
					//fmt.Printf(" D: %d, %d\n", hIdx, vIdx)
					matchCount++
				}
			}

			if hasSpaceUp && hasSpaceRight {
				if checkDiaUpRight(word, hIdx, vIdx, table) {
					//fmt.Printf("UR: %d, %d\n", hIdx, vIdx)
					matchCount++
				}
			}

			if hasSpaceUp && hasSpaceLeft {
				if checkDiaUpLeft(word, hIdx, vIdx, table) {
					//fmt.Printf("UL: %d, %d\n", hIdx, vIdx)
					matchCount++
				}
			}

			if hasSpaceDown && hasSpaceRight {
				if checkDiaDownRight(word, hIdx, vIdx, table) {
					//fmt.Printf("DR: %d, %d\n", hIdx, vIdx)
					matchCount++
				}
			}

			if hasSpaceDown && hasSpaceLeft {
				if checkDiaDownLeft(word, hIdx, vIdx, table) {
					//fmt.Printf("DL: %d, %d\n", hIdx, vIdx)
					matchCount++
				}
			}
		}
	}

	return matchCount
}

func checkHorRight(word string, hIdxStart, vIdxStart int, table []string) bool {

	hIdx := hIdxStart
	for wIdx := 0; wIdx < len(word); wIdx++ {
		if table[vIdxStart][hIdx] != word[wIdx] {
			return false
		}
		hIdx++
	}

	return true
}

func checkHorLeft(word string, hIdxStart, vIdxStart int, table []string) bool {

	hIdx := hIdxStart
	for wIdx := 0; wIdx < len(word); wIdx++ {
		if table[vIdxStart][hIdx] != word[wIdx] {
			return false
		}
		hIdx--
	}

	return true
}

func checkVerUp(word string, hIdxStart, vIdxStart int, table []string) bool {

	vIdx := vIdxStart
	for wIdx := 0; wIdx < len(word); wIdx++ {
		if table[vIdx][hIdxStart] != word[wIdx] {
			return false
		}
		vIdx--
	}

	return true
}

func checkVerDown(word string, hIdxStart, vIdxStart int, table []string) bool {

	vIdx := vIdxStart
	for wIdx := 0; wIdx < len(word); wIdx++ {
		if table[vIdx][hIdxStart] != word[wIdx] {
			return false
		}
		vIdx++
	}

	return true
}

func checkDiaUpRight(word string, hIdxStart, vIdxStart int, table []string) bool {

	vIdx := vIdxStart
	hIdx := hIdxStart
	for wIdx := 0; wIdx < len(word); wIdx++ {
		if table[vIdx][hIdx] != word[wIdx] {
			return false
		}
		vIdx--
		hIdx++
	}

	return true
}

func checkDiaUpLeft(word string, hIdxStart, vIdxStart int, table []string) bool {

	vIdx := vIdxStart
	hIdx := hIdxStart
	for wIdx := 0; wIdx < len(word); wIdx++ {
		if table[vIdx][hIdx] != word[wIdx] {
			return false
		}
		vIdx--
		hIdx--
	}

	return true
}

func checkDiaDownRight(word string, hIdxStart, vIdxStart int, table []string) bool {

	vIdx := vIdxStart
	hIdx := hIdxStart
	for wIdx := 0; wIdx < len(word); wIdx++ {
		if table[vIdx][hIdx] != word[wIdx] {
			return false
		}
		vIdx++
		hIdx++
	}

	return true
}

func checkDiaDownLeft(word string, hIdxStart, vIdxStart int, table []string) bool {

	vIdx := vIdxStart
	hIdx := hIdxStart
	for wIdx := 0; wIdx < len(word); wIdx++ {
		if table[vIdx][hIdx] != word[wIdx] {
			return false
		}
		vIdx++
		hIdx--
	}

	return true
}

//
// Part 2

func searchPattern(word string, table []string) (int, error) {

	if len(word)%2 != 1 {
		return 0, fmt.Errorf("the word needs to be an odd length")
	}

	var matchCount int

	halfWordMaxIdx := len(word) / 2

	vIdxMax := len(table) - 1
	hIdxMax := len(table[0]) - 1
	for vIdx := 0; vIdx <= vIdxMax; vIdx++ {
		for hIdx := 0; hIdx <= hIdxMax; hIdx++ {

			if table[vIdx][hIdx] != word[halfWordMaxIdx] {
				continue
			}

			hasSpaceUp := vIdx-halfWordMaxIdx >= 0
			hasSpaceDown := vIdx+halfWordMaxIdx <= vIdxMax
			hasSpaceLeft := hIdx-halfWordMaxIdx >= 0
			hasSpaceRight := hIdx+halfWordMaxIdx <= hIdxMax

			if hasSpaceUp && hasSpaceDown && hasSpaceLeft && hasSpaceRight {
				if checkPattern(word, hIdx, vIdx, table) {
					matchCount++
				}
			}

		}
	}

	return matchCount, nil
}

func checkPattern(word string, hIdxStart, vIdxStart int, table []string) bool {

	halfWordLength := len(word) / 2

	matchUpLeft := checkDiaUpLeft(word, hIdxStart+halfWordLength, vIdxStart+halfWordLength, table)
	matchDownRight := checkDiaDownRight(word, hIdxStart-halfWordLength, vIdxStart-halfWordLength, table)
	matchUpRight := checkDiaUpRight(word, hIdxStart-halfWordLength, vIdxStart+halfWordLength, table)
	matchDownLeft := checkDiaDownLeft(word, hIdxStart+halfWordLength, vIdxStart-halfWordLength, table)

	return (matchUpLeft || matchDownRight) && (matchUpRight || matchDownLeft)
}
