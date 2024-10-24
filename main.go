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

type Game struct {
	screen    tc.Screen
	debug     bool
	msg       string
	floor     int
	dungeon   maps.Floors
	player    *actors.Actor
	creatures actors.Actors
}

func (g *Game) Init() {
	var err error
	g.debug = false

	tc.SetEncodingFallback(tc.EncodingFallbackASCII)

	g.screen, err = tc.NewScreen()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err = g.screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func main() {
	var g Game
	g.Init()

	/* Colors */
	white := tc.StyleDefault.
		Foreground(tc.ColorWhite).
		Background(tc.ColorBlack)

	brown := tc.StyleDefault.
		Foreground(tc.ColorBrown).
		Background(tc.ColorBlack)

	pink := tc.StyleDefault.
		Foreground(tc.ColorPink).
		Background(tc.ColorBlack)

	gray := tc.StyleDefault.
		Foreground(tc.ColorGray).
		Background(tc.ColorBlack)

	burlyWood := tc.StyleDefault.
		Foreground(tc.ColorBurlyWood).
		Background(tc.ColorBlack)

	gold := tc.StyleDefault.
		Foreground(tc.ColorGold).
		Background(tc.ColorBlack)

	iron := tc.StyleDefault.
		Foreground(tc.ColorSlateGray).
		Background(tc.ColorBlack)

	/* Items */
	goldCoin := items.Item{
		ID:          "gold_coin",
		Rune:        '$',
		Description: "A pile of filthy lucre.",
		Color:       gold,
	}

	rustySword := items.Item{
		ID:          "rusty_sword",
		Rune:        '/',
		Description: "A rusty sword.",
		Color:       iron,
	}

	/* Cells */

	wall := maps.MapCell{
		Rune:        '#',
		Traversable: false,
		Description: "A rough-hewn stone wall.",
		Color:       gray,
	}

	ground := maps.MapCell{
		Rune:        '.',
		Traversable: true,
		Description: "A hard-packed dirt floor.",
		Color:       burlyWood,
	}

	ground_with_gold := maps.MapCell{
		Rune:        '.',
		Traversable: true,
		Description: "A hard-packed dirt floor.",
		Color:       burlyWood,
		Items: items.Items{
			goldCoin,
		},
	}

	ground_with_sword := maps.MapCell{
		Rune:        '.',
		Traversable: true,
		Description: "A hard-packed dirt floor.",
		Color:       burlyWood,
		Items: items.Items{
			rustySword,
		},
	}

	down_stairs := maps.MapCell{
		Rune:        '>',
		Traversable: true,
		Description: "A maze of twisty stairs leading down.",
		Color:       brown,
	}

	up_stairs := maps.MapCell{
		Rune:        '<',
		Traversable: true,
		Description: "A maze of twisty stairs leading up.",
		Color:       brown,
	}

	g.dungeon = append(g.dungeon, maps.Map{
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
		{wall, ground_with_sword, ground, ground, ground, ground, ground, ground, wall},
		{wall, ground, ground, ground, wall, ground, ground, ground, wall},
		{wall, ground, ground, ground, ground, ground, ground, ground, wall},
		{wall, ground, ground, ground, ground, ground, up_stairs, ground, wall},
		{wall, wall, ground, ground, ground, ground, ground, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, wall},
	})
	g.floor = 0
	current := &g.dungeon[g.floor]

	g.screen.SetStyle(tc.StyleDefault.
		Foreground(tc.ColorWhite).
		Background(tc.ColorBlack))

	g.screen.EnableMouse()
	g.screen.Clear()

	g.player = actors.NewActor(2, 2, 0, white, &player.Player{Rune: '@', Health: 10, Description: "Player"})

	g.creatures = append(
		g.creatures,
		actors.NewActor(2, 6, 0, pink, &cr.Pig{Rune: 'p', Health: 5, Description: "A Pig who loves straw"}),
		actors.NewActor(3, 6, 0, pink, &cr.Pig{Rune: 'p', Health: 5, Description: "A Pig who loves sticks"}),
		actors.NewActor(4, 6, 0, pink, &cr.Pig{Rune: 'p', Health: 5, Description: "A Pig who loves bricks"}),
		actors.NewActor(4, 3, 1, white, &cr.Rat{Rune: 'r', Health: 10, Description: "Lab Rat"}),
		actors.NewActor(4, 5, 1, brown, &cr.Rat{Rune: 'r', Health: 10, Description: "You Dirty Rat"}),
		actors.NewActor(3, 4, 1, white, &cr.Rat{Rune: 'r', Health: 10, Description: "Ratt *Plays Guitar Riff*"}),
		actors.NewActor(5, 4, 1, brown, &cr.Rat{Rune: 'r', Health: 10, Description: "Rat-tatooee"}),
	)

	quit := func() {
		g.screen.Fini()
		os.Exit(0)
	}

	// TODO: review creature movement so it doesn't move before the initial draw.
	// TODO: Adjust key events to only trigger a turn if the one of the movement/action keys are pressed.

	first_pass := true

	for {
		// Update Screen
		g.screen.Show()

		// Poll Event
		ev := g.screen.PollEvent()

		switch ev := ev.(type) {
		case *tc.EventResize:
			g.screen.Sync()

		case *tc.EventKey:
			switch ev.Key() {
			case tc.KeyRune:
				switch ev.Rune() {
				case ':':
					if (*current)[g.player.X][g.player.Y].Items != nil {
						g.msg = (*current)[g.player.X][g.player.Y].GetFirstItem().GetDescription()
					} else {
						g.msg = (*current)[g.player.X][g.player.Y].GetDescription()
					}

				case '>':
					r := (*current)[g.player.X][g.player.Y].GetRune()
					if r == '>' {
						g.floor++
						g.screen.Clear()
					}
				case '<':
					r := (*current)[g.player.X][g.player.Y].GetRune()
					if r == '<' {
						g.floor--
						g.screen.Clear()
					}
				}

			case tc.KeyRight:
				g.player.Move((*current), 1, 0)

			case tc.KeyLeft:
				g.player.Move((*current), -1, 0)

			case tc.KeyUp:
				g.player.Move((*current), 0, -1)

			case tc.KeyDown:
				g.player.Move((*current), 0, 1)

			case tc.KeyCtrlD:
				g.debug = !g.debug

			case tc.KeyCtrlL:
				g.screen.Clear()
				g.screen.Sync()

			case tc.KeyCtrlC, tc.KeyEscape:
				quit()

			}
		default:
			continue

		}

		// Creature Movement
		// First pass stops creature movement before the initial draw
		if !first_pass {
			for _, c := range g.creatures {
				if c.Floor == g.player.Floor {
					c.Wander(*current, rand.Intn(5)+1, g.floor)
				}
			}
		}

		first_pass = false

		// Process Event
		dbg := fmt.Sprintf("player floor: %d x: %d y: %d", g.player.Floor, g.player.X, g.player.Y)

		if g.debug {
			var yy int
			if g.player.Y == 0 {
				_, yy = g.screen.Size()
				yy--
			} else {
				yy = 0
			}

			ut.EmitStr(g.screen, 20, yy, white, dbg)
		}

		if g.floor == 0 {
			current = &g.dungeon[g.floor]
			g.player.Floor = 0
		} else if g.floor == 1 {
			current = &g.dungeon[g.floor]
			g.player.Floor = 1
		}

		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				var color tc.Style
				var map_rune rune

				if (*current)[i][j].Items != nil {
					// If there are items in the MapCell grab the first item's color and rune
					color = (*current)[i][j].GetFirstItem().GetColor()
					map_rune = (*current)[i][j].GetFirstItem().GetRune()
				} else {
					color = (*current)[i][j].GetColor()
					map_rune = (*current)[i][j].GetRune()
				}

				ut.EmitStr(g.screen, i, j, color, string(map_rune))
			}
		}

		for _, c := range g.creatures {
			if c.Floor == g.player.Floor {
				c.Draw(g.screen, g.floor)
			}
		}

		g.player.Draw(g.screen, g.floor)

		if g.msg != "" {
			ut.EmitStr(g.screen, 0, 0, white, g.msg)
		}

	}
}
