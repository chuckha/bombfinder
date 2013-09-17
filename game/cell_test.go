package game

import (
	"testing"
)

func TestCell(t *testing.T) {
	c := &Cell{}
	c.Data = one
	if c == nil {
		t.Errorf("Cell should not be nil")
	}
}

func TestDisplay(t *testing.T) {
	c := &Cell{}
	c.Data = bomb
	if c.String() != " " {
		t.Errorf("An unrevealed cell should be blank not %s", c.String())
	}
	c.Display = flag
	if c.String() != Flag {
		t.Errorf("A cell that has been set to flag should display a flag not %s", c.String())
	}
	c.Display = unknown
	if c.String() != Unknown {
		t.Errorf("A cell that has been set to unknown should display a question mark not %s", c.String())
	}
	c.Display = actual
	if c.String() != Bomb {
		t.Errorf("A cell that has been set to actual should display its actual value")
	}
}

func TestRightClick(t *testing.T) {
	c := &Cell{}
	c.Data = five
	if c.Display != none {
		t.Errorf("A new cell should display none")
	}
	c.RightClick()
	if c.Display != flag || c.String() != Flag {
		t.Errorf("Display should be flag after a right click on an empty cell takes place")
	}
	c.RightClick()
	if c.Display != unknown {
		t.Errorf("A flag that has been right clicked should be a question mark")
	}
	c.RightClick()
	if c.Display != none {
		t.Errorf("An unknown cell right clicked should now be none")
	}
}

func TestLeftClick(t *testing.T) {
	c := &Cell{}
	c.Data = four
	c.LeftClick()
	if c.Display != actual || c.String() != Four {
		t.Errorf("Cell should always display actual after a left click")
	}
	c = &Cell{}
	c.Data = bomb
	c.LeftClick()
	if c.Display != actual || c.String() != Bomb {
		t.Errorf("Cell should always display actual after a left click")
	}

}
