package game

import (
	"testing"
)

func buildBoard() *Board {
	placement := [][]int{
		[]int{2, 9, 9, 9, 2},
		[]int{9, 3, 4, 9, 3},
		[]int{3, 4, 4, 5, 9},
		[]int{9, 9, 9, 9, 9},
		[]int{9, 9, 8, 9, 4},
		[]int{9, 9, 9, 9, 2},
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
	return NewBoardWithField(board)
}

func CreateTestBoard() *Board {
	height := 4
	width := 5
	bombs := 10
	b, _ := NewBoard(height, width, bombs)
	return b
}

func TestZeroAroundWin(t *testing.T) {
	placement := [][]int{
		[]int{0, 0, 0},
		[]int{0, 1, 1},
		[]int{0, 1, 9},
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
	b := NewBoardWithField(board)
	b.LeftClick(0, 0)
	if !b.Won {
		t.Errorf("Zero around is able to win.: %d, %d", b.Numbers, b.Clicked)
	}
}

func TestClickedStartsAtZero(t *testing.T) {
	b := buildBoard()
	if b.Clicked != 0 {
		t.Errorf("We haven't clicked anything yet!")
	}
}
func TestMultipleClicks(t *testing.T) {
	b := buildBoard()
	b.LeftClick(0, 0)
	b.LeftClick(0, 0)
	b.LeftClick(0, 0)
	b.LeftClick(0, 0)
	b.LeftClick(0, 0)
	if b.Clicked != 1 {
		t.Errorf("Multiple Clicks should only count once.\nExpected: %d\nGot: %d", 1, b.Clicked)
	}
}

// If all the values are clicked, we have won
func TestWinCondition(t *testing.T) {
	b := buildBoard()
	b.LeftClick(0, 0)
	b.LeftClick(0, 4)
	b.LeftClick(1, 1)
	b.LeftClick(1, 2)
	b.LeftClick(1, 4)
	b.LeftClick(2, 0)
	b.LeftClick(2, 1)
	b.LeftClick(2, 2)
	b.LeftClick(2, 3)
	b.LeftClick(4, 2)
	b.LeftClick(4, 4)
	if b.Won {
		t.Errorf("We haven't won yet")
	}
	b.LeftClick(5, 4)
	if !b.Won {
		t.Errorf("We actually won.")
	}
	// Victory!!
	err := b.LeftClick(0, 1)
	if err != nil {
		t.Errorf("We should not do anything on a victory board")
	}
}

// If you click a finished board, nothing will happen
func TestClickFinishedBoard(t *testing.T) {
	b := buildBoard()
	b.LeftClick(0, 1)
	if b.Finished == false {
		t.Errorf("Game is actually finished")
	}
	err := b.LeftClick(1, 1)
	if err != nil {
		t.Errorf("A click on a finished board should not do anything")
	}
	err = b.RightClick(1, 1)
	if err != nil {
		t.Errorf("A right click on a finished board should not do anything")
	}
}

func TestLoseCondition(t *testing.T) {
	// If I leftclick on a bomb, I should lose
	b := buildBoard()
	err := b.LeftClick(0, 1)
	if err == nil {
		t.Errorf("You clicked on a bomb, you should have an error here")
	}
	if b.Finished == false || b.Won == true {
		t.Errorf("Should not have finished or should not have won.")
	}
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

func TestFieldInfo(t *testing.T) {
	b := buildBoard()
	height, width, bombs, numbers := FieldInfo(b.Field)
	if height != 6 {
		t.Errorf("Miscounted height")
	}
	if width != 5 {
		t.Errorf("Miscounted the width")
	}
	if bombs != 18 {
		t.Errorf("Miscounted the bombs")
	}
	if numbers != 12 {
		t.Errorf("Miscounted the numbers")
	}
}
