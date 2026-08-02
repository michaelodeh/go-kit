package nats_stream_consumer

import (
	"github.com/nats-io/nats.go/jetstream"
)

type NatsStreamOption struct {
	jetstream.ConsumerConfig

	MaxFetch int
}
