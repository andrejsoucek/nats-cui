package command

import (
	"github.com/jroimartin/gocui"
)

func ScrollView(g *gocui.Gui, v *gocui.View, dy int) {
	if v != nil {
		v.Autoscroll = false
		g.Cursor = true
		ox, oy := v.Origin()

		_, cy := v.Cursor()
		line, _ := v.Line(cy + 1)
		if len(line) == 0 {
			return
		}
		if cy == 0 {
			return
		}

		v.SetOrigin(ox, oy+dy)
	}
}
