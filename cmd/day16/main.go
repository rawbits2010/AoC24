package main

import (
	"fmt"
	"log"
	"math"

	"github.com/rawbits2010/AoC24/internal/inputhandler"
)

// TODO: Too slow to check every possible way through the maze
func main() {

	lines := inputhandler.ReadInput()
	maze := parseMaze(lines)
	/*
	   resultP1, path := maze.FindLowestScorePath()
	   _ = path
	*/
	resultP1 := maze.FindLowestScorePath()

	fmt.Printf("Result: %d", resultP1)
}

const EmptySpace = '.'
const Wall = '#'
const Start = 'S'
const End = 'E'

type CardinalDirection rune

const (
	EastFacing  CardinalDirection = '<'
	NorthFacing CardinalDirection = '^'
	WestFacing  CardinalDirection = '>'
	SouthFacing CardinalDirection = 'v'
)

type Position struct {
	X, Y int
}

type Raindeer struct {
	CurrPos Position
	Facing  CardinalDirection
}

type Tile struct {
	IsWalkable bool
	Pos        Position
	IsSelected bool
}

type Way struct {
	Facing    CardinalDirection
	IsOpen    bool
	IsChecked bool
	Pos       Position
}

type Junktion struct {
	PathIdx  int
	OpenWays []Way
}

// returns way and is depleted
func (j *Junktion) GetNextWay() (Way, bool) {
	for wIdx, way := range j.OpenWays {
		if !way.IsChecked {
			j.OpenWays[wIdx].IsChecked = true
			return way, false
		}
	}
	return Way{}, true
}

type Maze struct {
	Tiles         [][]Tile
	Raindeer      Raindeer
	Width, Height int
	StartPos      Position
	EndPos        Position
}

/*
func (m *Maze) FindLowestScorePath() (int, []Raindeer) {
		possiblePaths := m.findPathToEnd()

		if len(possiblePaths) == 0 {
			log.Fatal("didn't find a path")
		}

		pathScores := make([]int, len(possiblePaths))
		for pIdx, currPath := range possiblePaths {
			pathScores[pIdx] = calculateScore(currPath)
		}

		minScore := math.MaxInt
		minScoreIdx := -1
		for pIdx, score := range pathScores {
			if score < minScore {
				minScore = score
				minScoreIdx = pIdx
			}
		}

		m.visualizePath(possiblePaths[minScoreIdx], []Junktion{})

		return minScore, possiblePaths[minScoreIdx]
}
*/

func (m *Maze) FindLowestScorePath() int {
	return m.findPathToEnd()
}

func calculateScore(path []Raindeer) int {

	var pathScore int
	prevFacing := path[0].Facing
	for sIdx, step := range path {
		if sIdx == 0 {
			continue
		}

		if step.Facing != prevFacing {
			pathScore += 1000
		}

		pathScore += 1

		prevFacing = step.Facing
	}

	return pathScore
}

