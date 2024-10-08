package creatures

type Pig struct {
	Rune        rune   `json:"r,omitempty"`
	Health      int    `json:"health,omitempty"`
	Description string `json:"description,omitempty"`
}

func (p *Pig) GetRune() rune {
	return p.Rune
}

func (p *Pig) GetHealth() int {
	return p.Health
}

func (p *Pig) GetDescription() string {
	return p.Description
}

func (p *Pig) TakeDamage(dmg int) {
	p.Health -= dmg
}
