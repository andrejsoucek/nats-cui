package ui

import (
	"github.com/andrejsoucek/nats-cui/command"
	"github.com/jroimartin/gocui"
)

var previousView string

func BindKeys(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, command.SelectDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, command.SelectUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, command.Quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'l', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		previousView = g.CurrentView().Name()
		g.SetCurrentView("log")
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("buckets", gocui.KeyEnter, gocui.ModNone, command.SelectBucket); err != nil {
		return err
	}

	if err := g.SetKeybinding("keys", gocui.KeyEsc, gocui.ModNone, command.UnselectBucket); err != nil {
		return err
	}
	if err := g.SetKeybinding("keys", gocui.KeyEnter, gocui.ModNone, command.SelectKey); err != nil {
		return err
	}

	if err := g.SetKeybinding("value", gocui.KeyEsc, gocui.ModNone, command.UnselectKey); err != nil {
		return err
	}
	if err := g.SetKeybinding("value", gocui.KeyEnter, gocui.ModNone, command.SelectKey); err != nil {
		return err
	}
	if err := g.SetKeybinding("value", gocui.KeyDelete, gocui.ModNone, command.ConfirmDelete); err != nil {
		return err
	}

	if err := g.SetKeybinding("confirm", gocui.KeyEnter, gocui.ModNone, command.DeleteKey); err != nil {
		return err
	}
	if err := g.SetKeybinding("confirm", gocui.KeyEsc, gocui.ModNone, command.HideConfirmationDialog); err != nil {
		return err
	}

	if err := g.SetKeybinding("log", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			command.ScrollView(g, v, -1)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("log", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			command.ScrollView(g, v, 1)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("log", gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		v.Autoscroll = true
		g.SetCurrentView(previousView)
		return nil
	}); err != nil {
		return err
	}

	return nil
}
