package game

import (
	"fmt"
	"math/rand"
	"strings"
)

type Board struct {
	Height, Width, Bombs int
	Win                  bool
	Finished             bool
	Field                [][]*Cell
}

func NewBoard(h, w, bombs int) (*Board, error) {
	field, err := BuildField(h, w, bombs)
	if err != nil {
		return nil, err
	}
	return &Board{
		Height: h,
		Width:  w,
		Bombs:  bombs,
		Field:  field,
	}, nil
}

func (b *Board) String() string {
	out := ""
	for i := 0; i < b.Width*2+1; i++ {
		out += "-"
	}
	out += "\n"
	for i := range b.Field {
		out += "|"
		rowRep := make([]string, len(b.Field[i]))
		for j := range b.Field[i] {
			rowRep[j] = b.Field[i][j].String()
		}
		out += strings.Join(rowRep, "|")
		out += "|\n"
	}
	for i := 0; i < b.Width*2+1; i++ {
		out += "-"
	}
	return out
}

func (b *Board) RightClick(row, col int) {
	b.Field[row][col].RightClick()
}
func (b *Board) LeftClick(row, col int) {
	b.Field[row][col].LeftClick()
	if b.Field[row][col].Data == zero {
		b.ZeroAround(row, col)
	}
}

func (b *Board) ZeroAround(row, col int) {
	toZeroAround := make(chan *Cell, 8)
	leftClicker := func(c *Cell) {
		if c.Display == none {
			c.LeftClick()
			if c.Data == zero {
				toZeroAround <- c
			}
		}
	}
	EachNeighbor(b.Field, row, col, leftClicker)
	close(toZeroAround)
	for cell := range toZeroAround {
		b.ZeroAround(cell.Row, cell.Col)
	}
}

func BuildField(height, width, bombs int) ([][]*Cell, error) {
	if height*width <= bombs {
		return nil, fmt.Errorf("Impossible board, too many bombs.")
	}
	// Alloc
	field := make([][]*Cell, height)
	for i := 0; i < height; i++ {
		field[i] = make([]*Cell, width)
		for j := range field[i] {
			field[i][j] = &Cell{
				Row: i,
				Col: j,
			}
		}
	}

	// Place the bombs
	bombsPlaced := 0
	for bombsPlaced < bombs {
		row := rand.Intn(height)
		col := rand.Intn(width)
		if field[row][col].Data == bomb {
			continue
		}
		field[row][col].Data = bomb
		bombsPlaced += 1
	}

	for i := range field {
		for j := range field[i] {
			// If we don't have a bomb we need to calculate around it
			if field[i][j].Data != bomb {
				field[i][j].Data = CountBombsAround(field, i, j)
			}
		}
	}
	// Calculate the cells
	return field, nil
}

const (
	upleft = iota
	up
	upright
	right
	downright
	down
	downleft
	left
)

func CountBombsAround(field [][]*Cell, row, col int) int {
	bombs := make(chan int, 8)
	// The func that will run on each neighbor cell
	counter := func(c *Cell) {
		if c.Data == bomb {
			bombs <- 1
		}
	}
	EachNeighbor(field, row, col, counter)
	close(bombs)
	count := 0
	for v := range bombs {
		count += v
	}
	return count
}

// Given a row and col, find all valid cells and run a function
// on each valid cell.
func EachNeighbor(field [][]*Cell, row, col int, f func(*Cell)) {
	height := len(field) - 1
	width := len(field[0]) - 1

	//order of Countme: upleft, up, upright, right, downright, down, downleft, left
	directions := []bool{true, true, true, true, true, true, true, true}
	// Each of these ifs filters out the ones it should not count
	// At the end we will be left with only valid countable directions
	if row == 0 {
		directions[upleft] = false
		directions[up] = false
		directions[upright] = false
	}
	if col == 0 {
		directions[upleft] = false
		directions[downleft] = false
		directions[left] = false
	}
	if row == height {
		directions[downright] = false
		directions[down] = false
		directions[downleft] = false
	}
	if col == width {
		directions[upright] = false
		directions[right] = false
		directions[downright] = false
	}

	// We now only count the directions we know are valid.
	for i, shouldRun := range directions {
		if !shouldRun {
			continue
		}
		switch i {
		case upleft:
			f(field[row-1][col-1])
		case up:
			f(field[row-1][col])
		case upright:
			f(field[row-1][col+1])
		case right:
			f(field[row][col+1])
		case downright:
			f(field[row+1][col+1])
		case down:
			f(field[row+1][col])
		case downleft:
			f(field[row+1][col-1])
		case left:
			f(field[row][col-1])
		}
	}
}
