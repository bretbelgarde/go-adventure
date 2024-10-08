package creatures

type Rat struct {
	Rune        rune   `json:"r,omitempty"`
	Health      int    `json:"health,omitempty"`
	Description string `json:"description,omitempty"`
}

func (r *Rat) GetRune() rune {
	return r.Rune
}

func (r *Rat) GetHealth() int {
	return r.Health
}

func (r *Rat) GetDescription() string {
	return r.Description
}

func (r *Rat) TakeDamage(dmg int) {
	r.Health -= dmg
}
