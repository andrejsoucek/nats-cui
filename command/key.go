package command

import (
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
		return err
	}

	valueView, err := g.View("value")
	if err != nil {
		return err
	}

	valueView.Clear()

	fmt.Fprintln(valueView, bold("Bucket TTL:"))
	fmt.Fprintln(valueView, status.TTL())
	fmt.Fprintln(valueView, "─────────────────────────────────────────────")

	fmt.Fprintln(valueView, bold("Created:"))
	fmt.Fprintln(valueView, value.Created())
	fmt.Fprintln(valueView, bold("Expires:"))
	fmt.Fprintln(valueView, value.Created().Add(status.TTL()))
	fmt.Fprintln(valueView, bold("TTL:"))
	fmt.Fprintln(valueView, value.Created().Add(status.TTL()).Sub(time.Now()))
	fmt.Fprintln(valueView, "─────────────────────────────────────────────")

	fmt.Fprintln(valueView, bold("Value:"))
	fmt.Fprintln(valueView, fmt.Sprint(value.Value()))

	g.SetCurrentView("value")

	return nil
}

func UnselectKey(g *gocui.Gui, v *gocui.View) error {
	Log(g, "Unselected key")
	jetstream.UnselectKey()
	g.SetCurrentView("keys")

	valueView, err := g.View("value")
	if err != nil {
		return err
	}

	valueView.Clear()

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
		return nil, nil, err
	}

	return status, kve, nil
}

func bold(s string) string {
	return fmt.Sprintf("%s%s%s", "\033[1m", s, "\033[0m")
}