// func (m *Maze) findPathToEnd() [][]Raindeer {
func (m *Maze) findPathToEnd() int {

	//possiblePaths := make([][]Raindeer, 0, 20)
	minPathScore := math.MaxInt

	junktions := make([]Junktion, 0, 100)
	currPath := make([]Raindeer, 0, 100)

	m.Tiles[m.Raindeer.CurrPos.Y][m.Raindeer.CurrPos.X].IsSelected = true
	currPath = append(currPath, m.duplicateRaindeerState())

	for {

		if m.Raindeer.CurrPos.X == m.EndPos.X && m.Raindeer.CurrPos.Y == m.EndPos.Y {
			/*
				tempPath := make([]Raindeer, len(currPath))
				copy(tempPath, currPath)
				possiblePaths = append(possiblePaths, tempPath)
			*/
			currScore := calculateScore(currPath)
			if currScore < minPathScore {
				minPathScore = currScore
			}

		} else {

			switch m.Raindeer.Facing {

			case EastFacing:
				var isDeadEnd bool
				currPath, junktions, isDeadEnd = m.pickAWay(SouthFacing, EastFacing, NorthFacing, currPath, junktions)
				if !isDeadEnd {
					continue
				}

			case SouthFacing:
				var isDeadEnd bool
				currPath, junktions, isDeadEnd = m.pickAWay(SouthFacing, EastFacing, WestFacing, currPath, junktions)
				if !isDeadEnd {
					continue
				}

			case WestFacing:
				var isDeadEnd bool
				currPath, junktions, isDeadEnd = m.pickAWay(SouthFacing, NorthFacing, WestFacing, currPath, junktions)
				if !isDeadEnd {
					continue
				}

			case NorthFacing:
				var isDeadEnd bool
				currPath, junktions, isDeadEnd = m.pickAWay(WestFacing, EastFacing, NorthFacing, currPath, junktions)
				if !isDeadEnd {
					continue
				}

			}

		}

		// encountered a dead end
		if len(junktions) > 0 && len(currPath)-1 == junktions[len(junktions)-1].PathIdx {
			junktions = junktions[:len(junktions)-1]
		}

		if len(junktions) == 0 { // couldn't find a path
			break
		}

		lastJunktion := junktions[len(junktions)-1]
		for pIdx := lastJunktion.PathIdx; pIdx < len(currPath); pIdx++ {
			currStep := currPath[pIdx]
			m.Tiles[currStep.CurrPos.Y][currStep.CurrPos.X].IsSelected = false
		}
		currPath = currPath[:lastJunktion.PathIdx+1]

		m.Raindeer.CurrPos = currPath[lastJunktion.PathIdx].CurrPos
		m.Raindeer.Facing = currPath[lastJunktion.PathIdx].Facing
		continue
	}

	// return possiblePaths
	return minPathScore
}

func (m *Maze) pickAWay(direction1, direction2, direction3 CardinalDirection, currPath []Raindeer, junktions []Junktion) ([]Raindeer, []Junktion, bool) {

	if len(junktions) == 0 || len(currPath)-1 != junktions[len(junktions)-1].PathIdx {

		getWayForDirection := func(direction CardinalDirection) Way {

			switch direction {
			case SouthFacing:
				return m.isSouthWalkable()
			case NorthFacing:
				return m.isNorthWalkable()
			case EastFacing:
				return m.isEastWalkable()
			case WestFacing:
				return m.isWestWalkable()
			default:
				log.Fatalf("unknown direction found while picking a way '%v'", direction)
			}

			return Way{}
		}

		way1 := getWayForDirection(direction1)
		way2 := getWayForDirection(direction2)
		way3 := getWayForDirection(direction3)

		if IsJunktion(way1, way2, way3) {
			junktions = AddJunktion(len(currPath)-1, way1, way2, way3, junktions)
		} else {

			if way1.IsOpen {
				currPath = m.takeWay(way1, currPath)
				return currPath, junktions, false
			}
			if way2.IsOpen {
				currPath = m.takeWay(way2, currPath)
				return currPath, junktions, false
			}
			if way3.IsOpen {
				currPath = m.takeWay(way3, currPath)
				return currPath, junktions, false
			}

			return currPath, junktions, true
		}
	}

	if way, isDepleted := junktions[len(junktions)-1].GetNextWay(); !isDepleted {
		currPath = m.takeWay(way, currPath)
		return currPath, junktions, false
	}

	return currPath, junktions, true
}

func (m *Maze) takeWay(way Way, currPath []Raindeer) []Raindeer {

	m.Raindeer.CurrPos = way.Pos
	m.Raindeer.Facing = way.Facing

	m.Tiles[way.Pos.Y][way.Pos.X].IsSelected = true
	return append(currPath, m.duplicateRaindeerState())
}

func (m *Maze) duplicateRaindeerState() Raindeer {
	return Raindeer{
		CurrPos: Position{X: m.Raindeer.CurrPos.X, Y: m.Raindeer.CurrPos.Y},
		Facing:  m.Raindeer.Facing,
	}
}

func (m *Maze) isWalkable(facing CardinalDirection, targetPos Position) Way {

	if targetPos.X >= m.Width {
		return Way{
			Facing: facing,
			IsOpen: false,
			Pos:    targetPos,
		}
	}

	if !m.Tiles[targetPos.Y][targetPos.X].IsWalkable ||
		m.Tiles[targetPos.Y][targetPos.X].IsSelected {
		return Way{
			Facing: facing,
			IsOpen: false,
			Pos:    targetPos,
		}
	}

	return Way{
		Facing: facing,
		IsOpen: true,
		Pos:    targetPos,
	}
}

