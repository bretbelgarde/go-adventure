package main

import (
	"fmt"
	"os"

	"bretbelgarde.com/adventure/actors"
	cr "bretbelgarde.com/adventure/creatures"
	"bretbelgarde.com/adventure/player"
	ut "bretbelgarde.com/adventure/utils"
	tc "github.com/gdamore/tcell/v2"
	"golang.org/x/exp/rand"
)

func main() {
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

	brown := tc.StyleDefault.
		Foreground(tc.ColorBrown).
		Background(tc.ColorBlack)

	s.SetStyle(tc.StyleDefault.
		Foreground(tc.ColorWhite).
		Background(tc.ColorBlack))
	s.EnableMouse()
	s.Clear()

	player := actors.NewActor(3, 3, 1, white, &player.Player{Rune: '@', Health: 10, Description: "Player"})
	pig := actors.NewActor(3, 8, 1, brown, &cr.Pig{Rune: 'p', Health: 10, Description: "Pig 1"})

	quit := func() {
		s.Fini()
		os.Exit(0)
	}

	// Event Loop
	for {
		// Update Screen
		s.Show()

		// Poll Event
		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tc.EventResize:
			s.Sync()

		case *tc.EventKey:
			switch ev.Key() {
			case tc.KeyRune:
				switch ev.Rune() {
				case ':':
					g := current[player.X-1][player.Y-1]
					if g == '.' {
						msg = "You see some dirt."
					} else if g == '>' {
						msg = "You see stairs down."
					} else if g == '<' {
						msg = "You see stairs up."
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
				}
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
			case tc.KeyCtrlC, tc.KeyEscape:
				quit()
			}
		}
		// Creature Movement
		pig.Wander(rand.Intn(4)+1, level, s)

		// Process Event
		dbg := fmt.Sprintf("player level: %d x: %d y: %d", player.Floor, player.X, player.Y)

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
			player.Floor = 1
		} else if level == 2 {
			current = mapp2
			player.Floor = 2
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

		player.Draw(s, level)
		pig.Draw(s, level)

		if msg != "" {
			ut.EmitStr(s, 0, 0, white, msg)
		}

	}
}

// 	go func() {
// 		for {
// 			ev := s.PollEvent()

// 			switch ev := ev.(type) {
// 			case *tc.EventKey:
// 				switch ev.Key() {
// 				case tc.KeyRune:
// 					switch ev.Rune() {
// 					case ':':
// 						g := current[player.X-1][player.Y-1]
// 						if g == '.' {
// 							msg = "You see some dirt."
// 						} else if g == '>' {
// 							msg = "You see stairs down."
// 						} else if g == '<' {
// 							msg = "You see stairs up."
// 						}
// 					case '>':
// 						g := current[player.X-1][player.Y-1]
// 						if g == '>' {
// 							level++
// 							s.Clear()
// 						}
// 					case '<':
// 						g := current[player.X-1][player.Y-1]
// 						if g == '<' {
// 							level--
// 							s.Clear()
// 						}
// 					}

// 				case tc.KeyEscape, tc.KeyEnter:
// 					close(quit)
// 					return

// 				case tc.KeyRight:
// 					player.Move(1, 0, s)
// 				case tc.KeyLeft:
// 					player.Move(-1, 0, s)
// 				case tc.KeyUp:
// 					player.Move(0, -1, s)
// 				case tc.KeyDown:
// 					player.Move(0, 1, s)
// 				case tc.KeyCtrlD:
// 					debug = !debug
// 				case tc.KeyCtrlL:
// 					s.Clear()
// 					s.Sync()
// 				}
// 			case *tc.EventResize:
// 				s.Sync()
// 			}
// 		}
// 	}()

// loop:
// 	for {
// 		select {
// 		case <-quit:
// 			break loop
// 		case <-time.After(time.Millisecond * 50):
// 		}
// 		s.Clear()

// 		dbg := fmt.Sprintf("player level: %d x: %d y: %d", player.Floor, player.X, player.Y)

// 		if debug {
// 			var yy int
// 			if player.Y == 0 {
// 				_, yy = s.Size()
// 				yy--
// 			} else {
// 				yy = 0
// 			}

// 			ut.EmitStr(s, 0, yy, white, dbg)
// 		}

// 		if level == 1 {
// 			current = mapp
// 			player.Floor = 1
// 		} else if level == 2 {
// 			current = mapp2
// 			player.Floor = 2
// 		}

// 		var color tc.Style
// 		for i := 0; i < 9; i++ {
// 			for j := 0; j < 9; j++ {
// 				if mapp[i][j] == '#' {
// 					color = grey
// 				}

// 				if mapp[i][j] == '.' {
// 					color = burlyWood
// 				}

// 				ut.EmitStr(s, i+1, j+1, color, string(current[i][j]))
// 			}
// 		}

// 		player.Draw(s, level)
// 		pig.Draw(s, level)

// 		if msg != "" {
// 			ut.EmitStr(s, 0, 0, white, msg)
// 		}

// 		s.Show()
// 	}
