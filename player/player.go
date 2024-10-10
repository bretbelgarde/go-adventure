package player

type Player struct {
	Rune        rune
	Health      int
	Description string
}

func (p *Player) GetRune() rune {
	return p.Rune
}

func (p *Player) GetHealth() int {
	return p.Health
}

func (p *Player) GetDescription() string {
	return p.Description
}

func (p *Player) TakeDamage(dmg int) {
	p.Health -= dmg
}
