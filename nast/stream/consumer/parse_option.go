package nats_stream_consumer

import (
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

func parseOptions(options ...*jetstream.ConsumerConfig) jetstream.ConsumerConfig {
	// DeliverAll is important for an application-facing wrapper: messages
	// published before startup are still available. Applications that only want
	// future messages can explicitly set DeliverNewPolicy.
	cfg := jetstream.ConsumerConfig{
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		ReplayPolicy:  jetstream.ReplayInstantPolicy,
		AckWait:       30 * time.Second,
		MaxDeliver:    10,
		MaxAckPending: 1000,
	}

	for _, option := range options {
		if option != nil {
			mergeConsumerConfig(&cfg, *option)
		}
	}

	return cfg
}

func mergeConsumerConfig(dst *jetstream.ConsumerConfig, src jetstream.ConsumerConfig) {
	if src.Name != "" {
		dst.Name = src.Name
	}
	if src.Durable != "" {
		dst.Durable = src.Durable
	}
	if src.Description != "" {
		dst.Description = src.Description
	}
	if src.DeliverPolicy != jetstream.DeliverAllPolicy {
		dst.DeliverPolicy = src.DeliverPolicy
	}
	if src.OptStartSeq != 0 {
		dst.OptStartSeq = src.OptStartSeq
	}
	if src.OptStartTime != nil {
		t := *src.OptStartTime
		dst.OptStartTime = &t
	}
	if src.AckPolicy != jetstream.AckExplicitPolicy {
		dst.AckPolicy = src.AckPolicy
	}
	if src.AckWait != 0 {
		dst.AckWait = src.AckWait
	}
	if src.MaxDeliver != 0 {
		dst.MaxDeliver = src.MaxDeliver
	}
	if src.BackOff != nil {
		dst.BackOff = append([]time.Duration(nil), src.BackOff...)
	}
	if src.FilterSubject != "" {
		dst.FilterSubject = src.FilterSubject
	}
	if src.ReplayPolicy != jetstream.ReplayInstantPolicy {
		dst.ReplayPolicy = src.ReplayPolicy
	}
	if src.RateLimit != 0 {
		dst.RateLimit = src.RateLimit
	}
	if src.SampleFrequency != "" {
		dst.SampleFrequency = src.SampleFrequency
	}
	if src.MaxWaiting != 0 {
		dst.MaxWaiting = src.MaxWaiting
	}
	if src.MaxAckPending != 0 {
		dst.MaxAckPending = src.MaxAckPending
	}
	if src.HeadersOnly {
		dst.HeadersOnly = true
	}
	if src.MaxRequestBatch != 0 {
		dst.MaxRequestBatch = src.MaxRequestBatch
	}
	if src.MaxRequestExpires != 0 {
		dst.MaxRequestExpires = src.MaxRequestExpires
	}
	if src.MaxRequestMaxBytes != 0 {
		dst.MaxRequestMaxBytes = src.MaxRequestMaxBytes
	}
	if src.InactiveThreshold != 0 {
		dst.InactiveThreshold = src.InactiveThreshold
	}
	if src.Replicas != 0 {
		dst.Replicas = src.Replicas
	}
	if src.MemoryStorage {
		dst.MemoryStorage = true
	}
	if src.FilterSubjects != nil {
		dst.FilterSubjects = append([]string(nil), src.FilterSubjects...)
	}
	if src.Metadata != nil {
		dst.Metadata = cloneStringMap(src.Metadata)
	}
	if src.PauseUntil != nil {
		t := *src.PauseUntil
		dst.PauseUntil = &t
	}
	if src.PriorityPolicy != 0 {
		dst.PriorityPolicy = src.PriorityPolicy
	}
	if src.PinnedTTL != 0 {
		dst.PinnedTTL = src.PinnedTTL
	}
	if src.PriorityGroups != nil {
		dst.PriorityGroups = append([]string(nil), src.PriorityGroups...)
	}
	if src.DeliverSubject != "" {
		dst.DeliverSubject = src.DeliverSubject
	}
	if src.DeliverGroup != "" {
		dst.DeliverGroup = src.DeliverGroup
	}
	if src.FlowControl {
		dst.FlowControl = true
	}
	if src.IdleHeartbeat != 0 {
		dst.IdleHeartbeat = src.IdleHeartbeat
	}
}

func cloneConsumerConfig(src jetstream.ConsumerConfig) jetstream.ConsumerConfig {
	dst := src
	dst.BackOff = append([]time.Duration(nil), src.BackOff...)
	dst.FilterSubjects = append([]string(nil), src.FilterSubjects...)
	dst.PriorityGroups = append([]string(nil), src.PriorityGroups...)
	dst.Metadata = cloneStringMap(src.Metadata)
	if src.OptStartTime != nil {
		t := *src.OptStartTime
		dst.OptStartTime = &t
	}
	if src.PauseUntil != nil {
		t := *src.PauseUntil
		dst.PauseUntil = &t
	}
	return dst
}

func cloneStringMap(src map[string]string) map[string]string {
	if src == nil {
		return nil
	}
	dst := make(map[string]string, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}
