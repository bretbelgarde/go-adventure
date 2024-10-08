package creatures

type Rat struct {
	Rune        rune   `json:"r,omitempty"`
	Health      int    `json:"health,omitempty"`
	Description string `json:"description,omitempty"`
}

func (p *Rat) GetRune() rune {
	return p.Rune
}

func (p *Rat) GetHealth() int {
	return p.Health
}

func (p *Rat) GetDescription() string {
	return p.Description
}

func (p *Rat) TakeDamage(dmg int) {
	p.Health -= dmg
}
