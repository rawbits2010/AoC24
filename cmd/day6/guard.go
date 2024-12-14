package main

import "log"

type Direction string

const (
	UpDir    = "Up"
	RightDir = "Right"
	DownDir  = "Down"
	LeftDir  = "Left"
)

type Guard struct {
	X, Y           int
	Facing         Direction
	StartX, StartY int
	StartFacing    Direction
	StepCount      int
	TurnCount      int
}

func NewGuard(x, y int, facing Direction) *Guard {
	return &Guard{
		X:           x,
		Y:           y,
		Facing:      facing,
		StartX:      x,
		StartY:      y,
		StartFacing: facing,
		StepCount:   0,
		TurnCount:   0,
	}
}

// returns X, Y for targeted tile
func (g *Guard) GetTargetTileCoords() (int, int) {

	switch g.Facing {
	case UpDir:
		return g.X, g.Y - 1
	case RightDir:
		return g.X + 1, g.Y
	case DownDir:
		return g.X, g.Y + 1
	case LeftDir:
		return g.X - 1, g.Y
	default:
		log.Fatalf("invalid guard facing '%s'", g.Facing)
	}

	return -1, -1
}

func (g *Guard) Step() {

	switch g.Facing {
	case UpDir:
		g.Y--
	case RightDir:
		g.X++
	case DownDir:
		g.Y++
	case LeftDir:
		g.X--
	}

	g.StepCount++
}

func (g *Guard) Turn() {

	switch g.Facing {
	case UpDir:
		g.Facing = RightDir
	case RightDir:
		g.Facing = DownDir
	case DownDir:
		g.Facing = LeftDir
	case LeftDir:
		g.Facing = UpDir
	}

	g.TurnCount++
}

func (g *Guard) Reset() {
	g.X = g.StartX
	g.Y = g.StartY
	g.Facing = g.StartFacing
	g.StepCount = 0
	g.TurnCount = 0
}
