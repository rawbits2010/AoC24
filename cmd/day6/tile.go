package main

type Tile struct {
	IsEmpty     bool
	SteppedOn   []int
	OnTurn      []int
	WhileFacing []Direction
	StepCount   int
}

func NewTile(isEmpty bool) *Tile {
	return &Tile{
		IsEmpty:     isEmpty,
		SteppedOn:   make([]int, 0),
		OnTurn:      make([]int, 0),
		WhileFacing: make([]Direction, 0),
		StepCount:   0,
	}
}

func (t *Tile) Touch(currStep, currTurn int, facing Direction) {
	t.SteppedOn = append(t.SteppedOn, currStep)
	t.OnTurn = append(t.OnTurn, currTurn)
	t.WhileFacing = append(t.WhileFacing, facing)
	t.StepCount++
}

func (t *Tile) IsTouched() bool {
	return t.StepCount != 0
}

func (t *Tile) IsTouchedWhileFacing(facing Direction) bool {
	for _, oldFacing := range t.WhileFacing {
		if oldFacing == facing {
			return true
		}
	}
	return false
}

func (t *Tile) Block() {
	t.IsEmpty = false
}

func (t *Tile) UnBlock() {
	t.IsEmpty = true
}

func (t *Tile) ResetTouch() {
	t.SteppedOn = make([]int, 0)
	t.OnTurn = make([]int, 0)
	t.WhileFacing = make([]Direction, 0)
	t.StepCount = 0
}
