package items

import tc "github.com/gdamore/tcell/v2"

type Items []Item

type Item struct {
	ID          string
	Rune        rune
	Description string
	Color       tc.Style
}

func (i *Item) GetRune() rune {
	return i.Rune
}

func (i *Item) GetDescription() string {
	return i.Description
}

func (i *Item) GetColor() tc.Style {
	return i.Color
}

func (i *Item) GetID() string {
	return i.ID
}
