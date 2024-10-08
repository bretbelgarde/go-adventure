package creatures

type Creature interface {
	GetRune() rune
	GetHealth() int
	GetDescription() string
	TakeDamage(int)
}
