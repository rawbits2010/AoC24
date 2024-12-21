package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func solvePart2(robots []Robots, fieldWidth, fieldHeight int) []Constellation {

	outlierTest := NewZScore(7)
	minGroupCount := 10
	constellations := make([]Constellation, 0)

	startChecksum := generateChecksum(robots)
	var tIdx int
	for {
		tIdx++

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

		if startChecksum == generateChecksum(robots) {
			break
		}

		groupsCounts := make([]int, 0)
		for rIdx := 0; rIdx < len(robots); rIdx++ {
			count := countConnected(rIdx, robots)
			groupsCounts = append(groupsCounts, count)
		}
		sort.Ints(groupsCounts)
		largestGroupCount := groupsCounts[len(groupsCounts)-1]

		if outlierTest.Add(float64(largestGroupCount)) && largestGroupCount > minGroupCount {

			robotsCopy := createSnapshot(robots)

			tempConstellation := Constellation{
				RobotCount: largestGroupCount,
				AtTime:     tIdx,
				Snapshot:   robotsCopy,
			}
			constellations = append(constellations, tempConstellation)
		}
	}

	return constellations
}

type Constellation struct {
	RobotCount int
	AtTime     int
	Snapshot   []Robots
}

func countConnected(toCheckIdx int, robots []Robots) int {

	checked := make([]Robots, 0, 10)
	checked = append(checked, robots[toCheckIdx])
	robots[toCheckIdx].IsChecked = true

	for {
		lastIdx := len(checked) - 1

		for checkIdx := lastIdx; checkIdx < len(checked); checkIdx++ {

			checkedRobot := checked[checkIdx]

			for rIdx := 0; rIdx < len(robots); rIdx++ {
				if robots[rIdx].IsChecked {
					continue
				}

				// left
				if checkedRobot.CurrPosX-1 == robots[rIdx].CurrPosX {
					if checkedRobot.CurrPosY == robots[rIdx].CurrPosY {
						robots[rIdx].IsChecked = true
						checked = append(checked, robots[rIdx])
					}
				}

				// right
				if checkedRobot.CurrPosX+1 == robots[rIdx].CurrPosX {
					if checkedRobot.CurrPosY == robots[rIdx].CurrPosY {
						robots[rIdx].IsChecked = true
						checked = append(checked, robots[rIdx])
					}
				}

				// up
				if checkedRobot.CurrPosX == robots[rIdx].CurrPosX {
					if checkedRobot.CurrPosY-1 == robots[rIdx].CurrPosY {
						robots[rIdx].IsChecked = true
						checked = append(checked, robots[rIdx])
					}
				}

				// down
				if checkedRobot.CurrPosX == robots[rIdx].CurrPosX {
					if checkedRobot.CurrPosY+1 == robots[rIdx].CurrPosY {
						robots[rIdx].IsChecked = true
						checked = append(checked, robots[rIdx])
					}
				}
			}
		}

		if lastIdx == len(checked)-1 {
			break
		}
	}

	resetCheckedStatus(robots)

	return len(checked)
}

func resetCheckedStatus(robots []Robots) {
	for idx := 0; idx < len(robots); idx++ {
		robots[idx].IsChecked = false
	}
}

func generateChecksum(robots []Robots) int {

	var checksum int
	for rIdx, robot := range robots {
		checksum += (rIdx * 1000000) + (robot.CurrPosX * 1000) + robot.CurrPosY
	}

	return checksum
}

func ensureOutputDir(dirname string) error {

	if _, err := os.Stat(dirname); os.IsNotExist(err) {

		err := os.Mkdir(dirname, 0755)
		if err != nil {
			return fmt.Errorf("couldn't create directory '%s': %w", dirname, err)
		}
	} else {

		files, err := os.ReadDir(dirname)
		if err != nil {
			return fmt.Errorf("couldn't read directory contents '%s': %w", dirname, err)
		}

		for _, file := range files {
			fullpath := filepath.Join(dirname, file.Name())
			err = os.RemoveAll(fullpath)
			if err != nil {
				return fmt.Errorf("couldn't remove file '%s': %w", fullpath, err)
			}
		}
	}

	return nil
}

func saveSnapshotToFile(filename string, robots []Robots, width, height int) error {

	if _, err := os.Stat(filename); os.IsExist(err) {
		return fmt.Errorf("error saving snapshot file '%s': already exists", filename)
	}

	outFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error saving snapshot to file '%s': %w", filename, err)
	}
	defer outFile.Close()

	line := make([]rune, width+1)
	line[len(line)-1] = '\n'

	for vIdx := 0; vIdx < height; vIdx++ {

		for hIdx := 0; hIdx < width; hIdx++ {
			line[hIdx] = ' '
		}
		for _, robot := range robots {
			if robot.CurrPosY == vIdx {
				line[robot.CurrPosX] = '*'
			}
		}

		outFile.WriteString(string(line))
	}

	return nil
}
