package main

import (
	"fmt"
	"os"
	"time"

	"bretbelgarde.com/adventure/creatures"
	ut "bretbelgarde.com/adventure/utils"
	tc "github.com/gdamore/tcell/v2"
	//rw "github.com/mattn/go-runewidth"
	//tv "github.com/rivo/tview"
)

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
	var msg string

	debug := false
	player := player{rune: 'P', x: 0, y: 0, health: 10, level: 1}

	mapp := [9][9]rune{
		{'#', '#', '#', '#', '#', '#', '#', '#', '#'},
		{'#', '.', '.', '.', '.', '.', '.', '.', '#'},
		{'#', '.', '.', '.', '.', '.', '.', '.', '#'},
		{'#', '.', '.', '.', '#', '.', '.', '.', '#'},
		{'#', '.', '.', '#', '#', '#', '.', '.', '#'},
		{'#', '.', '.', '.', '#', '.', '.', '.', '#'},
		{'#', '.', '.', '.', '.', '.', '.', '.', '#'},
		{'#', '.', '.', '.', '.', '.', '.', '.', '#'},
		{'#', '#', '#', '#', '#', '#', '#', '#', '#'},
	}

	mapp2 := [9][9]rune{
		{'#', '#', '#', '#', '#', '#', '#', '#', '#'},
		{'#', '#', '.', '.', '.', '.', '.', '#', '#'},
		{'#', '.', '.', '.', '.', '.', '.', '.', '#'},
		{'#', '.', '.', '.', '.', '.', '.', '.', '#'},
		{'#', '.', '.', '.', '#', '.', '.', '.', '#'},
		{'#', '.', '.', '.', '.', '.', '.', '.', '#'},
		{'#', '.', '.', '.', '.', '.', '<', '.', '#'},
		{'#', '#', '.', '.', '.', '.', '.', '#', '#'},
		{'#', '#', '#', '#', '#', '#', '#', '#', '#'},
	}

	level := 1
	current := mapp

	tc.SetEncodingFallback(tc.EncodingFallbackASCII)

	s, e := tc.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	white := tc.StyleDefault.
		Foreground(tc.ColorWhite).
		Background(tc.ColorBlack)

	grey := tc.StyleDefault.
		Foreground(tc.ColorGrey).
		Background(tc.ColorBlack)

	burlyWood := tc.StyleDefault.
		Foreground(tc.ColorBurlyWood).
		Background(tc.ColorBlack)

	brown := tc.StyleDefault.
		Foreground(tc.ColorBrown).
		Background(tc.ColorBlack)

	s.SetStyle(tc.StyleDefault.
		Foreground(tc.ColorWhite).
		Background(tc.ColorBlack))
	s.EnableMouse()
	s.Clear()

	a := &Actors{}
	a.Actors = append(a.Actors, NewActor(4, 3, 1, brown, &creatures.Pig{Rune: 'p', Health: 10, Description: "Pig 1"}))
	a.Actors = append(a.Actors, NewActor(5, 2, 1, brown, &creatures.Pig{Rune: 'p', Health: 10, Description: "Pig 2"}))

	a.Actors = append(a.Actors, NewActor(4, 3, 2, brown, &creatures.Rat{Rune: 'r', Health: 10, Description: "Rat 1"}))
	a.Actors = append(a.Actors, NewActor(5, 2, 2, brown, &creatures.Rat{Rune: 'r', Health: 10, Description: "Rat 2"}))
	a.Actors = append(a.Actors, NewActor(4, 3, 2, brown, &creatures.Rat{Rune: 'r', Health: 10, Description: "Rat 3"}))
	a.Actors = append(a.Actors, NewActor(6, 3, 2, brown, &creatures.Rat{Rune: 'r', Health: 10, Description: "Rat 4"}))

	quit := make(chan struct{})
	go func() {
		for {
			x, y := s.Size()
			ev := s.PollEvent()

			switch ev := ev.(type) {
			case *tc.EventKey:
				switch ev.Key() {
				case tc.KeyRune:
					switch ev.Rune() {
					case ':':
						g := current[player.x-1][player.y-1]
						if g == '.' {
							msg = "You see some dirt."
						} else if g == '>' {
							msg = "You see stairs down"
						} else if g == '<' {
							msg = "You see stairs up"
						}

					case '>':
						g := current[player.x-1][player.y-1]
						if g == '>' {
							level++
							s.Clear()
						}

					case '<':
						g := current[player.x-1][player.y-1]
						if g == '<' {
							level--
							s.Clear()
						}

					case 'h':
						r, _, _, _ := s.GetContent(player.x-1, player.y)
						if r == '#' {
							// Do a thing
						} else if player.x-1 >= 0 {
							player.x--
						}
					case 'l':
						r, _, _, _ := s.GetContent(player.x+1, player.y)
						if r == '#' {
							// Do a thing
						} else if player.x+1 < x {
							player.x++
						}
					case 'k':
						r, _, _, _ := s.GetContent(player.x, player.y-1)
						if r == '#' {
							// Do a thing
						} else if player.y-1 >= 0 {
							player.y--
						}
					case 'j':
						r, _, _, _ := s.GetContent(player.x, player.y+1)
						if r == '#' {
							// Do a thing
						} else if player.y+1 < y {
							player.y++
						}
					}

				case tc.KeyEscape, tc.KeyEnter:
					close(quit)
					return
				case tc.KeyRight:
					r, _, _, _ := s.GetContent(player.x+1, player.y)
					if r == '#' {
						// Do a thing
					} else if player.x+1 < x {
						player.x++
					}
				case tc.KeyLeft:
					r, _, _, _ := s.GetContent(player.x-1, player.y)
					if r == '#' {
						// Do a thing
					} else if player.x-1 >= 0 {
						player.x--
					}
				case tc.KeyUp:
					r, _, _, _ := s.GetContent(player.x, player.y-1)
					if r == '#' {
						// Do a thing
					} else if player.y-1 >= 0 {
						player.y--
					}
				case tc.KeyDown:
					r, _, _, _ := s.GetContent(player.x, player.y+1)
					if r == '#' {
						// Do a thing
					} else if player.y+1 < y {
						player.y++
					}
				case tc.KeyCtrlD:
					debug = !debug
				case tc.KeyCtrlL:
					s.Clear()
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
		//s.Clear()

		dbg := fmt.Sprintf("player x: %d y: %d", player.x, player.y)
		if debug {
			var yy int
			if player.y == 0 {
				_, yy = s.Size()
				yy--
			} else {
				yy = 0
			}

			ut.EmitStr(s, 0, yy, white, dbg)
		}

		if level == 1 {
			current = mapp
		} else if level == 2 {
			current = mapp2
		}

		var color tc.Style
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if mapp[i][j] == '#' {
					color = grey
				}

				if mapp[i][j] == '.' {
					color = burlyWood
				}
				ut.EmitStr(s, i+1, j+1, color, string(current[i][j]))
			}
		}

		ut.EmitStr(s, 0, 0, white, msg)

		for i := range a.Actors {
			a.Actors[i].Draw(s, level)
		}

		ut.EmitStr(s, player.x, player.y, white, string(player.rune))
		s.Show()
	}

	s.Fini()
}
