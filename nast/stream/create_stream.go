package nats_stream

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

func (ns *NastJetStreamClient) CreateStream(ctx context.Context, cfg jetstream.StreamConfig) error {

	fmt.Println("Stream name: ", cfg.Name)
	fmt.Println("Stream config: ", cfg.Subjects)

	if _, err := (*ns.js).Stream(ctx, cfg.Name); err == nil {
		return nil // already exists
	}

	_, err := (*ns.js).CreateStream(ctx, cfg)

	if err != nil {
		return fmt.Errorf("failed to create stream: %v", err)
	}

	fmt.Println("Stream created successfully")

	return nil
}
