package natsultil

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// NatsClient connects to NATS and creates a JetStream client.
//
// jetstream.JetStream is already an interface. Returning it as a pointer to an
// interface makes the client harder to use and provides no additional safety,
// so the library deliberately returns the interface directly.
func NatsClient(url string) (*nats.Conn, jetstream.JetStream, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, nil, err
	}

	return nc, js, nil
}
