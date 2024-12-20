package ui

import (
	"fmt"

	"bretbelgarde.com/adventure/actors"
	"bretbelgarde.com/adventure/items"
	"bretbelgarde.com/adventure/maps"
	ut "bretbelgarde.com/adventure/utils"
	tc "github.com/gdamore/tcell/v2"
)

type Cursor struct {
	X        int
	Y        int
	floor    int
	Rune     rune
	Color    tc.Style
	IsActive bool
	Current  *maps.Map
}

func NewCursor(x, y, f int, m *maps.Map) *Cursor {
	return &Cursor{
		X:        x,
		Y:        y,
		floor:    f,
		Rune:     tc.RuneBlock,
		Color:    tc.StyleDefault.Foreground(tc.ColorGray).Background(tc.ColorDarkGoldenrod),
		Current:  m,
		IsActive: false,
	}
}

func (c *Cursor) Draw(s tc.Screen) {
	ut.EmitStr(s, c.X, c.Y, c.Color, string(c.Rune))
}

func (c *Cursor) SetCurrentFloor(f int, m maps.Map) {
	c.Current = &m
	c.floor = f
}

func (c *Cursor) SetRune(r rune) {
	c.Rune = r
}

func (c *Cursor) Look(a *actors.Actors) string {
	var seen string
	var selected *maps.MapCell
	var items *items.Items
	var actor_desc string
	var item_desc string
	var current_rune rune

	actor := a.GetActorFromLocaion(c.X, c.Y)
	actor_desc = ""
	item_desc = ""

	if c.X < 0 || c.Y < 0 || c.X >= c.Current.GetWidth() || c.Y >= c.Current.GetHeight() {
		return "You don't see anything."
	} else {
		selected = c.Current.GetCell(c.X, c.Y)
		items = selected.GetItems()
		current_rune = selected.GetRune()

		if len(*items) > 1 {
			item_desc = "There is a stack of items on the ground here."
		} else if len(*items) == 1 {
			item_desc = selected.GetFirstItem().GetDescription()
		}

		if len(*items) > 0 {
			current_rune = selected.GetFirstItem().GetRune()
		}

		if actor != nil && actor.Floor == c.floor {
			actor_desc = actor.Type.GetDescription()
			current_rune = actor.Type.GetRune()
		}

		seen = fmt.Sprintf("%s %s %s", actor_desc, item_desc, selected.GetDescription())

	}

	c.SetRune(current_rune)

	return seen
}

func (c *Cursor) Move(x, y int) {
	w := c.Current.GetWidth()
	h := c.Current.GetHeight()

	if c.X+x < 0 || c.Y+y < 0 || c.X+x >= w || c.Y+y >= h {
		return
	}

	c.X += x
	c.Y += y
}
