package main

import (
	"fmt"

	"github.com/rawbits2010/AoC24/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()
	garden := parsePlot(lines)

	trackFences := make(map[int]Fences)

	var resultP1 int
	var regionId int
	for vIdx, row := range garden.Plot {
		for hIdx, plant := range row {

			if plant.IsChecked {
				continue
			}

			area, perimeter := findRegion(regionId, Position{X: hIdx, Y: vIdx}, garden)
			trackFences[regionId] = Fences{Area: area}

			//fmt.Printf("plant: %c, area: %d, perimeter: %d\n", plant.PlantType, area, perimeter)
			resultP1 += area * perimeter

			garden.Reset()
			regionId++
		}
	}

	findRangeSides(garden, trackFences)
	var resultP2 int
	for _, fences := range trackFences {

		sides := fences.BottomCount + fences.TopCount + fences.LeftCount + fences.RightCount
		//fmt.Printf("range: %d, sides: %d, area: %d\n", rangeId, sides, fences.Area)
		resultP2 += sides * fences.Area
	}

	fmt.Printf("Result - Part 1: %d, Part 2: %d\n", resultP1, resultP2)
}

type Fences struct {
	Area                                                 int
	OngoingTop, OngoingBottom, OngoingLeft, OngoingRight bool
	TopCount, BottomCount, LeftCount, RightCount         int
}

type Position struct {
	X, Y int
}

type Plant struct {
	PlantType   rune
	IsSelected  bool
	IsChecked   bool
	RegionId    int
	FenceTop    bool
	FenceBottom bool
	FenceLeft   bool
	FenceRight  bool
}

func (p *Plant) Reset() {
	p.IsSelected = false
}

type Garden struct {
	Plot          [][]Plant
	Width, Height int
}

func (g *Garden) Reset() {
	for vIdx := 0; vIdx < g.Height; vIdx++ {
		for hIdx := 0; hIdx < g.Width; hIdx++ {
			g.Plot[vIdx][hIdx].Reset()
		}
	}
}

func parsePlot(lines []string) Garden {

	garden := Garden{
		Height: len(lines),
		Width:  len(lines[0]),
	}

	plot := make([][]Plant, garden.Height)
	for idx := 0; idx < len(plot); idx++ {
		plot[idx] = make([]Plant, garden.Width)
	}

	for vIdx, row := range lines {
		for hIdx, plant := range row {
			plot[vIdx][hIdx].PlantType = plant
		}
	}

	garden.Plot = plot

	return garden
}

func findRegion(regionId int, startPos Position, garden Garden) (int, int) {

	var countArea int
	var countOuterEdge int
	toCheck := make([]Position, 0, 20)
	toCheck = append(toCheck, startPos)

	for {

		if len(toCheck) == 0 {
			break
		}

		//get one out
		currSpot := toCheck[0]
		toCheck = toCheck[1:]

		currPos := Position{X: currSpot.X, Y: currSpot.Y}
		currPlant := garden.Plot[currPos.Y][currPos.X].PlantType

		garden.Plot[currPos.Y][currPos.X].IsChecked = true
		garden.Plot[currPos.Y][currPos.X].RegionId = regionId

		countArea++

		// left
		if currPos.X-1 >= 0 {
			if garden.Plot[currPos.Y][currPos.X-1].PlantType == currPlant {
				if !garden.Plot[currPos.Y][currPos.X-1].IsChecked && !garden.Plot[currPos.Y][currPos.X-1].IsSelected {

					toCheck = append(toCheck, Position{X: currPos.X - 1, Y: currPos.Y})
					garden.Plot[currPos.Y][currPos.X-1].IsSelected = true
				}
			} else {
				garden.Plot[currPos.Y][currPos.X].FenceLeft = true
				countOuterEdge++
			}
		} else {
			garden.Plot[currPos.Y][currPos.X].FenceLeft = true
			countOuterEdge++
		}

		// right
		if currPos.X+1 < garden.Width {
			if garden.Plot[currPos.Y][currPos.X+1].PlantType == currPlant {
				if !garden.Plot[currPos.Y][currPos.X+1].IsChecked && !garden.Plot[currPos.Y][currPos.X+1].IsSelected {

					toCheck = append(toCheck, Position{X: currPos.X + 1, Y: currPos.Y})
					garden.Plot[currPos.Y][currPos.X+1].IsSelected = true
				}
			} else {
				garden.Plot[currPos.Y][currPos.X].FenceRight = true
				countOuterEdge++
			}
		} else {
			garden.Plot[currPos.Y][currPos.X].FenceRight = true
			countOuterEdge++
		}

		// up
		if currPos.Y-1 >= 0 {
			if garden.Plot[currPos.Y-1][currPos.X].PlantType == currPlant {
				if !garden.Plot[currPos.Y-1][currPos.X].IsChecked && !garden.Plot[currPos.Y-1][currPos.X].IsSelected {

					toCheck = append(toCheck, Position{X: currPos.X, Y: currPos.Y - 1})
					garden.Plot[currPos.Y-1][currPos.X].IsSelected = true
				}
			} else {
				garden.Plot[currPos.Y][currPos.X].FenceTop = true
				countOuterEdge++
			}
		} else {
			garden.Plot[currPos.Y][currPos.X].FenceTop = true
			countOuterEdge++
		}

		// down
		if currPos.Y+1 < garden.Height {
			if garden.Plot[currPos.Y+1][currPos.X].PlantType == currPlant {
				if !garden.Plot[currPos.Y+1][currPos.X].IsChecked && !garden.Plot[currPos.Y+1][currPos.X].IsSelected {

					toCheck = append(toCheck, Position{X: currPos.X, Y: currPos.Y + 1})
					garden.Plot[currPos.Y+1][currPos.X].IsSelected = true
				}
			} else {
				garden.Plot[currPos.Y][currPos.X].FenceBottom = true
				countOuterEdge++
			}
		} else {
			garden.Plot[currPos.Y][currPos.X].FenceBottom = true
			countOuterEdge++
		}
	}

	return countArea, countOuterEdge
}

