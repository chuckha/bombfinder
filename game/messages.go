package game

import ()

type InMessage struct {
	Type  string
	Value Message
}

type Message struct {
	Click string
	Row   int
	Col   int
}

type OutMessage struct {
	Type  string
	Value interface{}
}
