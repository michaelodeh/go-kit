package nats_stream_consumer

import "context"

// Context is kept for source compatibility with the original wrapper. New
// code should pass an application context to ConsumeWithContext or
// ConsumeFuncWithContext instead.
func (c *NastStreamConsumer) Context() context.Context {
	return context.Background()
}
