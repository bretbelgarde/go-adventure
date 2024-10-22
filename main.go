package main

import (
	"fmt"
	"os"

	"bretbelgarde.com/adventure/actors"
	cr "bretbelgarde.com/adventure/creatures"
	"bretbelgarde.com/adventure/items"
	"bretbelgarde.com/adventure/maps"
	"bretbelgarde.com/adventure/player"
	ut "bretbelgarde.com/adventure/utils"
	tc "github.com/gdamore/tcell/v2"
	"golang.org/x/exp/rand"
)

func main() {
	var msg string
	debug := false

	wall := maps.MapCell{
		Rune:        '#',
		Traversable: false,
		Description: "A rough-hewn stone wall.",
		Color: tc.StyleDefault.
			Foreground(tc.ColorGray).
			Background(tc.ColorBlack),
	}

	ground := maps.MapCell{
		Rune:        '.',
		Traversable: true,
		Description: "A hard-packed dirt floor.",
		Color: tc.StyleDefault.
			Foreground(tc.ColorBurlyWood).
			Background(tc.ColorBlack),
	}

	ground_with_gold := maps.MapCell{
		Rune:        '.',
		Traversable: true,
		Description: "A hard-packed dirt floor.",
		Color: tc.StyleDefault.
			Foreground(tc.ColorBurlyWood).
			Background(tc.ColorBlack),
		Items: items.Items{
			items.Item{
				ID:          "gold",
				Description: "a pile of gold coins",
				Rune:        '$',
				Color: tc.StyleDefault.
					Foreground(tc.ColorGold).
					Background(tc.ColorBlack),
			},
		},
	}

	down_stairs := maps.MapCell{
		Rune:        '>',
		Traversable: true,
		Description: "A maze of twisty stairs leading down.",
		Color: tc.StyleDefault.
			Foreground(tc.ColorBrown).
			Background(tc.ColorBlack),
	}

	up_stairs := maps.MapCell{
		Rune:        '<',
		Traversable: true,
		Description: "A maze of twisty stairs leading up.",
		Color: tc.StyleDefault.
			Foreground(tc.ColorBrown).
			Background(tc.ColorBlack),
	}

	var dungeon maps.Floors

	dungeon = append(dungeon, maps.Map{
		{wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, ground, ground, ground, ground, ground, ground, ground, wall},
		{wall, ground, ground, ground, ground, ground, ground, ground, wall},
		{wall, ground, ground, ground, wall, ground, ground, ground, wall},
		{wall, ground, ground, wall, wall, wall, ground, ground, wall},
		{wall, ground, ground_with_gold, ground, wall, ground, ground, ground, wall},
		{wall, ground, ground, ground, ground, ground, down_stairs, ground, wall},
		{wall, ground, ground, ground, ground, ground, ground, ground, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, wall},
	}, maps.Map{
		{wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, ground, ground, ground, ground, ground, wall, wall},
		{wall, ground, ground, ground, ground, ground, ground, ground, wall},
		{wall, ground, ground, ground, ground, ground, ground, ground, wall},
		{wall, ground, ground, ground, wall, ground, ground, ground, wall},
		{wall, ground, ground, ground, ground, ground, ground, ground, wall},
		{wall, ground, ground, ground, ground, ground, up_stairs, ground, wall},
		{wall, wall, ground, ground, ground, ground, ground, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, wall},
	})

	level := 0
	current := dungeon[level]

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

	brown := tc.StyleDefault.
		Foreground(tc.ColorBrown).
		Background(tc.ColorBlack)

	pink := tc.StyleDefault.
		Foreground(tc.ColorPink).
		Background(tc.ColorBlack)

	s.SetStyle(tc.StyleDefault.
		Foreground(tc.ColorWhite).
		Background(tc.ColorBlack))

	s.EnableMouse()
	s.Clear()

	player := actors.NewActor(3, 3, 0, white, &player.Player{Rune: '@', Health: 10, Description: "Player"})

	var creatures actors.Actors
	creatures = append(
		creatures,
		actors.NewActor(3, 8, 0, pink, &cr.Pig{Rune: 'p', Health: 5, Description: "A Pig who loves straw"}),
		actors.NewActor(5, 8, 0, pink, &cr.Pig{Rune: 'p', Health: 5, Description: "A Pig who loves sticks"}),
		actors.NewActor(5, 8, 0, pink, &cr.Pig{Rune: 'p', Health: 5, Description: "A Pig who loves bricks"}),
		actors.NewActor(5, 3, 1, white, &cr.Rat{Rune: 'r', Health: 10, Description: "Lab Rat"}),
		actors.NewActor(5, 5, 1, brown, &cr.Rat{Rune: 'r', Health: 10, Description: "You Dirty Rat"}),
		actors.NewActor(4, 4, 1, white, &cr.Rat{Rune: 'r', Health: 10, Description: "Ratt *Plays Guitar Riff*"}),
		actors.NewActor(6, 4, 1, brown, &cr.Rat{Rune: 'r', Health: 10, Description: "Rat-tatooee"}),
	)

	quit := func() {
		s.Fini()
		os.Exit(0)
	}

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
					if current[player.X-1][player.Y-1].Items != nil {
						msg = current[player.X-1][player.Y-1].GetFirstItem().Description
					} else {
						msg = current[player.X-1][player.Y-1].GetDescription()
					}

				case '>':
					g := current[player.X-1][player.Y-1].GetRune()
					if g == '>' {
						level++
						s.Clear()
					}
				case '<':
					g := current[player.X-1][player.Y-1].GetRune()
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
		default:
			continue

		}

		// Creature Movement
		for _, c := range creatures {
			if c.Floor == player.Floor {
				c.Wander(rand.Intn(4)+1, level, s)

			}
		}

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

		if level == 0 {
			current = dungeon[level]
			player.Floor = 0
		} else if level == 1 {
			current = dungeon[level]
			player.Floor = 1
		}

		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				var color tc.Style
				var map_rune rune

				if current[i][j].Items != nil {
					// If there are items in the MapCell grab the first item's color and rune
					color = current[i][j].GetFirstItem().GetColor()
					map_rune = current[i][j].GetFirstItem().GetRune()
				} else {
					color = current[i][j].GetColor()
					map_rune = current[i][j].GetRune()
				}

				ut.EmitStr(s, i+1, j+1, color, string(map_rune))
			}
		}

		for _, c := range creatures {
			if c.Floor == player.Floor {
				c.Draw(level, s)
			}
		}

		player.Draw(level, s)

		if msg != "" {
			ut.EmitStr(s, 0, 0, white, msg)
		}

	}
}
