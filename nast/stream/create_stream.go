package nats_stream

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

// CreateStream ensures that a stream with cfg exists.
//
// The old implementation treated every Stream lookup error as "not found" and
// printed directly to stdout. CreateOrUpdateStream lets JetStream distinguish
// a missing stream from authentication, connection, and configuration errors.
func (ns *NastJetStreamClient) CreateStream(ctx context.Context, cfg jetstream.StreamConfig) error {
	if ns == nil || ns.js == nil {
		return errors.New("nats stream client is nil")
	}
	if ctx == nil {
		return errors.New("nats stream context is nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := ns.js.CreateOrUpdateStream(ctx, cfg); err != nil {
		return fmt.Errorf("ensure stream %q: %w", cfg.Name, err)
	}

	return nil
}
