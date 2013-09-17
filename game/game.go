package game

import ()

type Game struct {
	Board *Board
}

func NewGame() *Game {
	return &Game{
		Board: NewBoard(10, 10, 10),
	}
}
