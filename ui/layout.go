package ui

import (
	"fmt"

	"github.com/andrejsoucek/nats-cui/command"
	"github.com/andrejsoucek/nats-cui/text"
	"github.com/jroimartin/gocui"
)

func CreateLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("buckets", 1, 1, int(0.2*float32(maxX)), maxY-16); err != nil {
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
	if v, err := g.SetView("keys", int(0.2*float32(maxX)), 1, int(0.7*float32(maxX)), maxY-16); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Keys"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	if v, err := g.SetView("bucketInfo", int(0.7*float32(maxX)), 1, maxX-1, 5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Bucket Info"
		v.Wrap = true
	}
	if v, err := g.SetView("keyInfo", int(0.7*float32(maxX)), 5, maxX-1, 15); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Key Info"
		v.Wrap = true
	}
	if v, err := g.SetView("value", int(0.7*float32(maxX)), 15, maxX-1, maxY-16); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Value"
		v.Wrap = true
	}

	if v, err := g.SetView("log", 1, maxY-15, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Log"
		v.Wrap = true
		v.Autoscroll = true
	}

	if v, err := g.SetView("help", 1, maxY-2, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprintln(v, text.Bold("↑↓ - move cursor | Enter - select/refresh value | Esc - unselect | Del - remove selected key | L - browse log | Ctrl+C - quit"))
	}

	return nil
}
