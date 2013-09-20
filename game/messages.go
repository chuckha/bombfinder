package game

import ()

type InMessage struct {
	Type  string
	Value MoveMessage
}

type MoveMessage struct {
	Click string
	Row   int
	Col   int
}

type OutMessage struct {
	Type  string
	Value interface{}
}

type InfoMessage struct {
	Required int
	Have     int
}

type PlayerMessage struct {
	Players []*Player
}

type ErrorMessage struct {
	Type  string
	Value string
}
