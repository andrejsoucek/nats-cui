package command

import (
	"errors"
	"fmt"
	"time"

	"github.com/andrejsoucek/nats-cui/jetstream"
	"github.com/andrejsoucek/nats-cui/text"
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

	biv, kiv, vv, err := clearViews(g)
	if err != nil {
		return err
	}

	fmt.Fprintln(biv, text.Bold("Bucket TTL:"))
	fmt.Fprintln(biv, status.TTL())

	fmt.Fprintln(kiv, text.Bold("Created:"))
	fmt.Fprintln(kiv, value.Created())
	fmt.Fprintln(kiv, text.Bold("Expires:"))
	fmt.Fprintln(kiv, value.Created().Add(status.TTL()))
	fmt.Fprintln(kiv, text.Bold("TTL:"))
	fmt.Fprintln(kiv, value.Created().Add(status.TTL()).Sub(time.Now()))

	fmt.Fprintln(vv, text.Bold("Value:"))
	fmt.Fprintln(vv, string(value.Value()))
	fmt.Fprintln(vv, text.Bold("Raw Value:"))
	fmt.Fprintln(vv, fmt.Sprintf("%#v", value.Value()))

	g.SetCurrentView("value")

	return nil
}

func UnselectKey(g *gocui.Gui, v *gocui.View) error {
	jetstream.UnselectKey()
	g.SetCurrentView("keys")
	clearViews(g)

	return nil
}

func ConfirmDelete(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("confirm", maxX/2-150, maxY/2-10, maxX/2+150, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		fmt.Fprintln(v, text.Bold("Are you sure you want to delete the key?"))
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
		fmt.Fprintln(v, text.Green(text.Bold("Enter - confirm"))+text.Bold(text.Default(" | "))+text.Bold(text.Red("Esc - cancel")))
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

func clearViews(g *gocui.Gui) (*gocui.View, *gocui.View, *gocui.View, error) {

	biv, err := g.View("bucketInfo")
	if err != nil {
		return nil, nil, nil, err
	}

	kiv, err := g.View("keyInfo")
	if err != nil {
		return nil, nil, nil, err
	}

	vv, err := g.View("value")
	if err != nil {
		return nil, nil, nil, err
	}

	biv.Clear()
	kiv.Clear()
	vv.Clear()

	return biv, kiv, vv, nil
}
