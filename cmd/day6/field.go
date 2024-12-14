package main

import "log"

type Field struct {
	Height        int
	Width         int
	Tiles         [][]Tile
	Elf           Guard
	TempObstacleX int
	TempObstacleY int
}

func NewField(field []string) *Field {
	var temp Field

	temp.TempObstacleX = -1
	temp.TempObstacleY = -1

	temp.Height = len(field)
	temp.Width = len(field[0])

	temp.Tiles = make([][]Tile, temp.Height)
	for idx := 0; idx < temp.Height; idx++ {
		temp.Tiles[idx] = make([]Tile, temp.Width)
	}

	for vIdx, line := range field {
		for hIdx, char := range line {

			switch char {
			case empty:
				temp.Tiles[vIdx][hIdx] = *NewTile(true)

			case obstacle:
				temp.Tiles[vIdx][hIdx] = *NewTile(false)

			case guard:
				temp.Elf = *NewGuard(hIdx, vIdx, UpDir)
				temp.Tiles[vIdx][hIdx] = *NewTile(true)
				temp.Tiles[vIdx][hIdx].Touch(temp.Elf.StepCount, temp.Elf.TurnCount, temp.Elf.Facing)

			default:
				log.Fatalf("unknown object on the field '%s'", string(char))
			}
		}
	}

	return &temp
}

// returns if step/turn happened and if in a loop
func (f *Field) Step() (bool, bool) {

	if f.CheckIfStepExits() {
		return false, false
	}

	targetX, targetY := f.Elf.GetTargetTileCoords()
	targetTile := &f.Tiles[targetY][targetX]

	if targetTile.IsTouchedWhileFacing(f.Elf.Facing) {
		return true, true
	}

	if targetTile.IsEmpty {

		f.Elf.Step()
		targetTile.Touch(f.Elf.StepCount, f.Elf.TurnCount, f.Elf.Facing)
	} else {

		f.Elf.Turn()
	}

	return true, false
}

func (f *Field) CheckIfStepExits() bool {

	switch f.Elf.Facing {
	case UpDir:
		return f.Elf.Y == 0
	case RightDir:
		return f.Elf.X == f.Width-1
	case DownDir:
		return f.Elf.Y == f.Height-1
	case LeftDir:
		return f.Elf.X == 0
	default:
		log.Fatalf("invalid guard facing '%s'", f.Elf.Facing)
	}

	return false
}

func (f *Field) PlaceTempObstacle(x, y int) {

	if f.TempObstacleX != -1 {
		f.Tiles[f.TempObstacleY][f.TempObstacleX].UnBlock()
	}

	if f.Tiles[y][x].IsEmpty {

		f.TempObstacleX = x
		f.TempObstacleY = y
		f.Tiles[y][x].Block()
	} else {

		f.TempObstacleX = -1
		f.TempObstacleY = -1
	}
}

func (f *Field) Reset() {

	if f.TempObstacleX != -1 {
		f.Tiles[f.TempObstacleY][f.TempObstacleX].UnBlock()
	}

	for vIdx := 0; vIdx < f.Height; vIdx++ {
		for hIdx := 0; hIdx < f.Width; hIdx++ {
			f.Tiles[vIdx][hIdx].ResetTouch()
		}
	}

	f.Elf.Reset()
}
