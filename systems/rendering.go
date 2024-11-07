package systems

import (
	"bretbelgarde.com/adventure/components"
	"github.com/andygeiss/ecs"
	"github.com/gdamore/tcell/v2"
)

type renderingSystem struct {
	err    error
	screen tcell.Screen
	width  int
	height int
}

func (r *renderingSystem) Error() error {
	return r.err
}

func (r *renderingSystem) Setup() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	r.screen, r.err = tcell.NewScreen()
	r.err = r.screen.Init()

	r.screen.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))

}

func (r *renderingSystem) Process(em ecs.EntityManager) (state int) {
	r.screen.Clear()
	ev := r.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventResize:
		r.screen.Sync()
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			return ecs.StateEngineStop
		}
	}
	r.renderEntities(em)
	r.screen.Show()

	return ecs.StateEngineContinue
}

func (r *renderingSystem) Teardown() {

	r.screen.Fini()
}

func (r *renderingSystem) Width(w int) *renderingSystem {
	r.width = w
	return r
}

func (r *renderingSystem) Height(h int) *renderingSystem {
	r.height = h
	return r
}

func (r *renderingSystem) renderEntities(em ecs.EntityManager) {
	for _, e := range em.FilterByMask(components.MaskPosition | components.MaskSymbol) {
		posistion := e.Get(components.MaskPosition).(*components.Position)
		symbol := e.Get(components.MaskSymbol).(*components.Symbol)
		style := tcell.StyleDefault.
			Foreground(tcell.GetColor(symbol.Fg)).
			Background(tcell.GetColor(symbol.Bg))

		r.screen.SetContent(posistion.X, posistion.Y, symbol.Symbol, nil, style)
	}
}

func NewRenderingSystem() *renderingSystem {
	return &renderingSystem{}
}
