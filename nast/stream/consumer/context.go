package nats_stream_consumer

import (
	"context"
)

func (ns *NastStreamConsumer) Context() context.Context {
	return context.Background()
}
