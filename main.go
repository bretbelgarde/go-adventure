package main

import (
	"fmt"
	"os"
	"time"

	//"bretbelgarde.com/adventure/creatures"
	"bretbelgarde.com/adventure/player"
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

	mapp := [9][9]rune{
		{'#', '#', '#', '#', '#', '#', '#', '#', '#'},
		{'#', '.', '.', '.', '.', '.', '.', '.', '#'},
		{'#', '.', '.', '.', '.', '.', '.', '.', '#'},
		{'#', '.', '.', '.', '#', '.', '.', '.', '#'},
		{'#', '.', '.', '#', '#', '#', '.', '.', '#'},
		{'#', '.', '.', '.', '#', '.', '.', '.', '#'},
		{'#', '.', '.', '.', '.', '.', '>', '.', '#'},
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

	/* brown := tc.StyleDefault.
	Foreground(tc.ColorBrown).
	Background(tc.ColorBlack) */

	s.SetStyle(tc.StyleDefault.
		Foreground(tc.ColorWhite).
		Background(tc.ColorBlack))
	s.EnableMouse()
	s.Clear()

	player := NewActor(3, 3, 1, white, &player.Player{Rune: '@', Health: 10, Description: "Player"})

	/* a := &Actors{}

	a.Actors = append(a.Actors, NewActor(4, 3, 1, brown, &creatures.Pig{Rune: 'p', Health: 10, Description: "Pig 1"}))
	a.Actors = append(a.Actors, NewActor(5, 2, 1, brown, &creatures.Pig{Rune: 'p', Health: 10, Description: "Pig 2"}))

	a.Actors = append(a.Actors, NewActor(4, 3, 2, brown, &creatures.Rat{Rune: 'r', Health: 10, Description: "Rat 1"}))
	a.Actors = append(a.Actors, NewActor(5, 2, 2, brown, &creatures.Rat{Rune: 'r', Health: 10, Description: "Rat 2"}))
	a.Actors = append(a.Actors, NewActor(5, 4, 2, brown, &creatures.Rat{Rune: 'r', Health: 10, Description: "Rat 3"}))
	a.Actors = append(a.Actors, NewActor(6, 3, 2, brown, &creatures.Rat{Rune: 'r', Health: 10, Description: "Rat 4"})) */

	quit := make(chan struct{})
	go func() {
		for {
			//x, y := s.Size()
			ev := s.PollEvent()

			switch ev := ev.(type) {
			case *tc.EventKey:
				switch ev.Key() {
				case tc.KeyRune:
					switch ev.Rune() {
					case ':':
						g := current[player.X-1][player.Y-1]
						if g == '.' {
							msg = "You see some dirt."
						} else if g == '>' {
							msg = "You see stairs down"
						} else if g == '<' {
							msg = "You see stairs up"
						}
					case '>':
						g := current[player.X-1][player.Y-1]
						if g == '>' {
							level++
							s.Clear()
						}
					case '<':
						g := current[player.X-1][player.Y-1]
						if g == '<' {
							level--
							s.Clear()
						}
					case 'h':
						player.Move(-1, 0, s)
					case 'l':
						player.Move(1, 0, s)
					case 'k':
						player.Move(0, -1, s)
					case 'j':
						player.Move(0, 1, s)
					}

				case tc.KeyEscape, tc.KeyEnter:
					close(quit)
					return

				case tc.KeyRight:
					player.Move(1, 0, s)

				case tc.KeyLeft:
					player.Move(-1, 0, s)

				case tc.KeyUp:
					player.Move(0, -1, s)

				case tc.KeyDown:
					player.Move(0, 1, s)

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
		s.Clear()

		dbg := fmt.Sprintf("player x: %d y: %d", player.X, player.Y)
		if debug {
			var yy int
			if player.Y == 0 {
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

		ut.EmitStr(s, player.X, player.Y, white, string(player.Type.GetRune()))
		s.Show()
	}

	s.Fini()
}
