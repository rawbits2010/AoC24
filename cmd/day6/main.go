package main

import (
	"fmt"

	"github.com/rawbits2010/AoC24/internal/inputhandler"
)

const guard = '^'
const obstacle = '#'
const empty = '.'

func main() {

	lines := inputhandler.ReadInput()

	field := NewField(lines)
	for {
		if isSuccess, _ := field.Step(); !isSuccess {
			break
		}
	}

	var resultP1 int
	for _, row := range field.Tiles {
		for _, tile := range row {
			if tile.IsTouched() {
				resultP1++
			}
		}
	}

	tempField := NewField(lines)
	var resultP2 int
	for vIdx, row := range field.Tiles {
		for hIdx, tile := range row {

			if tile.IsTouched() {

				tempField.PlaceTempObstacle(hIdx, vIdx)

				for {
					isSuccess, isLoop := tempField.Step()

					if isLoop {
						resultP2++
						break
					}

					if !isSuccess {
						break
					}
				}

				tempField.Reset()
			}
		}
	}

	fmt.Printf("Result - Part1: %d Part2: %d", resultP1, resultP2)
}
