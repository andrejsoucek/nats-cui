package jetstream

import (
	"log"

	"github.com/nats-io/nats.go"
)

var js nats.JetStreamContext

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
