package main

import (
	tv "github.com/rivo/tview"
)

func main() {
	app := tv.NewApplication()

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
	}
}
