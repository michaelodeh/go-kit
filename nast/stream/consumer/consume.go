package nats_stream_consumer

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

// Consumer is implemented by application workers that need a context-aware
// message handler. A nil error acknowledges the message; a non-nil error
// causes the wrapper to Nak it for redelivery (unless AckNonePolicy is used).
type Consumer interface {
	Consume(context.Context, jetstream.Msg) error
}

// ConsumerFunc adapts a function to Consumer.
type ConsumerFunc func(context.Context, jetstream.Msg) error

func (f ConsumerFunc) Consume(ctx context.Context, msg jetstream.Msg) error {
	return f(ctx, msg)
}