func (m *Maze) isEastWalkable() Way {

	targetPos := Position{X: m.Raindeer.CurrPos.X - 1, Y: m.Raindeer.CurrPos.Y}
	return m.isWalkable(EastFacing, targetPos)
}

func (m *Maze) isWestWalkable() Way {

	targetPos := Position{X: m.Raindeer.CurrPos.X + 1, Y: m.Raindeer.CurrPos.Y}
	return m.isWalkable(WestFacing, targetPos)
}

func (m *Maze) isNorthWalkable() Way {

	targetPos := Position{X: m.Raindeer.CurrPos.X, Y: m.Raindeer.CurrPos.Y - 1}
	return m.isWalkable(NorthFacing, targetPos)
}

func (m *Maze) isSouthWalkable() Way {

	targetPos := Position{X: m.Raindeer.CurrPos.X, Y: m.Raindeer.CurrPos.Y + 1}
	return m.isWalkable(SouthFacing, targetPos)
}

func (m *Maze) visualizePath(path []Raindeer, junktions []Junktion) {

	maze := make([][]rune, m.Height)
	for rIdx := 0; rIdx < m.Height; rIdx++ {
		maze[rIdx] = make([]rune, m.Width)
	}

	for vIdx, row := range m.Tiles {
		for hIdx, tile := range row {
			if tile.IsWalkable {
				maze[vIdx][hIdx] = EmptySpace
			} else {
				maze[vIdx][hIdx] = Wall
			}
		}
	}

	for _, junktion := range junktions {
		for _, openWay := range junktion.OpenWays {
			maze[openWay.Pos.Y][openWay.Pos.X] = '+'
		}
	}

	for _, step := range path {
		maze[step.CurrPos.Y][step.CurrPos.X] = rune(step.Facing)
	}

	for _, row := range maze {
		fmt.Println(string(row))
	}
}

func IsJunktion(way1, way2, way3 Way) bool {
	return way1.IsOpen && way2.IsOpen ||
		way1.IsOpen && way3.IsOpen ||
		way2.IsOpen && way3.IsOpen
}

func AddJunktion(currPathIdx int, way1, way2, way3 Way, junktions []Junktion) []Junktion {

	tempJunktion := Junktion{
		PathIdx:  currPathIdx,
		OpenWays: make([]Way, 0, 3),
	}

	if way1.IsOpen {
		tempJunktion.OpenWays = append(tempJunktion.OpenWays, way1)
	}

	if way2.IsOpen {
		tempJunktion.OpenWays = append(tempJunktion.OpenWays, way2)
	}

	if way3.IsOpen {
		tempJunktion.OpenWays = append(tempJunktion.OpenWays, way3)
	}

	return append(junktions, tempJunktion)
}

func parseMaze(lines []string) Maze {

	tempMaze := Maze{
		Height: len(lines),
		Width:  len(lines[0]),
	}

	tiles := make([][]Tile, tempMaze.Height)
	for vIdx := 0; vIdx < tempMaze.Height; vIdx++ {
		tiles[vIdx] = make([]Tile, tempMaze.Width)
	}

	for vIdx, row := range lines {
		for hIdx, spot := range row {

			tiles[vIdx][hIdx].Pos = Position{X: hIdx, Y: vIdx}

			switch spot {

			case Wall:
				tiles[vIdx][hIdx].IsWalkable = false

			case Start:
				tempMaze.StartPos = Position{X: hIdx, Y: vIdx}
				tiles[vIdx][hIdx].IsWalkable = true

			case End:
				tempMaze.EndPos = Position{X: hIdx, Y: vIdx}
				tiles[vIdx][hIdx].IsWalkable = true

			case EmptySpace:
				tiles[vIdx][hIdx].IsWalkable = true

			default:
				log.Fatalf("unknown tile symbol in maze while parsing '%c'", spot)
			}
		}
	}
	tempMaze.Tiles = tiles

	tempMaze.Raindeer = Raindeer{
		CurrPos: Position{X: tempMaze.StartPos.X, Y: tempMaze.StartPos.Y},
		Facing:  EastFacing,
	}

	return tempMaze
}
