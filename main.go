package main

import (
	"log"

	"github.com/andrejsoucek/nats-cui/config"
	"github.com/andrejsoucek/nats-cui/jetstream"
	"github.com/andrejsoucek/nats-cui/ui"
	"github.com/jroimartin/gocui"
)

func main() {
	config.New()
	jetstream.GetJetStream()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	// g.Cursor = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	g.InputEsc = true

	g.SetManagerFunc(ui.CreateLayout)

	if err := ui.BindKeys(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
