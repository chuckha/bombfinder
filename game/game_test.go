package game

import (
	"testing"
)

func TestRemovePlayer(t *testing.T) {
	g := NewGame(4, 10, 10, 10)
	p1 := &Player{}
	p2 := &Player{}
	p3 := &Player{}
	p4 := &Player{}

	g.AddPlayer(p1)
	g.AddPlayer(p2)
	g.AddPlayer(p3)
	g.AddPlayer(p4)
	g.RemovePlayer(p2)
	if len(g.Players) != 3 {
		t.Errorf("Should only have 3 players, found: %d", len(g.Players))
	}
	if g.Players[2] != p4 {
		t.Errorf("Player4 should be in the second position")
	}
}
func TestRemovePlayerEdge(t *testing.T) {
	g := NewGame(2, 10, 10, 10)
	p1 := &Player{}
	g.AddPlayer(p1)
	g.RemovePlayer(p1)
	if len(g.Players) != 0 {
		t.Errorf("Expecting 0 players, got: %d", len(g.Players))
	}
}
