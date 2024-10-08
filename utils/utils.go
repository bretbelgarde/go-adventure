package utils

import (
	tc "github.com/gdamore/tcell/v2"
	rw "github.com/mattn/go-runewidth"
)

func EmitStr(s tc.Screen, x, y int, style tc.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := rw.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}
