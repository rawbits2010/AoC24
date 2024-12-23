package main

import (
	"fmt"
	"log"

	"github.com/rawbits2010/AoC24/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()
	warehouse, wideWarehouse, movesList := parseWarehouse(lines)

	var boxCount int
	for _, row := range wideWarehouse.Area {
		for _, symbol := range row {
			if symbol == BoxLeftSymbol {
				boxCount++
			}
		}
	}
	fmt.Println(boxCount)

	warehouse.DoMoves(movesList)

	boxCount = 0
	for _, row := range wideWarehouse.Area {
		for _, symbol := range row {
			if symbol == BoxLeftSymbol {
				boxCount++
			}
		}
	}
	fmt.Println(boxCount)

	resultP1 := warehouse.SumBoxGPSCoords()
	warehouse.Display()

	wideWarehouse.DoMoves(movesList)
	resultP2 := wideWarehouse.SumBoxGPSCoords()
	wideWarehouse.Display()

	fmt.Printf("Result - Part 1: %d, Part 2: %d", resultP1, resultP2)
}

const RobotSymbol = '@'
const WallSymbol = '#'
const BoxSymbol = 'O'
const EmptySymbol = '.'

const BoxLeftSymbol = '['
const BoxRightSymbol = ']'

const UpMove = '^'
const DownMove = 'v'
const LeftMove = '<'
const RightMove = '>'

type StepDir struct {
	xDir, yDir int
}

var StepUp = StepDir{xDir: 0, yDir: -1}
var StepDown = StepDir{xDir: 0, yDir: 1}
var StepLeft = StepDir{xDir: -1, yDir: 0}
var StepRight = StepDir{xDir: 1, yDir: 0}

type Position struct {
	X, Y int
}

func parseWarehouse(lines []string) (Warehouse, WideWarehouse, []string) {

	layout := make([][]rune, 0, len(lines))
	wideLayout := make([][]rune, 0, len(lines))
	var movesList []string

	for lIdx, line := range lines {
		if len(line) == 0 {

			for _, line := range lines[:lIdx] {

				layout = append(layout, []rune(line))

				tempWide := make([]rune, len(line)*2)
				for lIdx, spot := range line {

					switch spot {

					case WallSymbol:
						fallthrough
					case EmptySymbol:
						tempWide[lIdx*2] = spot
						tempWide[(lIdx*2)+1] = spot
						continue

					case RobotSymbol:
						tempWide[lIdx*2] = spot
						tempWide[(lIdx*2)+1] = EmptySymbol
						continue

					case BoxSymbol:
						tempWide[lIdx*2] = BoxLeftSymbol
						tempWide[(lIdx*2)+1] = BoxRightSymbol
						continue

					default:
						log.Fatalf("unknown symbol found while parsing warehouse layout '%c'", spot)
					}
				}

				wideLayout = append(wideLayout, tempWide)
			}

			movesList = lines[lIdx+1:]
			break
		}
	}

	var rXIdx, rYIdx int
	for vIdx, row := range layout {
		for hIdx, spot := range row {
			if spot == RobotSymbol {
				rXIdx = hIdx
				rYIdx = vIdx
				break
			}
		}
	}

	warehouse := Warehouse{
		Area:     layout,
		Width:    len(layout[0]),
		Height:   len(layout),
		RobotPos: Position{X: rXIdx, Y: rYIdx},
	}

	wideWarehouse := WideWarehouse{
		Area:     wideLayout,
		Width:    len(wideLayout[0]),
		Height:   len(wideLayout),
		RobotPos: Position{X: rXIdx * 2, Y: rYIdx},
	}

	return warehouse, wideWarehouse, movesList
}
