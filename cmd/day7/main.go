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

	equations, err := parseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	var resultP1 int
	for eqIdx, equation := range equations {

		counter := NewBinaryCounter(len(equation.Numbers) - 1)
		for {

			tempVal, err := equation.Eval(*counter)
			if err == nil {

				if tempVal == equation.TestValue {
					resultP1 += equation.TestValue
					equations[eqIdx].IsDone = true
					break
				}
			}

			counter.Increment()
			if counter.IsZero() {
				break
			}
		}
	}

	// NOTE: Part 2 just desperately needs multithreading
	resultP2 := resultP1
nextEquation:
	for eqIdx, equation := range equations {

		if equation.IsDone {
			continue
		}

		counter := NewBinaryCounter(len(equation.Numbers) - 1)
		for {

			concatCounter := NewBinaryCounter(len(equation.Numbers) - 1)
			for {
				tempVal, err := equation.EvalWithConcat(*counter, *concatCounter)
				if err == nil {

					if tempVal == equation.TestValue {
						resultP2 += equation.TestValue
						equations[eqIdx].IsDone = true
						continue nextEquation
					}
				}

				concatCounter.Increment()
				if concatCounter.IsZero() {
					break
				}
			}

			counter.Increment()
			if counter.IsZero() {
				break
			}
		}
	}

	fmt.Printf("Result - Part 1: %d, Part 2: %d\n", resultP1, resultP2)
}

type BinaryCounter struct {
	Bits       []bool
	IsNegative bool
}

func NewBinaryCounter(length int) *BinaryCounter {
	return &BinaryCounter{
		Bits:       make([]bool, length),
		IsNegative: false,
	}
}

func (bc *BinaryCounter) Increment() {

	for idx := 0; idx < len(bc.Bits); idx++ {

		bc.Bits[idx] = !bc.Bits[idx]

		if bc.IsNegative {
			if !bc.Bits[idx] {
				break
			}
		} else {
			if bc.Bits[idx] {
				break
			}
		}
	}

	if bc.IsAllTrue() {
		bc.IsNegative = true
	}
}

func (bc *BinaryCounter) IsZero() bool {
	for _, bit := range bc.Bits {
		if bit {
			return false
		}
	}
	return true
}

func (bc *BinaryCounter) IsAllTrue() bool {
	for _, bit := range bc.Bits {
		if !bit {
			return false
		}
	}
	return true
}

type Equation struct {
	TestValue int
	Numbers   []int
	IsDone    bool
}

// Part 1
func (e *Equation) Eval(counter BinaryCounter) (int, error) {

	tempVal := e.Numbers[0]
	var err error
	for bIdx := 0; bIdx < len(counter.Bits); bIdx++ {

		if counter.Bits[bIdx] {
			tempVal, err = mulInt(tempVal, e.Numbers[bIdx+1])
			if err != nil {
				return 0, err
			}
		} else {
			tempVal, err = addInt(tempVal, e.Numbers[bIdx+1])
			if err != nil {
				return 0, err
			}
		}

		if tempVal > e.TestValue {
			return 0, fmt.Errorf("overshot")
		}
	}

	return tempVal, nil
}

// Part 2
func (e *Equation) EvalWithConcat(counter BinaryCounter, concatCounter BinaryCounter) (int, error) {

	tempVal := e.Numbers[0]
	var err error
	for bIdx := 0; bIdx < len(counter.Bits); bIdx++ {

		if concatCounter.Bits[bIdx] {

			concatTemp := fmt.Sprintf("%d%d", tempVal, e.Numbers[bIdx+1])
			tempVal, err = strconv.Atoi(concatTemp)
			if err != nil {
				log.Fatalf("error converting concat value '%s' - %s", concatTemp, err)
			}
		} else {

			if counter.Bits[bIdx] {
				tempVal, err = mulInt(tempVal, e.Numbers[bIdx+1])
				if err != nil {
					return 0, err
				}
			} else {
				tempVal, err = addInt(tempVal, e.Numbers[bIdx+1])
				if err != nil {
					return 0, err
				}
			}
		}

		if tempVal > e.TestValue {
			return 0, fmt.Errorf("overshot")
		}
	}

	return tempVal, nil
}

func parseInput(lines []string) ([]Equation, error) {

	equations := make([]Equation, len(lines))
	for lineIdx, line := range lines {

		temp := strings.Split(line, ": ")
		if len(temp) != 2 {
			return nil, fmt.Errorf("invalid equation format in '%s'", line)
		}

		testVal, err := strconv.Atoi(temp[0])
		if err != nil {
			return nil, fmt.Errorf("error converting test value '%s' at line %d", temp[0], lineIdx)
		}

		tempNums := strings.Split(temp[1], " ")
		if len(tempNums) < 2 {
			// this is not in the text
			return nil, fmt.Errorf("not enough numbers in '%s'", line)
		}

		numbers := make([]int, len(tempNums))
		for tnIdx, tnStr := range tempNums {
			tn, err := strconv.Atoi(tnStr)
			if err != nil {
				return nil, fmt.Errorf("error converting number '%s' at line %d", tnStr, lineIdx)
			}
			numbers[tnIdx] = tn
		}

		equations[lineIdx] = Equation{
			TestValue: testVal,
			Numbers:   numbers,
		}
	}

	return equations, nil
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
