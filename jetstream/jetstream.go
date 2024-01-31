package jetstream

import (
	"fmt"

	"github.com/andrejsoucek/nats-cui/config"
	"github.com/nats-io/nats.go"
)

var js nats.JetStreamContext
var selectedBucket = ""
var selectedKey = ""

func GetJetStream() nats.JetStreamContext {
	cfg := config.GetConfig()
	if js == nil {
		nc, err := nats.Connect(
			fmt.Sprintf("nats://%s:%d", cfg.Host, cfg.Port),
			nats.UserInfo(cfg.Username, cfg.Password),
		)
		if err != nil {
			panic(err)
		}
		js, err = nc.JetStream()
		if err != nil {
			panic(err)
		}
	}

	return js
}

func GetSelectedBucket() string {
	return selectedBucket
}

func SelectBucket(b string) {
	selectedBucket = b
}

func UnselectBucket() {
	selectedBucket = ""
}

func GetSelectedKey() string {
	return selectedKey
}

func SelectKey(k string) {
	selectedKey = k
}

func UnselectKey() {
	selectedKey = ""
}
