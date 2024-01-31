package ui

import (
	"fmt"

	"github.com/andrejsoucek/nats-cui/command"
	"github.com/jroimartin/gocui"
)

func CreateLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("buckets", 1, 1, int(0.2*float32(maxX)), maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Buckets"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for storeName := range command.GetBuckets() { // TODO move
			fmt.Fprintln(v, storeName)
		}
		g.SetCurrentView("buckets")
	}
	if v, err := g.SetView("keys", int(0.2*float32(maxX)), 1, int(0.7*float32(maxX)), maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Keys"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	if v, err := g.SetView("value", int(0.7*float32(maxX)), 1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Value"
		v.Wrap = true
	}

	if v, err := g.SetView("log", 1, maxY-15, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Log"
		v.Wrap = true
		v.Autoscroll = true
	}

	return nil
}
