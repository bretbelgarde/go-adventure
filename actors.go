package main

import (
	"bretbelgarde.com/adventure/die"
	ut "bretbelgarde.com/adventure/utils"
	tc "github.com/gdamore/tcell/v2"
)

type Actors struct {
	Actors []Actor
}

type ActorType interface {
	GetRune() rune
	GetHealth() int
	GetDescription() string
	TakeDamage(int)
}

type Actor struct {
	ID    string
	X     int
	Y     int
	Floor int
	Color tc.Style
	Type  ActorType
}

func (a *Actor) Move(x, y int, s tc.Screen) {
	l, _, _, _ := s.GetContent(a.X+x, a.Y+y)

	if l != '#' {
		a.X += x
		a.Y += y
	}
}

func (a *Actor) Wander(floor int, s tc.Screen) {
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
		ut.EmitStr(s, a.X, a.Y, a.Color, string(a.Type.GetRune()))
	}
}

func (a *Actor) GetLocation() (x, y int) {
	return a.X, a.Y
}

func NewActor(x, y, floor int, color tc.Style, actor ActorType) Actor {
	return Actor{
		X:     x,
		Y:     y,
		Floor: floor,
		Color: color,
		Type:  actor,
	}
}
