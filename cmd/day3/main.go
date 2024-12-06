package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/rawbits2010/AoC24/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()

	// Part 1
	instStr, err := findInsatructions(lines, `mul\(\d{1,3},\d{1,3}\)`)
	if err != nil {
		log.Fatal(err)
	}

	instructions, err := parseInstructions(instStr)
	if err != nil {
		log.Fatal(err)
	}

	var resultP1 int
	for _, currInst := range instructions {
		resultP1 += currInst.Val1 * currInst.Val2
	}

	// Part2
	instStr, err = findInsatructions(lines, `(mul\(\d{1,3},\d{1,3}\))|(do\(\))|(don't\(\))`)
	if err != nil {
		log.Fatal(err)
	}

	instructions, err = parseInstructions(instStr)
	if err != nil {
		log.Fatal(err)
	}

	var resultP2 int
	dodo := true
	for _, currInst := range instructions {
		switch currInst.Op {
		case Do:
			dodo = true
		case Dont:
			dodo = false
		default:
			if dodo {
				resultP2 += currInst.Val1 * currInst.Val2
			}
		}
	}

	fmt.Printf("Result - Part1: %d Part2: %d", resultP1, resultP2)
}

func findInsatructions(lines []string, expr string) ([]string, error) {

	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, fmt.Errorf("error finding instructions: %w", err)
	}

	instStrList := make([]string, 0, 1000)

	for _, line := range lines {
		temp := re.FindAllString(line, -1)
		if temp == nil {
			continue
		}
		instStrList = append(instStrList, temp...)
	}

	return instStrList, nil
}

func parseInstructions(instStrList []string) ([]Instruction, error) {

	re, err := regexp.Compile(`\d{1,3}`)
	if err != nil {
		return nil, fmt.Errorf("error finding values: %w", err)
	}

	instrList := make([]Instruction, len(instStrList))

	for instIdx, instStr := range instStrList {

		if !strings.HasPrefix(instStr, Multiplication) {

			if strings.HasPrefix(instStr, Dont) {
				instrList[instIdx].Op = Dont
			} else {
				instrList[instIdx].Op = Do
			}

			continue
		}

		temp := re.FindAllString(instStr, -1)
		if len(temp) != 2 {
			return nil, fmt.Errorf("invalid number of operands '%s'", instStr)
		}

		val1, err := strconv.Atoi(temp[0])
		if err != nil {
			return nil, fmt.Errorf("invalid first operand in '%s'", instStr)
		}

		val2, err := strconv.Atoi(temp[1])
		if err != nil {
			return nil, fmt.Errorf("invalid second operand in '%s'", instStr)
		}

		instrList[instIdx].Op = Multiplication
		instrList[instIdx].Val1 = val1
		instrList[instIdx].Val2 = val2
	}

	return instrList, nil
}

type Operation string

const (
	Multiplication = "mul"
	Do             = "do"
	Dont           = "don't"
)

type Instruction struct {
	Op   Operation
	Val1 int
	Val2 int
}
