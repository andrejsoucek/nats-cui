package command

import (
	"errors"
	"fmt"
	"time"

	"github.com/andrejsoucek/nats-cui/jetstream"
	"github.com/jroimartin/gocui"
	"github.com/nats-io/nats.go"
)

func SelectKey(g *gocui.Gui, v *gocui.View) error {
	if g.CurrentView().Name() != "value" {
		_, cy := v.Cursor()
		line, err := v.Line(cy)
		if err != nil {
			return err
		}
		jetstream.SelectKey(line)
		Log(g, "Selected key: "+line)
	} else {
		Log(g, "Refreshing value")
	}

	status, value, err := getValue(g, jetstream.GetSelectedKey())
	if err != nil {
		if errors.Is(err, nats.ErrKeyNotFound) {
			Log(g, fmt.Sprintf("Key %s not found", jetstream.GetSelectedKey()))
			UnselectKey(g, v)
			err := RefreshKeys(g)
			if err != nil {
				if errors.Is(err, nats.ErrNoKeysFound) {
					Log(g, fmt.Sprintf("Bucket %s is empty", jetstream.GetSelectedBucket()))
					g.SetCurrentView("buckets")
					return nil
				}
				return err
			}
			return nil
		}
		return err
	}

	valueView, err := g.View("value")
	if err != nil {
		return err
	}

	valueView.Clear()

	fmt.Fprintln(valueView, "────────────── BUCKET INFORMATION ─────────────────────")
	fmt.Fprintln(valueView, bold("Bucket TTL:"))
	fmt.Fprintln(valueView, status.TTL())

	fmt.Fprintln(valueView, "──────────────── KEY INFORMATION ──────────────────────")
	fmt.Fprintln(valueView, bold("Created:"))
	fmt.Fprintln(valueView, value.Created())
	fmt.Fprintln(valueView, bold("Expires:"))
	fmt.Fprintln(valueView, value.Created().Add(status.TTL()))
	fmt.Fprintln(valueView, bold("TTL:"))
	fmt.Fprintln(valueView, value.Created().Add(status.TTL()).Sub(time.Now()))
	fmt.Fprintln(valueView, "──────────────────── VALUE ────────────────────────────")
	fmt.Fprintln(valueView, bold("Value:"))
	fmt.Fprintln(valueView, string(value.Value()))
	fmt.Fprintln(valueView, bold("Raw Value:"))
	fmt.Fprintln(valueView, value.Value())

	g.SetCurrentView("value")

	return nil
}

func UnselectKey(g *gocui.Gui, v *gocui.View) error {
	jetstream.UnselectKey()
	g.SetCurrentView("keys")

	valueView, err := g.View("value")
	if err != nil {
		return err
	}

	valueView.Clear()

	return nil
}

func ConfirmDelete(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("confirm", maxX/2-150, maxY/2-10, maxX/2+150, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		fmt.Fprintln(v, bold("Are you sure you want to delete the key?"))
		fmt.Fprintln(v, jetstream.GetSelectedKey())
		if _, err := g.SetCurrentView("confirm"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("confirmHelp", maxX/2-150, maxY/2, maxX/2+150, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprintln(v, "Enter - confirm | Esc - cancel")
		if _, err := g.SetCurrentView("confirm"); err != nil {
			return err
		}
	}

	return nil
}

func HideConfirmationDialog(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView("confirm")
	g.DeleteView("confirmHelp")
	g.SetCurrentView("value")

	return nil
}

func DeleteKey(g *gocui.Gui, v *gocui.View) error {
	HideConfirmationDialog(g, v)
	js := jetstream.GetJetStream()
	kv, err := js.KeyValue(jetstream.GetSelectedBucket())
	if err != nil {
		return err
	}
	kv.Delete(jetstream.GetSelectedKey())
	Log(g, "Deleted key: "+jetstream.GetSelectedKey())
	UnselectKey(g, v)
	err = RefreshKeys(g)
	if err != nil {
		if errors.Is(err, nats.ErrNoKeysFound) {
			Log(g, fmt.Sprintf("Bucket %s is empty", jetstream.GetSelectedBucket()))
			g.SetCurrentView("buckets")
			return nil
		}
		return err
	}
	g.SetCurrentView("keys")

	return nil
}

func getValue(g *gocui.Gui, key string) (nats.KeyValueStatus, nats.KeyValueEntry, error) {
	js := jetstream.GetJetStream()
	kv, err := js.KeyValue(jetstream.GetSelectedBucket())
	if err != nil {
		return nil, nil, err
	}

	status, err := kv.Status()
	if err != nil {
		return nil, nil, err
	}

	kve, err := kv.Get(key)
	if err != nil {
		if err != nil {
			return nil, nil, err
		}
		return nil, nil, err
	}

	return status, kve, nil
}

func bold(s string) string {
	return fmt.Sprintf("%s%s%s", "\033[1m", s, "\033[0m")
}
