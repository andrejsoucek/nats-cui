package command

import (
	"errors"
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

	err = RefreshKeys(g)
	if err != nil {
		if errors.Is(err, nats.ErrNoKeysFound) {
			Log(g, fmt.Sprintf("Bucket %s is empty", jetstream.GetSelectedBucket()))
			return nil
		}
		return err
	}
	g.SetCurrentView("keys")

	return nil
}

func RefreshKeys(g *gocui.Gui) error {
	keysView, err := g.View("keys")
	if err != nil {
		return err
	}
	keysView.Clear()

	keys, err := getKeys(g, jetstream.GetSelectedBucket())
	if err != nil {
		return err
	}

	for _, key := range keys {
		fmt.Fprintln(keysView, key)
	}

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
		return nil, err
	}

	return keys, nil
}
