package game

import (
	"testing"
)

func CreateTestBoard() *Board {
	height := 4
	width := 5
	bombs := 10
	b, _ := NewBoard(height, width, bombs)
	return b
}

func TestBoardCreation(t *testing.T) {
	height := 4
	width := 5
	bombs := 10
	b, err := NewBoard(height, width, bombs)
	if err != nil {
		t.Errorf("Got an error creating a board: %s", err.Error())
	}
	if b.Height == 0 {
		t.Errorf("Height must be > 0")
	}
	if b.Width == 0 {
		t.Errorf("Width must be > 0")
	}
	if b.Height != height {
		t.Errorf("Board set the height wrong")
	}
	if b.Width != width {
		t.Errorf("Board set the width wrong")
	}
	if b.Bombs != bombs {
		t.Errorf("Board did not set bombs correctly")
	}
	if len(b.Field) != height {
		t.Errorf("Board did not set field height correctly")
	}
	if len(b.Field[0]) != width {
		t.Errorf("Board did not set width of row correctly")
	}
	if b.Field[3][4].Row != 3 {
		t.Errorf("Cell's row is wrong")
	}
	if b.Field[3][4].Col != 4 {
		t.Errorf("Cell's Col is wrong")
	}
}

func TestBuildField(t *testing.T) {
	field, err := BuildField(50, 20, 100)
	if err != nil {
		t.Errorf("BuildField returned an error with valid inputs: %s", err.Error())
	}
	bombCount := 0
	for _, row := range field {
		for _, col := range row {
			if col.Data == bomb {
				bombCount += 1
			}
		}
	}
	if bombCount != 100 {
		t.Errorf("Needed 100 bombs, got %d", bombCount)
	}
	field, err = BuildField(1, 1, 100)
	if err == nil {
		t.Errorf("An impossible board was created")
	}
}

func TestZeroAround(t *testing.T) {
	placement := [][]int{
		[]int{0, 0, 2, 9, 2},
		[]int{0, 0, 2, 9, 3},
		[]int{2, 3, 4, 5, 9},
		[]int{9, 9, 9, 9, 9},
		[]int{9, 9, 8, 9, 4},
		[]int{9, 9, 9, 9, 2},
	}
	board := make([][]*Cell, len(placement))
	for i := range placement {
		board[i] = make([]*Cell, len(placement[i]))
		for j := range placement[i] {
			board[i][j] = &Cell{
				Row:  i,
				Col:  j,
				Data: placement[i][j],
			}
		}
	}
	b := &Board{
		Height: 6,
		Width:  5,
		Bombs:  15,
		Field:  board,
	}
	b.Field[0][0].LeftClick()
	b.ZeroAround(0, 0)
	if b.Field[1][0].Display != actual {
		t.Errorf("Zero around isn't going downwards")
	}
	if b.Field[2][0].Display != actual {
		t.Errorf("Zero around isn't going downwards twice")
	}
	if b.Field[3][0].Display == actual {
		t.Errorf("Zero around is going too far downwards")
	}
	if b.Field[2][2].Display != actual {
		t.Errorf("Zero around isn't going out far enough")
	}
}

func TestCountAround(t *testing.T) {
	placement := [][]int{
		[]int{0, 9, 9, 9, 0},
		[]int{9, 0, 0, 9, 0},
		[]int{0, 0, 0, 0, 9},
		[]int{9, 9, 9, 9, 9},
		[]int{9, 9, 0, 9, 0},
		[]int{9, 9, 9, 9, 0},
	}
	board := make([][]*Cell, len(placement))
	for i := range placement {
		board[i] = make([]*Cell, len(placement[i]))
		for j := range placement[i] {
			board[i][j] = &Cell{
				Data: placement[i][j],
			}
		}
	}
	tc := []struct {
		row, col, expected int
	}{
		{0, 0, 2},
		{0, 4, 2},
		{5, 0, 3},
		{5, 4, 2},
		{4, 2, 8},
	}
	for _, testcase := range tc {
		got := CountBombsAround(board, testcase.row, testcase.col)
		if got != testcase.expected {
			t.Errorf("\nExpected: %d\nGot: %d\n", testcase.expected, got)
		}
	}

}
