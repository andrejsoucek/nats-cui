package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/nats-io/nats.go"
)

var (
	logView *gocui.View
	js      nats.JetStreamContext
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("buckets", 1, 1, int(0.2*float32(maxX)), maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Buckets"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for storeName := range getBuckets() {
			fmt.Fprintln(v, storeName)
		}
		g.SetCurrentView("buckets")
	}
	if v, err := g.SetView("keys", int(0.2*float32(maxX)), 1, int(0.4*float32(maxX)), maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Keys"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	if v, err := g.SetView("main", int(0.4*float32(maxX)), 1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Value"
		v.Wrap = true
	}

	if v, err := g.SetView("log", -1, maxY-5, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Log"
		v.Wrap = true
		v.Autoscroll = true
		logView = v
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	// g.Cursor = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	g.InputEsc = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func selectBucket(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return err
	}
	fmt.Fprintln(logView, "Selected bucket: "+line)

	keys, err := getKeys(line)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}

	keysView, err := g.View("keys")
	if err != nil {
		return err
	}

	keysView.Clear()
	for _, key := range keys {
		fmt.Fprintln(keysView, key)
	}
	g.SetCurrentView("keys")

	return nil
}

func unselectBucket(g *gocui.Gui, v *gocui.View) error {
	logString("Unselected bucket")
	g.SetCurrentView("buckets")

	keysView, err := g.View("keys")
	if err != nil {
		return err
	}

	keysView.Clear()

	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("buckets", gocui.KeyEnter, gocui.ModNone, selectBucket); err != nil {
		return err
	}
	if err := g.SetKeybinding("keys", gocui.KeyEsc, gocui.ModNone, unselectBucket); err != nil {
		return err
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func getBuckets() <-chan string {
	nc, err := nats.Connect("nats://localhost:4222", nats.UserInfo("trust", "super-secret-nats-password-admin"))
	if err != nil {
		log.Fatal(err)
	}
	js, err = nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	return js.KeyValueStoreNames()
}

func getKeys(bucket string) ([]string, error) {
	if js == nil {
		logString("JetStream not initialized")
	}
	kv, err := js.KeyValue(bucket)
	if err != nil {
		return nil, err
	}

	keys, err := kv.Keys()
	if err != nil {
		if err == nats.ErrNoKeysFound {
			logString("No keys found in bucket " + bucket)
			return []string{}, nil
		}
		return nil, err
	}

	return keys, nil
}

func logString(s string) {
	fmt.Fprintln(logView, s)
}
