package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/rawbits2010/AoC24/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()

	terrain, err := parseTerrain(lines)
	if err != nil {
		log.Fatal(err)
	}

	var resultP1 int
	for vIdx, row := range terrain.HeightMap {
		for hIdx, spot := range row {

			if spot.Height == 0 {
				terrain.Reset()
				resultP1 += countPaths(Position{X: hIdx, Y: vIdx}, terrain, true)
			}
		}
	}

	var resultP2 int
	for vIdx, row := range terrain.HeightMap {
		for hIdx, spot := range row {

			if spot.Height == 0 {
				terrain.Reset()
				resultP2 += countPaths(Position{X: hIdx, Y: vIdx}, terrain, false)
			}
		}
	}
	fmt.Printf("Result - Part 1: %d, Part 2: %d", resultP1, resultP2)
}

func parseTerrain(lines []string) (Field, error) {

	tempField := Field{
		Height: len(lines),
		Width:  len(lines[0]),
	}

	terrain := make([][]Spot, tempField.Height)
	for idx := 0; idx < len(terrain); idx++ {
		terrain[idx] = make([]Spot, tempField.Width)
	}

	for vIdx, line := range lines {
		for hIdx, spot := range line {

			val, err := strconv.Atoi(string(spot))
			if err != nil {
				return Field{}, fmt.Errorf("error parsing terrain at x: %d, y:%d", hIdx, vIdx)
			}

			terrain[vIdx][hIdx].Height = val
			terrain[vIdx][hIdx].isChecked = false
		}
	}

	tempField.HeightMap = terrain

	return tempField, nil
}

type Field struct {
	HeightMap     [][]Spot
	Width, Height int
}

func (f *Field) Reset() {
	for vIdx, row := range f.HeightMap {
		for hIdx, _ := range row {
			f.HeightMap[vIdx][hIdx].Reset()
		}
	}
}

type Spot struct {
	Height     int
	isChecked  bool
	isSelected bool
}

func (s *Spot) Reset() {
	s.isChecked = false
	s.isSelected = false
}

type Position struct {
	X, Y int
}

func countPaths(startPos Position, terrain Field, checkSpotOnce bool) int {

	var countTrail int
	toCheck := make([]Position, 0, 20)
	toCheck = append(toCheck, startPos)
	//fmt.Printf("Start - X: %d, Y: %d\n", startPos.X, startPos.Y)

	for {

		if len(toCheck) == 0 {
			break
		}

		//get one out
		currSpot := toCheck[0]
		toCheck = toCheck[1:]

		currPos := Position{X: currSpot.X, Y: currSpot.Y}
		currHeight := terrain.HeightMap[currPos.Y][currPos.X].Height

		terrain.HeightMap[currPos.Y][currPos.X].isChecked = true
		if terrain.HeightMap[currPos.Y][currPos.X].Height == 9 {
			countTrail++
			//fmt.Printf("End - X: %d, Y: %d\n", currPos.X, currPos.Y)
			continue
		}

		// left
		if currPos.X-1 >= 0 {
			if (checkSpotOnce && !terrain.HeightMap[currPos.Y][currPos.X-1].isChecked && !terrain.HeightMap[currPos.Y][currPos.X-1].isSelected) || !checkSpotOnce {
				if terrain.HeightMap[currPos.Y][currPos.X-1].Height == currHeight+1 {
					toCheck = append(toCheck, Position{X: currPos.X - 1, Y: currPos.Y})
					terrain.HeightMap[currPos.Y][currPos.X-1].isSelected = true
				}
			}
		}

		// right
		if currPos.X+1 < terrain.Width {
			if (checkSpotOnce && !terrain.HeightMap[currPos.Y][currPos.X+1].isChecked && !terrain.HeightMap[currPos.Y][currPos.X+1].isSelected) || !checkSpotOnce {
				if terrain.HeightMap[currPos.Y][currPos.X+1].Height == currHeight+1 {
					toCheck = append(toCheck, Position{X: currPos.X + 1, Y: currPos.Y})
					terrain.HeightMap[currPos.Y][currPos.X+1].isSelected = true
				}
			}
		}

		// up
		if currPos.Y-1 >= 0 {
			if (checkSpotOnce && !terrain.HeightMap[currPos.Y-1][currPos.X].isChecked && !terrain.HeightMap[currPos.Y-1][currPos.X].isSelected) || !checkSpotOnce {
				if terrain.HeightMap[currPos.Y-1][currPos.X].Height == currHeight+1 {
					toCheck = append(toCheck, Position{X: currPos.X, Y: currPos.Y - 1})
					terrain.HeightMap[currPos.Y-1][currPos.X].isSelected = true
				}
			}
		}

		// down
		if currPos.Y+1 < terrain.Height {
			if (checkSpotOnce && !terrain.HeightMap[currPos.Y+1][currPos.X].isChecked && !terrain.HeightMap[currPos.Y+1][currPos.X].isSelected) || !checkSpotOnce {
				if terrain.HeightMap[currPos.Y+1][currPos.X].Height == currHeight+1 {
					toCheck = append(toCheck, Position{X: currPos.X, Y: currPos.Y + 1})
					terrain.HeightMap[currPos.Y+1][currPos.X].isSelected = true
				}
			}
		}
	}

	return countTrail
}
