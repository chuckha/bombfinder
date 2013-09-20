package game

import (
	"code.google.com/p/go.net/websocket"
)

const (
	red   = "#ee0a0a"
	green = "#0aee0a"
	blue  = "#0a0aee"
)

var players = 0
var colors = []string{red, green, blue}

type Player struct {
	WS      *websocket.Conn `json:"-"`
	Name    string
	Playing bool
	Color   string
}

func NewPlayer(ws *websocket.Conn, name string) *Player {
	return &Player{
		WS:      ws,
		Playing: true,
		Name:    name,
		Color:   colors[players%len(colors)],
	}
}

func (p *Player) Die() {
	p.Playing = false
}

func (p *Player) Send(msg *OutMessage) {
	websocket.JSON.Send(p.WS, msg)
}
