package main

import (
	"fmt"
	"os"
	"time"

	tc "github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	//tv "github.com/rivo/tview"
)

type player struct {
	x int
	y int
}

func emit_str(s tc.Screen, x, y int, style tc.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func main() {
	/* app := tv.NewApplication()

	newPrimitive := func(text string, align int) tv.Primitive {
		return tv.NewTextView().
			SetTextAlign(align).
			SetText(text)
	}

	main_win := tv.NewTextArea()
	sidebar := newPrimitive("Sidebar", tv.AlignCenter)
	status := newPrimitive("Status Bar", tv.AlignCenter)

	grid := tv.NewGrid().
		SetRows(0, 1).
		SetColumns(0, 20).
		SetBorders(true).
		AddItem(main_win, 0, 0, 1, 1, 0, 0, true).
		AddItem(sidebar, 0, 1, 1, 1, 0, 0, false).
		AddItem(status, 1, 0, 1, 2, 0, 0, false)

	err := app.SetRoot(grid, true).SetFocus(grid).Run()

	if err != nil {
		panic(err)
	} */

	tc.SetEncodingFallback(tc.EncodingFallbackASCII)

	s, e := tc.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
	}

	white := tc.StyleDefault.
		Foreground(tc.ColorWhite).
		Background(tc.ColorBlack)

	s.SetStyle(tc.StyleDefault.
		Foreground(tc.ColorWhite).
		Background(tc.ColorBlack))
	s.Clear()

	quit := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tc.EventKey:
				switch ev.Key() {
				case tc.KeyEscape, tc.KeyEnter:
					close(quit)
					return
				case tc.KeyCtrlL:
					s.Sync()
				}
			case *tc.EventResize:
				s.Sync()
			}
		}
	}()

loop:
	for {
		select {
		case <-quit:
			break loop
		case <-time.After(time.Millisecond * 50):
		}
		s.Clear()
		emit_str(s, 0, 0, white, "@")
		s.Show()
	}

	s.Fini()
}
