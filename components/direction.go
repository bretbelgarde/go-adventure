package components

type Direction struct {
	DX int `json:"dx"`
	DY int `json:"dy"`
}

func (d *Direction) Mask() uint64 {
	return MaskDirection
}

func (d *Direction) WithX(dx int) *Direction {
	d.DX = dx
	return d
}

func (d *Direction) WithY(dy int) *Direction {
	d.DY = dy
	return d
}

func NewDirection() *Direction {
	return &Direction{}
}
