package nats_stream

import (
	nats_stream_consumer "github.com/michaelodeh/go-kit/nast/stream/consumer"
	nats_stream_publisher "github.com/michaelodeh/go-kit/nast/stream/publisher"
	"github.com/nats-io/nats.go/jetstream"
)

// type publisher interface {
// 	New() *nats_stream.NastJetStreamPublish
// }

type NastJetStreamClient struct {
	js *jetstream.JetStream
}

func NewNastJetStreamClient(js *jetstream.JetStream) *NastJetStreamClient {
	return &NastJetStreamClient{js: js}
}

func (ns *NastJetStreamClient) NewConsumer(streamName string, opts ...*jetstream.ConsumerConfig) *nats_stream_consumer.NastStreamConsumer {
	return nats_stream_consumer.NewNastStreamConsumer(ns.js, streamName, opts...)
}

func (ns *NastJetStreamClient) NewPublisher() *nats_stream_publisher.NastJetStreamPublisher {
	return nats_stream_publisher.NewNastJetStreamPublisher(ns.js)
}
