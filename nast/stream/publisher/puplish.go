package nats_stream_publisher

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

type NastJetStreamPublisher struct {
	js jetstream.JetStream
}

func NewNastJetStreamPublisher(js jetstream.JetStream) *NastJetStreamPublisher {
	return &NastJetStreamPublisher{js: js}
}

func (ns *NastJetStreamPublisher) Publish(ctx context.Context, subject string, payload []byte) error {
	if ns == nil || ns.js == nil {
		return errors.New("nats publisher is nil")
	}
	if ctx == nil {
		return errors.New("nats publisher context is nil")
	}

	_, err := ns.js.Publish(ctx, subject, payload)
	if err != nil {
		return fmt.Errorf("publish %q: %w", subject, err)
	}

	return nil
}
