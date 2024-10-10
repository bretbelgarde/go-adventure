package main

import (
	cr "bretbelgarde.com/adventure/creatures"
	"bretbelgarde.com/adventure/die"
	ut "bretbelgarde.com/adventure/utils"
	tc "github.com/gdamore/tcell/v2"
)

type Actors struct {
	Actors []Actor
}

type Actor struct {
	ID       string
	X        int
	Y        int
	Floor    int
	Color    tc.Style
	Creature cr.Creature
}

func (a *Actor) Move(floor int, s tc.Screen) {
	/*
		Random Wander
		1 = Up
		2 = Right
		3 = Down
		4 = Left
	*/

	if floor == a.Floor {
		var xx, yy int
		d, err := die.Roll("1d4")

		if err != nil {
			panic("Die roll failed")
		}

		switch d {
		case 1:
			yy = -1
		case 2:
			xx = 1
		case 3:
			yy = 1
		case 4:
			xx = -1
		}

		l, _, _, _ := s.GetContent(a.X+xx, a.Y+yy)

		if l != '#' {
			a.X += xx
			a.Y += yy
		}
	}
}

func (a *Actor) Draw(s tc.Screen, f int) {
	if f == a.Floor {
		ut.EmitStr(s, a.X, a.Y, a.Color, string(a.Creature.GetRune()))
	}
}

func (a *Actor) GetLocation() (x, y int) {
	return a.X, a.Y
}

func NewActor(x, y, floor int, color tc.Style, creature cr.Creature) Actor {
	return Actor{
		X:        x,
		Y:        y,
		Floor:    floor,
		Color:    color,
		Creature: creature,
	}
}

type player struct {
	rune   rune
	x      int
	y      int
	health int
	level  int
}
