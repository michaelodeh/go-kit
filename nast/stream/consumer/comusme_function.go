package nats_stream_consumer

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

// ConsumeFunc starts a consumer using context.Background(). Retain the
// returned ConsumeContext and stop or drain it during shutdown.
func (c *NastStreamConsumer) ConsumeFunc(topic string, handler func(jetstream.Msg) error) (jetstream.ConsumeContext, error) {
	return c.ConsumeFuncWithContext(context.Background(), topic, handler)
}

// ConsumeFuncWithContext is the function-oriented equivalent of
// ConsumeWithContext. The handler is acknowledged automatically when it
// returns nil and Nak'ed when it returns an error.
func (c *NastStreamConsumer) ConsumeFuncWithContext(ctx context.Context, topic string, handler func(jetstream.Msg) error) (jetstream.ConsumeContext, error) {
	if handler == nil {
		return nil, ErrHandlerRequired
	}
	return c.start(ctx, topic, handler)
}
