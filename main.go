package main

import (
	"fmt"
	"os"

	"bretbelgarde.com/adventure/actors"
	cr "bretbelgarde.com/adventure/creatures"
	"bretbelgarde.com/adventure/items"
	"bretbelgarde.com/adventure/maps"
	"bretbelgarde.com/adventure/player"
	ui "bretbelgarde.com/adventure/ui"
	ut "bretbelgarde.com/adventure/utils"
	tc "github.com/gdamore/tcell/v2"
)

type Game struct {
	screen  tc.Screen
	debug   bool
	msg     string
	floor   int
	dungeon maps.Floors
	actors  actors.Actors
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

func (g *Game) HandleMovement(current *maps.Map, x, y int) {
	// There may be a more graceful way to do this, but this will do for now
	p := g.actors.GetActorFromID("player")
	origx, origy, _ := p.GetLocation()
	p.Move(*current, x, y)
	newx, newy, _ := p.GetLocation()

	if newx != origx || newy != origy {
		for _, c := range g.actors {
			if c.Floor == p.Floor && c.ID != "player" {
				c.Wander(*current, g.floor)
			}
		}
	}
}

func main() {
	var g Game
	g.Init()

	var (
		/* Colors */
		white = tc.StyleDefault.
			Foreground(tc.ColorWhite).
			Background(tc.ColorBlack)

		brown = tc.StyleDefault.
			Foreground(tc.ColorBrown).
			Background(tc.ColorBlack)

		pink = tc.StyleDefault.
			Foreground(tc.ColorPink).
			Background(tc.ColorBlack)

		gray = tc.StyleDefault.
			Foreground(tc.ColorGray).
			Background(tc.ColorBlack)

		burlyWood = tc.StyleDefault.
				Foreground(tc.ColorBurlyWood).
				Background(tc.ColorBlack)

		gold = tc.StyleDefault.
			Foreground(tc.ColorGold).
			Background(tc.ColorBlack)

		iron = tc.StyleDefault.
			Foreground(tc.ColorSlateGray).
			Background(tc.ColorBlack)

		/* Items */
		goldCoin = items.Item{
			ID:          "gold_coin",
			Rune:        '$',
			Description: "A pile of filthy lucre.",
			Color:       gold,
		}

		rustySword = items.Item{
			ID:          "rusty_sword",
			Rune:        '/',
			Description: "A rusty sword.",
			Color:       iron,
		}

		/* Cells */

		wall = maps.MapCell{
			Rune:        '#',
			Traversable: false,
			Description: "A rough-hewn stone wall.",
			Color:       gray,
		}

		ground = maps.MapCell{
			Rune:        '.',
			Traversable: true,
			Description: "A hard-packed dirt floor.",
			Color:       burlyWood,
		}

		ground_with_gold = maps.MapCell{
			Rune:        '.',
			Traversable: true,
			Description: "A hard-packed dirt floor.",
			Color:       burlyWood,
			Items: items.Items{
				goldCoin,
				rustySword,
			},
		}

		ground_with_sword = maps.MapCell{
			Rune:        '.',
			Traversable: true,
			Description: "A hard-packed dirt floor.",
			Color:       burlyWood,
			Items: items.Items{
				rustySword,
			},
		}

		down_stairs = maps.MapCell{
			Rune:        '>',
			Traversable: true,
			Description: "A maze of twisty stairs leading down.",
			Color:       brown,
		}

		up_stairs = maps.MapCell{
			Rune:        '<',
			Traversable: true,
			Description: "A maze of twisty stairs leading up.",
			Color:       brown,
		}
	)

	g.dungeon = append(g.dungeon, maps.Map{
		{wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, ground, ground, ground, ground, ground, ground, ground, wall},
		{wall, ground, ground, ground, ground, ground, down_stairs, ground, wall},
		{wall, ground, ground, ground, wall, ground, ground, ground, wall},
		{wall, ground, ground, wall, wall, wall, ground, ground, wall},
		{wall, ground, ground_with_gold, ground, wall, ground, ground, ground, wall},
		{wall, ground, ground, ground, ground, ground, ground, ground, wall},
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

	g.screen.Clear()

	g.actors = append(
		g.actors,
		actors.NewActor("player", 2, 2, 0, white, &player.Player{Rune: '@', Health: 10, Description: "Player"}),
		actors.NewActor("pig_1", 2, 6, 0, pink, &cr.Pig{Rune: 'p', Health: 5, Description: "A Pig who loves straw"}),
		actors.NewActor("pig_2", 3, 6, 0, pink, &cr.Pig{Rune: 'p', Health: 5, Description: "A Pig who loves sticks"}),
		actors.NewActor("pig_3", 4, 6, 0, pink, &cr.Pig{Rune: 'p', Health: 5, Description: "A Pig who loves bricks"}),
		actors.NewActor("rat_1", 4, 3, 1, white, &cr.Rat{Rune: 'r', Health: 10, Description: "Lab Rat"}),
		actors.NewActor("rat_2", 4, 5, 1, brown, &cr.Rat{Rune: 'r', Health: 10, Description: "You Dirty Rat"}),
		actors.NewActor("rat_3", 3, 4, 1, white, &cr.Rat{Rune: 'r', Health: 10, Description: "Ratt *Plays Guitar Riff*"}),
		actors.NewActor("rat_4", 5, 4, 1, brown, &cr.Rat{Rune: 'r', Health: 10, Description: "Rat-tatooee"}),
	)

	p := g.actors.GetActorFromID("player")

	cursor := ui.NewCursor(p.X, p.Y, g.floor, &g.dungeon[g.floor])

	quit := func() {
		g.screen.Fini()
		os.Exit(0)
	}

	for {
		// Update Screen
		g.screen.Clear()

		g.msg = ""

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
					cursor.IsActive = !cursor.IsActive

					if !cursor.IsActive {
						cursor.X = p.X
						cursor.Y = p.Y
						cursor.SetCurrentFloor(p.Floor, g.dungeon[p.Floor])
					}

				case '>':
					r := g.dungeon[g.floor].GetCell(p.X, p.Y).GetRune()

					if r == '>' {
						g.floor++
					}
				case '<':
					r := g.dungeon[g.floor].GetCell(p.X, p.Y).GetRune()

					if r == '<' {
						g.floor--
					}
				}

			case tc.KeyRight:
				if !cursor.IsActive {
					g.HandleMovement(current, 1, 0)
				}

				cursor.Move(1, 0)

			case tc.KeyLeft:
				if !cursor.IsActive {
					g.HandleMovement(current, -1, 0)
				}

				cursor.Move(-1, 0)

			case tc.KeyUp:
				if !cursor.IsActive {
					g.HandleMovement(current, 0, -1)
				}

				cursor.Move(0, -1)

			case tc.KeyDown:
				if !cursor.IsActive {
					g.HandleMovement(current, 0, 1)
				}

				cursor.Move(0, 1)

			case tc.KeyCtrlD:
				g.debug = !g.debug

			case tc.KeyCtrlL:
				g.screen.Clear()
				g.screen.Sync()

			case tc.KeyCtrlQ:
				quit()

			}
		default:
			continue

		}

		if cursor.IsActive {
			g.msg = cursor.Look(&g.actors)
		}

		// Process Event
		dbg := fmt.Sprintf("player floor: %d x: %d y: %d", p.Floor, p.X, p.Y)

		if g.debug {
			ut.EmitStr(g.screen, 20, 1, white, dbg)
		}

		if g.floor == 0 {
			cursor.SetCurrentFloor(p.Floor, g.dungeon[g.floor])
			p.Floor = 0
		} else if g.floor == 1 {
			cursor.SetCurrentFloor(p.Floor, g.dungeon[g.floor])
			p.Floor = 1
		}

		for row := 0; row < g.dungeon[g.floor].GetWidth(); row++ {
			for col := 0; col < g.dungeon[g.floor].GetHeight(); col++ {

				var color tc.Style
				var map_rune rune

				if (g.dungeon[g.floor]).GetCell(col, row).Items != nil {
					// If there are items in the MapCell grab the first item's color and rune
					color = (g.dungeon[g.floor]).GetCell(col, row).GetFirstItem().GetColor()
					map_rune = (g.dungeon[g.floor]).GetCell(col, row).GetFirstItem().GetRune()
				} else {
					color = (g.dungeon[g.floor]).GetCell(col, row).GetColor()
					map_rune = (g.dungeon[g.floor]).GetCell(col, row).GetRune()
				}

				ut.EmitStr(g.screen, col, row, color, string(map_rune))
			}
		}

		for _, c := range g.actors {
			if c.ID != "player" && c.Floor == p.Floor {
				c.Draw(g.screen, g.floor)
			}
		}

		p.Draw(g.screen, g.floor)

		if cursor.IsActive {
			cursor.Draw(g.screen)
		}

		if g.msg != "" {
			ut.EmitStr(g.screen, 20, 0, white, g.msg)
		}

		g.screen.Show()
	}
}
