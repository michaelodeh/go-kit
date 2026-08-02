package nats_stream_consumer

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

type Consumer interface {
	Consume(ctx context.Context, data jetstream.Msg) error
}

func (ns *NastStreamConsumer) Consume(topic string, consumer Consumer) {
	ctx := ns.Context()
	fmt.Println("topic: ", topic)
	fmt.Println("consumer: ", consumer.Consume)
	ns.ConsumeFunc(ctx, topic, func(data jetstream.Msg) error {

		return consumer.Consume(ctx, data)
	})
}
