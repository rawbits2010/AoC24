package main

import (
	"fmt"
	"slices"

	"github.com/rawbits2010/AoC24/internal/inputhandler"
)

const empty = '.'

func main() {

	lines := inputhandler.ReadInput()

	field := Field{XMax: len(lines[0]) - 1, YMax: len(lines) - 1}
	antennas := extractAntennas(lines)

	resultP1 := countAntinodes(antennas, field)
	resultP2 := countAntinodesWithResonance(antennas, field)

	fmt.Printf("Result - Part 1: %d, Part 2: %d\n", resultP1, resultP2)
}

type Field struct {
	XMax, YMax int
}

type AntennaMap map[rune][]AntennaPosition

type AntennaPosition struct {
	Frequency  rune
	XPos, YPos int
}

type AntiNode struct {
	Antenna1   AntennaPosition
	Antenna2   AntennaPosition
	XPos, YPos int
	Harmonic   int
}

func extractAntennas(lines []string) AntennaMap {

	antennas := make(AntennaMap, 50)

	for yIdx, row := range lines {
		for xIdx, place := range row {

			if place == empty {
				continue
			}

			_, ok := antennas[place]
			if !ok {
				antennas[place] = make([]AntennaPosition, 0, 50)
			}

			antennas[place] = append(antennas[place],
				AntennaPosition{
					Frequency: place,
					XPos:      xIdx,
					YPos:      yIdx,
				})

		}
	}

	return antennas
}

func findAntinodes(antennas AntennaMap, filed Field) []AntiNode {

	antinodes := make([]AntiNode, 0, 100)
	for _, posList := range antennas {

		for currIdx := 0; currIdx < len(posList); currIdx++ {
			for targetIdx := 0; targetIdx < len(posList); targetIdx++ {

				if targetIdx == currIdx {
					continue
				}

				var antinodeX int
				if posList[currIdx].XPos < posList[targetIdx].XPos {

					xDiff := posList[targetIdx].XPos - posList[currIdx].XPos
					antinodeX = posList[currIdx].XPos - xDiff

				} else if posList[currIdx].XPos > posList[targetIdx].XPos {

					xDiff := posList[currIdx].XPos - posList[targetIdx].XPos
					antinodeX = posList[currIdx].XPos + xDiff

				} else {
					antinodeX = posList[currIdx].XPos
				}

				if antinodeX < 0 || antinodeX > filed.XMax {
					continue
				}

				var antinodeY int
				if posList[currIdx].YPos < posList[targetIdx].YPos {

					yDiff := posList[targetIdx].YPos - posList[currIdx].YPos
					antinodeY = posList[currIdx].YPos - yDiff

				} else if posList[currIdx].YPos > posList[targetIdx].YPos {

					yDiff := posList[currIdx].YPos - posList[targetIdx].YPos
					antinodeY = posList[currIdx].YPos + yDiff

				} else {
					antinodeY = posList[currIdx].YPos
				}

				if antinodeY < 0 || antinodeY > filed.YMax {
					continue
				}

				temp := AntiNode{
					Antenna1: posList[currIdx],
					Antenna2: posList[targetIdx],
					XPos:     antinodeX,
					YPos:     antinodeY,
				}

				antinodes = append(antinodes, temp)
			}
		}
	}

	return antinodes
}

func findAntinodesWithResonance(antennas AntennaMap, filed Field) []AntiNode {

	antinodes := make([]AntiNode, 0, 100)
	for _, posList := range antennas {

		for currIdx := 0; currIdx < len(posList); currIdx++ {
			for targetIdx := 0; targetIdx < len(posList); targetIdx++ {

				if targetIdx == currIdx {
					continue
				}

				harmonic := 0
				temp := AntiNode{
					Antenna1: posList[currIdx],
					Antenna2: posList[targetIdx],
					XPos:     posList[currIdx].XPos,
					YPos:     posList[currIdx].YPos,
					Harmonic: harmonic,
				}

				antinodes = append(antinodes, temp)

				xDiff := posList[targetIdx].XPos - posList[currIdx].XPos
				yDiff := posList[targetIdx].YPos - posList[currIdx].YPos

				prevPosX := posList[currIdx].XPos
				prevPosY := posList[currIdx].YPos
				for {
					harmonic++

					antinodeX := prevPosX - xDiff
					if antinodeX < 0 || antinodeX > filed.XMax {
						break
					}

					antinodeY := prevPosY - yDiff
					if antinodeY < 0 || antinodeY > filed.YMax {
						break
					}

					temp := AntiNode{
						Antenna1: posList[currIdx],
						Antenna2: posList[targetIdx],
						XPos:     antinodeX,
						YPos:     antinodeY,
						Harmonic: harmonic,
					}

					antinodes = append(antinodes, temp)

					prevPosX = antinodeX
					prevPosY = antinodeY
				}
			}
		}
	}

	return antinodes
}

func countUniqueAntinodes(antinodeList []AntiNode) int {

	slices.SortFunc(antinodeList,
		func(a, b AntiNode) int {

			if a.YPos < b.YPos {
				return -1
			} else if a.YPos > b.YPos {
				return 1
			} else {

				if a.XPos < b.XPos {
					return -1
				} else if a.XPos > b.XPos {
					return 1
				}
			}
			return 0
		})

	var countUnique int
	prev := AntiNode{
		XPos: -1,
		YPos: -1,
	}
	for _, antinode := range antinodeList {

		if prev.XPos == antinode.XPos && prev.YPos == antinode.YPos {
			continue
		}

		countUnique++

		prev.XPos = antinode.XPos
		prev.YPos = antinode.YPos
	}

	return countUnique
}

func countAntinodes(antennas AntennaMap, field Field) int {

	antinodeList := findAntinodes(antennas, field)
	result := countUniqueAntinodes(antinodeList)

	return result
}

func countAntinodesWithResonance(antennas AntennaMap, field Field) int {

	antinodeList := findAntinodesWithResonance(antennas, field)
	result := countUniqueAntinodes(antinodeList)

	return result
}
