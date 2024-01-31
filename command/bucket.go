package command

import (
	"fmt"

	"github.com/andrejsoucek/nats-cui/jetstream"
	"github.com/jroimartin/gocui"
	"github.com/nats-io/nats.go"
)

func GetBuckets() <-chan string {
	js := jetstream.GetJetStream()
	return js.KeyValueStoreNames()
}

func SelectBucket(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return err
	}
	jetstream.SelectBucket(line)
	Log(g, "Selected bucket: "+line)

	keys, err := getKeys(g, jetstream.GetSelectedBucket())
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

func UnselectBucket(g *gocui.Gui, v *gocui.View) error {
	jetstream.UnselectBucket()
	Log(g, "Unselected bucket")
	g.SetCurrentView("buckets")

	keysView, err := g.View("keys")
	if err != nil {
		return err
	}

	keysView.Clear()

	return nil
}

func getKeys(g *gocui.Gui, bucket string) ([]string, error) {
	js := jetstream.GetJetStream()
	kv, err := js.KeyValue(bucket)
	if err != nil {
		return nil, err
	}

	keys, err := kv.Keys()
	if err != nil {
		if err == nats.ErrNoKeysFound {
			Log(g, "No keys found in bucket "+bucket)
			return []string{}, nil
		}
		return nil, err
	}

	return keys, nil
}
