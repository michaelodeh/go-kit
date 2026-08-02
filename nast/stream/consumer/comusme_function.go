package nats_stream_consumer

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

func (ns *NastStreamConsumer) ConsumeFunc(ctx context.Context, topic string, handlerFunc func(data jetstream.Msg) error) {

	println("consume function called")
	go func() {
		if err := ns.rawConsumer(ctx, topic, func(msg jetstream.Msg) {
			if err := handlerFunc(msg); err != nil {
				println(err.Error())
				msg.Nak()
				return
			}
			msg.Ack()
		}); err != nil {
			println("rawConsumer:", err.Error())
		}
	}()

}
