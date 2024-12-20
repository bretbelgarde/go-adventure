package maps

import (
	"bretbelgarde.com/adventure/items"
	tc "github.com/gdamore/tcell/v2"
)

type Floors []Map
type Map [][]MapCell

func (m *Map) GetCell(x, y int) *MapCell {
	return &(*m)[y][x]
}
func (m *Map) GetHeight() int {
	return len(*m)
}

func (m *Map) GetWidth() int {
	return len((*m)[0])
}

type MapCell struct {
	Rune        rune   `json:"rune,omitempty"`
	Traversable bool   `json:"traversable,omitempty"`
	Description string `json:"description,omitempty"`
	Color       tc.Style
	Items       items.Items
}

func (mc *MapCell) GetRune() rune {
	return mc.Rune
}

func (mc *MapCell) GetDescription() string {
	return mc.Description
}

func (mc *MapCell) IsTraversable() bool {
	return mc.Traversable
}

func (mc *MapCell) GetColor() tc.Style {
	return mc.Color
}

func (mc *MapCell) GetFirstItem() *items.Item {
	return &mc.Items[0]
}

func (mc *MapCell) GetItems() *items.Items {
	return &mc.Items
}

func (mc *MapCell) GetItemFromStack(index int) *items.Item {
	return &mc.Items[index]
}

func (mc *MapCell) AddItem(item items.Item) {
	mc.Items = append(mc.Items, item)
}

func (mc *MapCell) RemoveItem(index int) {
	mc.Items = append(mc.Items[:index], mc.Items[index+1:]...)
}
