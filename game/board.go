package game

import (
	"fmt"
	"math/rand"
	"strings"
)

type Board struct {
	Height, Width, Bombs int
	Finished             bool
	Won                  bool
	Numbers              int
	Clicked              int
	Field                [][]*Cell
}

func NewBoardWithField(field [][]*Cell) *Board {
	height, width, bombCount, numbersCount := FieldInfo(field)
	return &Board{
		Height:  height,
		Width:   width,
		Bombs:   bombCount,
		Field:   field,
		Numbers: numbersCount,
	}
}

func NewBoard(h, w, bombs int) (*Board, error) {
	field, err := BuildField(h, w, bombs)
	if err != nil {
		return nil, err
	}
	_, _, _, numbers := FieldInfo(field)
	return &Board{
		Height:  h,
		Width:   w,
		Bombs:   bombs,
		Field:   field,
		Numbers: numbers,
	}, nil
}

// Return
// Height, Width, BombCount, NumbersCount
func FieldInfo(field [][]*Cell) (int, int, int, int) {
	height := len(field)
	width := len(field[0])
	bombCount := 0
	numbersCount := 0
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			switch field[i][j].Data {
			case bomb:
				bombCount += 1
			case zero:
				continue
			default:
				numbersCount += 1
			}
		}
	}
	return height, width, bombCount, numbersCount
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
			rowRep[j] = fmt.Sprintf("%d", b.Field[i][j].Data)
		}
		out += strings.Join(rowRep, "|")
		out += "|\n"
	}
	for i := 0; i < b.Width*2+1; i++ {
		out += "-"
	}
	return out
}

func (b *Board) RightClick(row, col int) error {
	if b.Finished {
		return nil
	}
	b.Field[row][col].RightClick()
	return nil
}
func (b *Board) LeftClick(row, col int) error {
	// If the board is finished, don't do anything
	if b.Finished {
		return nil
	}
	// If it's already been clicked, don't do anything
	if b.Field[row][col].Pressed {
		return nil
	}

	b.Field[row][col].LeftClick()

	// What did we click?
	switch b.Field[row][col].Data {
	case zero:
		// If we get a zero, we need to reveal all zeroes around
		b.ZeroAround(row, col)
		return nil
	case bomb:
		// If we click a bomb, we lose the game
		b.Finished = true
		b.Won = false
		return fmt.Errorf("Boom!")
	default:
		// Otherwise we clicked a number. One step closer to victory
		b.Clicked += 1
		if b.Clicked == b.Numbers {
			b.Finished = true
			b.Won = true
		}
		return nil
	}
}

func (b *Board) ZeroAround(row, col int) {
	toZeroAround := make(chan *Cell, 8)
	leftClicker := func(c *Cell) {
		if c.Display == none {
			c.LeftClick()
			switch c.Data {
			case zero:
				toZeroAround <- c
			case one, two, three, four, five, six, seven, eight:
				b.Clicked += 1
				if b.Clicked == b.Numbers {
					b.Finished = true
					b.Won = true
				}
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