func findRangeSides(garden Garden, trackFences map[int]Fences) {

	for _, row := range garden.Plot {
		prevRegion := -1
		for hIdx, plant := range row {

			if prevRegion != -1 && plant.RegionId != prevRegion {

				if trackFences[prevRegion].OngoingTop {
					if entry, ok := trackFences[prevRegion]; ok {
						entry.TopCount++
						entry.OngoingTop = false
						trackFences[prevRegion] = entry
					}
				}

				if trackFences[prevRegion].OngoingBottom {
					if entry, ok := trackFences[prevRegion]; ok {
						entry.BottomCount++
						entry.OngoingBottom = false
						trackFences[prevRegion] = entry
					}
				}
			}

			if trackFences[plant.RegionId].OngoingTop {

				if !plant.FenceTop {
					if entry, ok := trackFences[plant.RegionId]; ok {
						entry.TopCount++
						entry.OngoingTop = false
						trackFences[plant.RegionId] = entry
					}
				}
			} else {

				if plant.FenceTop {
					if entry, ok := trackFences[plant.RegionId]; ok {
						entry.OngoingTop = true
						trackFences[plant.RegionId] = entry
					}
				}
			}

			if trackFences[plant.RegionId].OngoingBottom {

				if !plant.FenceBottom {
					if entry, ok := trackFences[plant.RegionId]; ok {
						entry.BottomCount++
						entry.OngoingBottom = false
						trackFences[plant.RegionId] = entry
					}
				}
			} else {

				if plant.FenceBottom {
					if entry, ok := trackFences[plant.RegionId]; ok {
						entry.OngoingBottom = true
						trackFences[plant.RegionId] = entry
					}
				}
			}

			if hIdx+1 == garden.Width {

				if trackFences[plant.RegionId].OngoingTop {
					if entry, ok := trackFences[plant.RegionId]; ok {
						entry.TopCount++
						entry.OngoingTop = false
						trackFences[plant.RegionId] = entry
					}
				}

				if trackFences[plant.RegionId].OngoingBottom {
					if entry, ok := trackFences[plant.RegionId]; ok {
						entry.BottomCount++
						entry.OngoingBottom = false
						trackFences[plant.RegionId] = entry
					}
				}
			}

			prevRegion = plant.RegionId
		}
	}

	for hIdx := 0; hIdx < garden.Height; hIdx++ {
		prevRegion := -1
		for vIdx := 0; vIdx < garden.Width; vIdx++ {
			plant := garden.Plot[vIdx][hIdx]

			if prevRegion != -1 && plant.RegionId != prevRegion {

				if trackFences[prevRegion].OngoingLeft {
					if entry, ok := trackFences[prevRegion]; ok {
						entry.LeftCount++
						entry.OngoingLeft = false
						trackFences[prevRegion] = entry
					}
				}

				if trackFences[prevRegion].OngoingRight {
					if entry, ok := trackFences[prevRegion]; ok {
						entry.RightCount++
						entry.OngoingRight = false
						trackFences[prevRegion] = entry
					}
				}
			}

			if trackFences[plant.RegionId].OngoingLeft {

				if !plant.FenceLeft {
					if entry, ok := trackFences[plant.RegionId]; ok {
						entry.LeftCount++
						entry.OngoingLeft = false
						trackFences[plant.RegionId] = entry
					}
				}
			} else {

				if plant.FenceLeft {
					if entry, ok := trackFences[plant.RegionId]; ok {
						entry.OngoingLeft = true
						trackFences[plant.RegionId] = entry
					}
				}
			}

			if trackFences[plant.RegionId].OngoingRight {

				if !plant.FenceRight {
					if entry, ok := trackFences[plant.RegionId]; ok {
						entry.RightCount++
						entry.OngoingRight = false
						trackFences[plant.RegionId] = entry
					}
				}
			} else {

				if plant.FenceRight {
					if entry, ok := trackFences[plant.RegionId]; ok {
						entry.OngoingRight = true
						trackFences[plant.RegionId] = entry
					}
				}
			}

			if vIdx+1 == garden.Height {

				if trackFences[plant.RegionId].OngoingLeft {
					if entry, ok := trackFences[plant.RegionId]; ok {
						entry.LeftCount++
						entry.OngoingLeft = false
						trackFences[plant.RegionId] = entry
					}
				}

				if trackFences[plant.RegionId].OngoingRight {
					if entry, ok := trackFences[plant.RegionId]; ok {
						entry.RightCount++
						entry.OngoingRight = false
						trackFences[plant.RegionId] = entry
					}
				}
			}

			prevRegion = plant.RegionId
		}
	}
}
