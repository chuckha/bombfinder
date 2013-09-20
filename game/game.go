package game

import (
	"fmt"
)

type Game struct {
	Board      *Board
	NumPlayers int
	Players    []*Player
}

func NewGame(numPlayers, h, w, bombs int) *Game {
	b, _ := NewBoard(h, w, bombs)
	return &Game{
		Board:      b,
		NumPlayers: numPlayers,
		Players:    make([]*Player, 0, numPlayers),
	}
}

func (g *Game) RemovePlayer(p *Player) {
	for i, player := range g.Players {
		if p == player {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			return
		}
	}
}

func (g *Game) ActivePlayers() int {
	return len(g.Players)
}

func (g *Game) AddPlayer(p *Player) {
	g.Players = append(g.Players, p)
}

func (g *Game) IsFull() bool {
	return len(g.Players) == g.NumPlayers
}

func (g *Game) SendAll(msg *OutMessage) {
	fmt.Printf("Sending: %v\n", msg)
	fmt.Printf("Sending that message to %d players\n", len(g.Players))
	for _, p := range g.Players {
		p.Send(msg)
	}
}

func (g *Game) SendBoard() {
	msg := &OutMessage{
		Type:  "board",
		Value: g.Board,
	}
	g.SendAll(msg)
}

func (g *Game) SendInfo() {
	msg := &OutMessage{
		Type: "info",
		Value: &InfoMessage{
			Required: g.NumPlayers,
			Have:     g.ActivePlayers(),
		},
	}
	g.SendAll(msg)
}

func (g *Game) SendPlayerInfo() {
	msg := &OutMessage{
		Type: "players",
		Value: &PlayerMessage{
			Players: g.Players,
		},
	}
	g.SendAll(msg)
}
