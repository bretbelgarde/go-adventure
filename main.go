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
)

type Cursor struct {
	X        int
	Y        int
	Rune     rune
	Color    tc.Style
	IsActive bool
	Current  *maps.Map
}

func NewCursor(x, y int, m *maps.Map) *Cursor {
	return &Cursor{
		X:        x,
		Y:        y,
		Rune:     tc.RuneBlock,
		Color:    tc.StyleDefault.Foreground(tc.ColorDarkGoldenrod).Background(tc.ColorBlack),
		Current:  m,
		IsActive: false,
	}
}

func (c *Cursor) Draw(s tc.Screen) {
	ut.EmitStr(s, c.X, c.Y, c.Color, string(c.Rune))
}

func (c *Cursor) SetCurrentFloor(floor maps.Map) {
	c.Current = &floor
}

func (c *Cursor) Look(a *actors.Actors) string {
	var seen string
	var selected *maps.MapCell
	var items *items.Items
	var actor_desc string
	var item_desc string

	actor := a.GetActorFromLocaion(c.X, c.Y)
	actor_desc = ""
	item_desc = ""

	if c.X < 0 || c.Y < 0 || c.X >= c.Current.GetWidth() || c.Y >= c.Current.GetHeight() {
		return "You don't see anything."
	} else {
		selected = c.Current.GetCell(c.X, c.Y)
		items = selected.GetItems()

		if actor != nil {
			actor_desc = actor.Type.GetDescription()
		}

		if len(*items) > 1 {
			item_desc = "There is a stack of items on the ground here."
		} else if len(*items) == 1 {
			item_desc = selected.GetFirstItem().GetDescription()
		}

		seen = fmt.Sprintf("%s %s %s", actor_desc, item_desc, selected.GetDescription())

	}

	return seen
}

func (c *Cursor) Move(x, y int) {
	w := c.Current.GetWidth()
	h := c.Current.GetHeight()

	if c.X+x < 0 || c.Y+y < 0 || c.X+x >= w || c.Y+y >= h {
		return
	}

	c.X += x
	c.Y += y
}

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

func (g *Game) HandleMovement(current *maps.Map, x, y int) {
	// There may be a more graceful way to do this, but this will do for now
	origx, origy, _ := g.player.GetLocation()
	g.player.Move(*current, x, y)
	newx, newy, _ := g.player.GetLocation()

	if newx != origx || newy != origy {
		for _, c := range g.creatures {
			if c.Floor == g.player.Floor {
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

	cursor := NewCursor(g.player.X, g.player.Y, &g.dungeon[g.floor])

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
						cursor.X = g.player.X
						cursor.Y = g.player.Y
					}

				case '>':
					r := (*current).GetCell(g.player.X, g.player.Y).GetRune()
					if r == '>' {
						g.floor++
						g.screen.Clear()
					}
				case '<':
					r := (*current).GetCell(g.player.X, g.player.Y).GetRune()
					if r == '<' {
						g.floor--
						g.screen.Clear()
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
			g.msg = cursor.Look(&g.creatures)
		}

		// Process Event
		dbg := fmt.Sprintf("player floor: %d x: %d y: %d", g.player.Floor, g.player.X, g.player.Y)

		if g.debug {
			ut.EmitStr(g.screen, 20, 1, white, dbg)
		}

		if g.floor == 0 {
			current = &g.dungeon[g.floor]
			cursor.SetCurrentFloor(*current)
			g.player.Floor = 0
		} else if g.floor == 1 {
			current = &g.dungeon[g.floor]
			cursor.SetCurrentFloor(*current)
			g.player.Floor = 1
		}

		for row := 0; row < 9; row++ {
			for col := 0; col < 9; col++ {

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

		for _, c := range g.creatures {
			if c.Floor == g.player.Floor {
				c.Draw(g.screen, g.floor)
			}
		}

		g.player.Draw(g.screen, g.floor)

		if cursor.IsActive {
			cursor.Draw(g.screen)
		}

		if g.msg != "" {
			ut.EmitStr(g.screen, 20, 0, white, g.msg)
		}

		g.screen.Show()
	}
}
