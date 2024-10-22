package actors

import (
	"fmt"

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

func (a *Actor) Move(x, y int, s tc.Screen) {
	// TODO: add more detailed collision detection
	l, _, _, _ := s.GetContent(a.X+x, a.Y+y)

	style := tc.StyleDefault.Foreground(tc.ColorRed).Background(tc.ColorBlack)

	if a.Type.GetRune() == '@' {
		ut.EmitStr(s, 30, 10, style, fmt.Sprintf("In Move before update: l:%s (%d,%d)", string(l), a.X, a.Y))
	}

	if l != '#' {
		a.X += x
		a.Y += y
	}

	if a.Type.GetRune() == '@' {
		ut.EmitStr(s, 30, 11, style, fmt.Sprintf("In Move after update: l:%s (%d,%d)", string(l), a.X, a.Y))
	}
}

func (a *Actor) Wander(roll, f int, s tc.Screen) {
	/*
		Random Wander
		1 = Up
		2 = Right
		3 = Down
		4 = Left
	*/

	if f == a.Floor {
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

func (a *Actor) Draw(f int, s tc.Screen) {
	if f == a.Floor {
		ut.EmitStr(s, a.X, a.Y, a.Color, string(a.Type.GetRune()))
	}
}

func (a *Actor) GetLocation() (z, x, y int) {
	return a.Floor, a.X, a.Y
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
