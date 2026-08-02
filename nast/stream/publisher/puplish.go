package nats_stream_publisher

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

type NastJetStreamPublisher struct {
	js *jetstream.JetStream
}

func NewNastJetStreamPublisher(js *jetstream.JetStream) *NastJetStreamPublisher {
	return &NastJetStreamPublisher{js: js}
}

func (ns *NastJetStreamPublisher) Publish(ctx context.Context, subject string, payload []byte) error {

	ack, err := (*ns.js).Publish(ctx, subject, payload)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	fmt.Println("Stored with sequence:", ack.Sequence)

	return nil
}
