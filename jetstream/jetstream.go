package jetstream

import (
	"log"

	"github.com/nats-io/nats.go"
)

var js nats.JetStreamContext
var selectedBucket = ""
var selectedKey = ""

func GetJetStream() nats.JetStreamContext {
	if js == nil {
		nc, err := nats.Connect("nats://localhost:4222", nats.UserInfo("trust", "super-secret-nats-password-admin"))
		if err != nil {
			log.Panicln(err)
		}
		js, err = nc.JetStream()
		if err != nil {
			log.Panicln(err)
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
