package command

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func Log(g *gocui.Gui, s string) {
	logView, err := g.View("log")
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(logView, s)
}
