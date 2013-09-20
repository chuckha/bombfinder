package game

import (
	"fmt"
)

// Actual data enums
const (
	zero = iota
	one
	two
	three
	four
	five
	six
	seven
	eight
	bomb
)

// display values
const (
	None    = " "
	Flag    = "⚑"
	Unknown = "?"
	Bomb    = "B"
	Zero    = "-"
	One     = "1" //"①"
	Two     = "2" //"②"
	Three   = "3" //"③"
	Four    = "4" //"④"
	Five    = "5" //"⑤"
	Six     = "6" //"⑥"
	Seven   = "7" //"⑦"
	Eight   = "8" //"⑧"
)

// Display data enum
const (
	none    = iota // Display nothing
	actual         // Display actual value
	flag           // Display flag
	unknown        // Display question mark
)

type Cell struct {
	Display  int `json:"-"`
	Pressed  bool
	Val      string
	Player   *Player
	Data     int `json:"-"`
	Row, Col int `json:"-"`
}

func (c *Cell) RightClick(player *Player) {
	// If it's pressed already, don't do anything
	if c.Pressed {
		return
	}
	// If it's a flag, only let the owner player change it
	if c.Display == flag && player != c.Player {
		return
	}
	switch c.Display {
	case none:
		c.Display = flag
		c.Val = Flag
		c.Player = player
	case flag:
		c.Display = none
		c.Val = None
	}
}

func (c *Cell) LeftClick() {
	if c.Pressed {
		return
	}
	c.Display = actual
	c.Pressed = true
	c.Player = &Player{} // Unset the player if there was one
	c.Val = StringData(c.Data)
}

func (c *Cell) String() string {
	switch c.Display {
	case none:
		return None
	case flag:
		return Flag
	case actual:
		return StringData(c.Data)
	default:
		return ""
	}
}

func StringData(data int) string {
	switch data {
	case bomb:
		return Bomb
	case zero:
		return Zero
	case one:
		return One
	case two:
		return Two
	case three:
		return Three
	case four:
		return Four
	case five:
		return Five
	case six:
		return Six
	case seven:
		return Seven
	case eight:
		return Eight
	default:
		return fmt.Sprintf("%d", data)
	}
}
