package components

type Symbol struct {
	Symbol rune   `json:"symbol"`
	Fg     string `json:"fg"`
	Bg     string `json:"bg"`
}

func (s *Symbol) Mask() uint64 {
	return MaskSymbol
}

func (s *Symbol) WithSymbol(symbol rune) *Symbol {
	s.Symbol = symbol
	return s
}

func (s *Symbol) Foreground(fg string) *Symbol {
	s.Fg = fg
	return s
}

func (s *Symbol) Background(bg string) *Symbol {
	s.Bg = bg
	return s
}

func NewSymbol() *Symbol {
	return &Symbol{}
}
