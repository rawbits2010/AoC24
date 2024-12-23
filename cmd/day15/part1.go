package main

import (
	"fmt"
	"log"
)

type Warehouse struct {
	Area     [][]rune
	Width    int
	Height   int
	RobotPos Position
}

func (w *Warehouse) DoMoves(movesList []string) {

	for _, line := range movesList {
		for _, direction := range line {
			switch direction {

			case UpMove:
				chainEndPos, isWalkable := w.searchForWallOrEmpty(StepUp)
				if isWalkable {
					w.pushBoxes(chainEndPos, StepUp)
				}

			case DownMove:
				chainEndPos, isWalkable := w.searchForWallOrEmpty(StepDown)
				if isWalkable {
					w.pushBoxes(chainEndPos, StepDown)
				}

			case LeftMove:
				chainEndPos, isWalkable := w.searchForWallOrEmpty(StepLeft)
				if isWalkable {
					w.pushBoxes(chainEndPos, StepLeft)
				}

			case RightMove:
				chainEndPos, isWalkable := w.searchForWallOrEmpty(StepRight)
				if isWalkable {
					w.pushBoxes(chainEndPos, StepRight)
				}
			}
		}
	}
}

func (w *Warehouse) searchForWallOrEmpty(step StepDir) (Position, bool) {

	for vIdx := w.RobotPos.Y + step.yDir; vIdx >= 0 && vIdx < w.Height; vIdx += step.yDir {
		for hIdx := w.RobotPos.X + step.xDir; hIdx >= 0 && hIdx < w.Width; hIdx += step.xDir {

			switch w.Area[vIdx][hIdx] {

			case WallSymbol:
				return Position{X: hIdx, Y: vIdx}, false

			case EmptySymbol:
				return Position{X: hIdx, Y: vIdx}, true
			}

			if step.xDir == 0 {
				break
			}
		}
		if step.yDir == 0 {
			break
		}
	}

	log.Fatal("out of bounds while serching for a wall or empty space")

	return Position{X: -1, Y: -1}, false
}

func (w *Warehouse) pushBoxes(startPos Position, step StepDir) {

	searchStep := StepDir{xDir: -step.xDir, yDir: -step.yDir}

	prevXIdx := startPos.X
	prevYIdx := startPos.Y
	for vIdx := startPos.Y + searchStep.yDir; vIdx >= 0 && vIdx < w.Height; vIdx += searchStep.yDir {
		for hIdx := startPos.X + searchStep.xDir; hIdx >= 0 && hIdx < w.Width; hIdx += searchStep.xDir {

			w.Area[prevYIdx][prevXIdx] = w.Area[vIdx][hIdx]

			if vIdx == w.RobotPos.Y && hIdx == w.RobotPos.X {

				w.Area[vIdx][hIdx] = EmptySymbol

				w.RobotPos.X = prevXIdx
				w.RobotPos.Y = prevYIdx

				return
			}

			prevXIdx = hIdx
			prevYIdx = vIdx

			if searchStep.xDir == 0 {
				break
			}
		}
		if searchStep.yDir == 0 {
			break
		}
	}

	log.Fatal("pushing boxes never reached the robot - telekinesis?")
}

func (w *Warehouse) SumBoxGPSCoords() int {

	var sum int
	for vIdx := 0; vIdx < w.Height; vIdx++ {
		for hIdx := 0; hIdx < w.Width; hIdx++ {
			if w.Area[vIdx][hIdx] == BoxSymbol {
				sum += (100 * vIdx) + hIdx
			}
		}
	}

	return sum
}

func (w *Warehouse) Display() {
	for _, row := range w.Area {
		fmt.Println(string(row))
	}
}
