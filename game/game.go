package game

import ()

type Game struct {
	Board *Board
}

func NewGame() *Game {
	b, _ := NewBoard(10, 10, 10)
	return &Game{
		Board: b,
	}
}
