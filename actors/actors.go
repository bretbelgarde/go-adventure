package actors

import (
	ut "bretbelgarde.com/adventure/utils"
	tc "github.com/gdamore/tcell/v2"
)

type Actors []Actor

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
	// TODO: add more detailed collision detection
	l, _, _, _ := s.GetContent(a.X+x, a.Y+y)

	if l != '#' {
		a.X += x
		a.Y += y
	}
}

func (a *Actor) Wander(roll, floor int, s tc.Screen) {
	/*
		Random Wander
		1 = Up
		2 = Right
		3 = Down
		4 = Left
	*/

	if floor == a.Floor {
		switch roll {
		case 1:
			a.Move(0, -1, s)
		case 2:
			a.Move(1, 0, s)
		case 3:
			a.Move(0, 1, s)
		case 4:
			a.Move(-1, 0, s)
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
