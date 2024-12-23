package main

import (
	"fmt"
	"log"
)

type WideWarehouse struct {
	Area     [][]rune
	Width    int
	Height   int
	RobotPos Position
}

func (w *WideWarehouse) DoMoves(movesList []string) {

	for _, line := range movesList {
		for _, direction := range line {

			// visualize steps
			w.Display()
			fmt.Println(string(direction))

			switch direction {

			case UpMove:
				boxLayers, isWalkable := w.searchForWallOrEmpty(StepUp)
				if isWalkable {
					w.pushBoxes(boxLayers, StepUp)
				}

			case DownMove:
				boxLayers, isWalkable := w.searchForWallOrEmpty(StepDown)
				if isWalkable {
					w.pushBoxes(boxLayers, StepDown)
				}

			case LeftMove:
				boxLayers, isWalkable := w.searchForWallOrEmpty(StepLeft)
				if isWalkable {
					w.pushBoxes(boxLayers, StepLeft)
				}

			case RightMove:
				boxLayers, isWalkable := w.searchForWallOrEmpty(StepRight)
				if isWalkable {
					w.pushBoxes(boxLayers, StepRight)
				}
			}
		}
	}
}

func (w *WideWarehouse) searchForWallOrEmpty(step StepDir) ([][]Position, bool) {

	// sanity check
	if step.xDir == 0 && step.yDir == 0 {
		return [][]Position{
			{
				{X: w.RobotPos.X, Y: w.RobotPos.Y},
			},
		}, false
	}

	boxLayers := make([][]Position, 0, 10)

	toCheck := make([]Position, 0, 20)
	toCheck = append(toCheck, w.RobotPos)

	boxLayers = append(boxLayers, toCheck)

	boxLayers, isWalkable := w.searchForWallOrEmptyRecurse(step, toCheck, boxLayers)

	return boxLayers, isWalkable
}

func (w *WideWarehouse) searchForWallOrEmptyRecurse(step StepDir, toCheck []Position, boxLayers [][]Position) ([][]Position, bool) {

	nextToCheck := make([]Position, 0, 20)

	for _, currPos := range toCheck {

		posToCheck := Position{X: currPos.X + step.xDir, Y: currPos.Y + step.yDir}

		switch w.Area[posToCheck.Y][posToCheck.X] {

		case WallSymbol:
			return boxLayers, false

		case EmptySymbol:
			// check the rest

		case BoxLeftSymbol:

			if step.xDir < 0 {
				log.Fatalf("found left side of a box while checking to the left at x: %d, y: %d", posToCheck.X, posToCheck.Y)
			} else if step.xDir > 0 {
				nextToCheck = append(nextToCheck, Position{X: posToCheck.X + 1, Y: posToCheck.Y})
			}

			if step.yDir != 0 {
				nextToCheck = append(nextToCheck, Position{X: posToCheck.X, Y: posToCheck.Y})
				if w.Area[currPos.Y][currPos.X] != w.Area[posToCheck.Y][posToCheck.X] {
					nextToCheck = append(nextToCheck, Position{X: posToCheck.X + 1, Y: posToCheck.Y})
				}
			}

		case BoxRightSymbol:

			if step.xDir > 0 {
				log.Fatalf("found right side of a box while checking to the right at x: %d, x: %d", posToCheck.X, posToCheck.Y)
			} else if step.xDir < 0 {
				nextToCheck = append(nextToCheck, Position{X: posToCheck.X - 1, Y: posToCheck.Y})
			}

			if step.yDir != 0 {
				if w.Area[currPos.Y][currPos.X] != w.Area[posToCheck.Y][posToCheck.X] {
					nextToCheck = append(nextToCheck, Position{X: posToCheck.X - 1, Y: posToCheck.Y})
				}
				nextToCheck = append(nextToCheck, Position{X: posToCheck.X, Y: posToCheck.Y})
			}

		}
	}

	if len(nextToCheck) == 0 {
		return boxLayers, true
	}

	boxLayers = append(boxLayers, nextToCheck)

	return w.searchForWallOrEmptyRecurse(step, nextToCheck, boxLayers)
}

func (w *WideWarehouse) pushBoxes(boxLayers [][]Position, step StepDir) {

	for lIdx := len(boxLayers) - 1; lIdx >= 0; lIdx-- {

		for _, currPos := range boxLayers[lIdx] {

			currSpot := w.Area[currPos.Y][currPos.X]

			if step.yDir != 0 {
				w.Area[currPos.Y+step.yDir][currPos.X] = currSpot
				w.Area[currPos.Y][currPos.X] = EmptySymbol
				continue
			}

			if step.xDir != 0 {
				switch w.Area[currPos.Y][currPos.X] {

				case BoxLeftSymbol:
					w.Area[currPos.Y][currPos.X+step.xDir] = currSpot
					w.Area[currPos.Y][currPos.X] = BoxRightSymbol
					w.Area[currPos.Y][currPos.X-step.xDir] = EmptySymbol

				case BoxRightSymbol:
					w.Area[currPos.Y][currPos.X+step.xDir] = currSpot
					w.Area[currPos.Y][currPos.X] = BoxLeftSymbol
					w.Area[currPos.Y][currPos.X-step.xDir] = EmptySymbol

				case RobotSymbol:
					w.Area[currPos.Y][currPos.X+step.xDir] = currSpot
					w.Area[currPos.Y][currPos.X] = EmptySymbol
				}
			}
		}
	}

	w.RobotPos.X += step.xDir
	w.RobotPos.Y += step.yDir
}

func (w *WideWarehouse) SumBoxGPSCoords() int {

	var sum int
	for vIdx := 0; vIdx < w.Height; vIdx++ {
		for hIdx := 0; hIdx < w.Width; hIdx++ {
			if w.Area[vIdx][hIdx] == BoxLeftSymbol {
				sum += (100 * vIdx) + hIdx
			}
		}
	}

	return sum
}

func (w *WideWarehouse) Display() {
	for _, row := range w.Area {
		fmt.Println(string(row))
	}
}
