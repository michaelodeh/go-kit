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
	js jetstream.JetStream
}

func NewNastJetStreamClient(js jetstream.JetStream) *NastJetStreamClient {
	return &NastJetStreamClient{js: js}
}

// NewConsumer creates a consumer for streamName. The supplied config is a
// template: the application subject passed to Consume becomes FilterSubject.
//
// A consumer is started once. Create another NastStreamConsumer when an
// application needs independent subscriptions or different filters.
func (ns *NastJetStreamClient) NewConsumer(streamName string, opts ...*jetstream.ConsumerConfig) *nats_stream_consumer.NastStreamConsumer {
	if ns == nil {
		return nats_stream_consumer.NewNastStreamConsumer(nil, streamName, opts...)
	}
	return nats_stream_consumer.NewNastStreamConsumer(ns.js, streamName, opts...)
}

func (ns *NastJetStreamClient) NewPublisher() *nats_stream_publisher.NastJetStreamPublisher {
	if ns == nil {
		return nats_stream_publisher.NewNastJetStreamPublisher(nil)
	}
	return nats_stream_publisher.NewNastJetStreamPublisher(ns.js)
}
