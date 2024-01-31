package command

import (
	"github.com/jroimartin/gocui"
)

func CursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		line, err := v.Line(cy + 1)
		if err != nil {
			return err
		}
		if len(line) == 0 {
			return nil
		}
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func CursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		line, err := v.Line(cy)
		if err != nil {
			return err
		}
		if len(line) == 0 {
			return nil
		}
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}
