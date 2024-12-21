package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/rawbits2010/AoC24/internal/inputhandler"
)

const ButtonACost = 3
const ButtonBCost = 1

func main() {

	lines := inputhandler.ReadInput()

	machines, err := parseMachines(lines)
	if err != nil {
		log.Fatal(err)
	}

	var resultP1 int
	for _, machine := range machines {
		if cost, grabbed := grabPrizeWithBruteforce(machine); grabbed {
			resultP1 += cost
		}
	}

	fmt.Printf("Result - Part 1: %d", resultP1)
}

type Position struct {
	X, Y int
}

type Button struct {
	X, Y int
}

type ClawMachine struct {
	A, B     Button
	PrizePos Position
}

func parseMachines(lines []string) ([]ClawMachine, error) {

	machines := make([]ClawMachine, (len(lines)/4)+1)

	var machineIdx int
	for _, line := range lines {
		if len(line) == 0 {
			machineIdx++
			continue
		}

		token := strings.Fields(line)
		if len(token) != 3 && len(token) != 4 {
			return nil, fmt.Errorf("invalid number of tokens '%d' in line '%s'", len(token), line)
		}

		if token[0] == "Button" {

			xToken := strings.FieldsFunc(token[2],
				func(r rune) bool {
					return r == '+' || r == ','
				})
			if len(xToken) != 2 {
				return nil, fmt.Errorf("invalid X value format '%s' in line '%s'", token[2], line)
			}

			xVal, err := strconv.Atoi(xToken[1])
			if err != nil {
				return nil, fmt.Errorf("invalid X value '%s' in line '%s': %w", xToken[1], line, err)
			}

			yToken := strings.FieldsFunc(token[3],
				func(r rune) bool {
					return r == '+'
				})
			if len(yToken) != 2 {
				return nil, fmt.Errorf("invalid Y value format '%s' in line '%s'", token[3], line)
			}

			yVal, err := strconv.Atoi(yToken[1])
			if err != nil {
				return nil, fmt.Errorf("invalid Y value '%s' in line '%s': %w", yToken[1], line, err)
			}

			if token[1] == "A:" {
				machines[machineIdx].A = Button{X: xVal, Y: yVal}
			} else if token[1] == "B:" {
				machines[machineIdx].B = Button{X: xVal, Y: yVal}
			} else {
				return nil, fmt.Errorf("unknown button '%s' in line '%s'", token[1], line)
			}

		} else if token[0] == "Prize:" {

			xToken := strings.FieldsFunc(token[1],
				func(r rune) bool {
					return r == '=' || r == ','
				})
			if len(xToken) != 2 {
				return nil, fmt.Errorf("invalid X value format '%s' in line '%s'", token[1], line)
			}

			xVal, err := strconv.Atoi(xToken[1])
			if err != nil {
				return nil, fmt.Errorf("invalid X value '%s' in line '%s': %w", xToken[1], line, err)
			}

			yToken := strings.FieldsFunc(token[2],
				func(r rune) bool {
					return r == '='
				})
			if len(yToken) != 2 {
				return nil, fmt.Errorf("invalid Y value format '%s' in line '%s'", token[2], line)
			}

			yVal, err := strconv.Atoi(yToken[1])
			if err != nil {
				return nil, fmt.Errorf("invalid Y value '%s' in line '%s': %w", yToken[1], line, err)
			}

			machines[machineIdx].PrizePos = Position{X: xVal, Y: yVal}
		} else {
			return nil, fmt.Errorf("unknown feature '%s' in line '%s'", token[0], line)
		}
	}

	return machines, nil
}

func grabPrizeWithBruteforce(machine ClawMachine) (int, bool) {

startPressingButtons:
	for currPressA := 0; currPressA <= 100; currPressA++ {
		for currPressB := 0; currPressB <= 100; currPressB++ {

			currXPos := (currPressA * machine.A.X) + (currPressB * machine.B.X)

			if currXPos > machine.PrizePos.X {
				if currPressB == 0 {
					return -1, false
				}
				continue startPressingButtons
			}

			if currXPos == machine.PrizePos.X {

				currYPos := (currPressA * machine.A.Y) + (currPressB * machine.B.Y)

				if currYPos == machine.PrizePos.Y {
					cost := currPressA*ButtonACost + currPressB*ButtonBCost
					return cost, true
				}
			}
		}
	}

	return -1, false
}

func grabPrizeWithThePowerOf10thGradeMath(machine ClawMachine) (int, bool) {

	return -1, false
}
