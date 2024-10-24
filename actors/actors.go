package actors

import (
	"bretbelgarde.com/adventure/maps"
	ut "bretbelgarde.com/adventure/utils"
	tc "github.com/gdamore/tcell/v2"
)

type Actors []*Actor

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

func (a *Actor) Move(m maps.Map, x, y int) {
	t := m.GetCell(a.X+x, a.Y+y).IsTraversable()

	if t {
		a.X += x
		a.Y += y
	}
}

func (a *Actor) Wander(m maps.Map, roll, f int) {
	/*
		Random Wander
		1 = Up
		2 = Right
		3 = Down
		4 = Left
		5 = Stay
	*/

	if f == a.Floor {
		switch roll {
		case 1:
			a.Move(m, 0, -1)
		case 2:
			a.Move(m, 1, 0)
		case 3:
			a.Move(m, 0, 1)
		case 4:
			a.Move(m, -1, 0)
		}
	}
}

func (a *Actor) Draw(s tc.Screen, f int) {
	if f == a.Floor {
		ut.EmitStr(s, a.X, a.Y, a.Color, string(a.Type.GetRune()))
	}
}

func (a *Actor) GetLocation() (x, y, z int) {
	return a.X, a.Y, a.Floor
}

func NewActor(x, y, floor int, color tc.Style, actor ActorType) *Actor {
	return &Actor{
		X:     x,
		Y:     y,
		Floor: floor,
		Color: color,
		Type:  actor,
	}
}
