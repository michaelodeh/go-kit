package nats_stream_consumer

import (
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

func parseOptions(opts ...*jetstream.ConsumerConfig) *jetstream.ConsumerConfig {
	option := &jetstream.ConsumerConfig{
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
		ReplayPolicy:  jetstream.ReplayInstantPolicy,
		AckWait:       5 * time.Second,
		MaxDeliver:    10,
		MaxAckPending: 100,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		cfg := *opt

		if cfg.Durable != "" {
			option.Durable = cfg.Durable
		}

		if cfg.DeliverGroup != "" {
			option.DeliverGroup = cfg.DeliverGroup
		}

		if cfg.DeliverSubject != "" {
			option.DeliverSubject = cfg.DeliverSubject
		}

		if cfg.FilterSubject != "" {
			option.FilterSubject = cfg.FilterSubject
		}

		if cfg.AckPolicy != 0 {
			option.AckPolicy = cfg.AckPolicy
		}

		if cfg.DeliverPolicy != 0 {
			option.DeliverPolicy = cfg.DeliverPolicy
		}

		if cfg.ReplayPolicy != 0 {
			option.ReplayPolicy = cfg.ReplayPolicy
		}

		if cfg.AckWait > 0 {
			option.AckWait = cfg.AckWait
		}

		if cfg.MaxDeliver > 0 {
			option.MaxDeliver = cfg.MaxDeliver
		}

		if cfg.MaxAckPending > 0 {
			option.MaxAckPending = cfg.MaxAckPending
		}

	}

	return option
}
