package nats_stream_consumer

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

type NastStreamConsumer struct {
	js         *jetstream.JetStream
	streamName string
	option     *jetstream.ConsumerConfig
}

func NewNastStreamConsumer(js *jetstream.JetStream, streamName string, option ...*jetstream.ConsumerConfig) *NastStreamConsumer {
	cfg := parseOptions(option...)
	return &NastStreamConsumer{js: js, streamName: streamName, option: cfg}
}

func (c *NastStreamConsumer) rawConsumer(
	ctx context.Context,
	topic string,
	handler func(jetstream.Msg),
) error {
	print("Raw consume function")
	consumer, err := (*c.js).CreateOrUpdateConsumer(
		ctx,
		c.streamName,
		jetstream.ConsumerConfig{
			Durable:       c.option.Durable,
			AckPolicy:     c.option.AckPolicy,
			DeliverGroup:  c.option.DeliverGroup,
			DeliverPolicy: c.option.DeliverPolicy,
			FilterSubject: topic,
			// DeliverSubject: topic,
			ReplayPolicy:  c.option.ReplayPolicy,
			AckWait:       c.option.AckWait,
			MaxDeliver:    c.option.MaxDeliver,
			MaxAckPending: c.option.MaxAckPending,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create or update consumer: %v", err)
	}

	_, err = consumer.Consume(handler)
	if err != nil {
		return fmt.Errorf("failed to consume: %v", err)
	}
	return nil
}
