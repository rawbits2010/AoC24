package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rawbits2010/AoC24/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()
	robots, err := parseRobots(lines)
	if err != nil {
		log.Fatal(err)
	}

	outputDirName := "out"

	fieldWidth := 101
	fieldHeight := 103
	//fieldWidth := 11
	//fieldHeight := 7
	timeInSec := 100

	resultP1 := solvePart1(robots, fieldWidth, fieldHeight, timeInSec)

	resetRobots(robots)
	constellations := solvePart2(robots, fieldWidth, fieldHeight)

	ensureOutputDir(outputDirName)
	maxRobotCount := 0
	maxRobotCountAtIdx := -1
	for cIdx, constellation := range constellations {

		if constellation.RobotCount >= maxRobotCount {
			maxRobotCount = constellation.RobotCount
			maxRobotCountAtIdx = cIdx
		}

		filename := strconv.Itoa(constellation.AtTime)
		filename += ".txt"

		filepath := filepath.Join(outputDirName, filename)
		err := saveSnapshotToFile(filepath, constellation.Snapshot, fieldWidth, fieldHeight)
		if err != nil {
			continue
		}
	}
	resultP2 := constellations[maxRobotCountAtIdx].AtTime

	displayArea(constellations[maxRobotCountAtIdx].Snapshot, fieldWidth, fieldHeight)

	fmt.Printf("Result - Part 1: %d\n", resultP1)
	fmt.Printf("Possible result - Part 2: %d (all suspected solutions are in the '%s' directory)\n", resultP2, outputDirName)
}

func solvePart1(robots []Robots, fieldWidth, fieldHeight int, timeInSec int) int {

	for tIdx := 1; tIdx <= timeInSec; tIdx++ {

		for rIdx := 0; rIdx < len(robots); rIdx++ {

			robots[rIdx].CurrPosX += robots[rIdx].VelocityX
			if robots[rIdx].CurrPosX < 0 {
				robots[rIdx].CurrPosX += fieldWidth
			}
			robots[rIdx].CurrPosX %= fieldWidth

			robots[rIdx].CurrPosY += robots[rIdx].VelocityY
			if robots[rIdx].CurrPosY < 0 {
				robots[rIdx].CurrPosY += fieldHeight
			}
			robots[rIdx].CurrPosY %= fieldHeight
		}
	}

	var q1, q2, q3, q4 int
	fieldMaxXIdx := fieldWidth - 1
	fieldMaxYIdx := fieldHeight - 1
	for _, robot := range robots {

		currX := robot.CurrPosX
		currY := robot.CurrPosY

		if currX < fieldMaxXIdx/2 {
			if currY < fieldMaxYIdx/2 {
				q1++
			} else if currY >= (fieldMaxYIdx/2)+1 {
				q3++
			}
		} else if currX >= (fieldMaxXIdx/2)+1 {
			if currY < fieldMaxYIdx/2 {
				q2++
			} else if currY >= (fieldMaxYIdx/2)+1 {
				q4++
			}
		}
	}

	return q1 * q2 * q3 * q4
}

type Robots struct {
	StartX, StartY       int
	CurrPosX, CurrPosY   int
	VelocityX, VelocityY int
	IsChecked            bool
}

func parseRobots(lines []string) ([]Robots, error) {

	robots := make([]Robots, len(lines))

	for lIdx, line := range lines {

		fields := strings.FieldsFunc(line,
			func(r rune) bool {
				return r == '=' || r == ',' || r == ' '
			})
		if len(fields) != 6 {
			return nil, fmt.Errorf("invalid robot data format '%s'", line)
		}

		pxVal, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("error converting robot position x '%s' in line '%s': %w", fields[1], line, err)
		}

		pyVal, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("error converting robot position y '%s' in line '%s': %w", fields[2], line, err)
		}

		vxVal, err := strconv.Atoi(fields[4])
		if err != nil {
			return nil, fmt.Errorf("error converting robot velocity x '%s' in line '%s': %w", fields[4], line, err)
		}

		vyVal, err := strconv.Atoi(fields[5])
		if err != nil {
			return nil, fmt.Errorf("error converting robot velocity x '%s' in line '%s': %w", fields[5], line, err)
		}

		robots[lIdx] = Robots{
			StartX:    pxVal,
			StartY:    pyVal,
			CurrPosX:  pxVal,
			CurrPosY:  pyVal,
			VelocityX: vxVal,
			VelocityY: vyVal,
		}
	}

	return robots, nil
}

func resetRobots(robots []Robots) {
	for rIdx := 0; rIdx < len(robots); rIdx++ {
		robots[rIdx].CurrPosX = robots[rIdx].StartX
		robots[rIdx].CurrPosY = robots[rIdx].StartY
	}
}

func displayArea(robots []Robots, width, height int) {

	line := make([]rune, width)

	for vIdx := 0; vIdx < height; vIdx++ {

		for hIdx := 0; hIdx < width; hIdx++ {
			line[hIdx] = ' '
		}
		for _, robot := range robots {
			if robot.CurrPosY == vIdx {
				line[robot.CurrPosX] = '*'
			}
		}

		fmt.Println(string(line))
	}
}

func createSnapshot(robots []Robots) []Robots {
	robotsCopy := make([]Robots, len(robots))
	copy(robotsCopy, robots)
	return robotsCopy
}
