// Package nats_stream_consumer provides a small, application-oriented wrapper
// around JetStream consumers.
package nats_stream_consumer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var (
	// ErrConsumerAlreadyStarted is returned when a wrapper is started more than
	// once. Use a separate NastStreamConsumer for each independent subscription.
	ErrConsumerAlreadyStarted = errors.New("nats consumer already started")

	// ErrHandlerRequired is returned when a nil handler is supplied.
	ErrHandlerRequired = errors.New("nats consumer handler is required")
)

// NastStreamConsumer is an application-facing JetStream consumer.
//
// A wrapper represents one JetStream consumer and one subscription. The
// wrapper creates a pull consumer by default. Set DeliverGroup or
// DeliverSubject in the ConsumerConfig to opt into push mode; DeliverSubject
// is an internal NATS delivery subject and is never used as the application
// subject. The subject passed to Consume is always installed as FilterSubject.
type NastStreamConsumer struct {
	js         jetstream.JetStream
	streamName string
	option     jetstream.ConsumerConfig

	mu           sync.Mutex
	started      bool
	subscription jetstream.ConsumeContext

	errorHandler func(error)
	logger       *slog.Logger
}

// NewNastStreamConsumer creates a consumer with sensible JetStream defaults.
// Multiple config pointers are merged for compatibility with the original
// API; later non-zero fields take precedence.
func NewNastStreamConsumer(js jetstream.JetStream, streamName string, options ...*jetstream.ConsumerConfig) *NastStreamConsumer {
	return &NastStreamConsumer{
		js:         js,
		streamName: streamName,
		option:     parseOptions(options...),
		logger:     slog.Default(),
	}
}

// SetErrorHandler installs a callback for runtime consumption errors such as
// connection loss, missing heartbeats, or a server-side consumer error.
// Setup errors are returned directly by Consume and ConsumeFunc.
//
// SetErrorHandler must be called before starting the consumer.
func (c *NastStreamConsumer) SetErrorHandler(handler func(error)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errorHandler = handler
}

// SetLogger changes the logger used when no error handler is installed. Pass
// nil to disable default logging.
func (c *NastStreamConsumer) SetLogger(logger *slog.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger = logger
}

// Consume is the simple, process-oriented API. It starts the consumer and
// keeps the calling goroutine alive until SIGINT or SIGTERM is received. This
// means a worker can simply call:
//
//	consumer.Consume("payment.created", worker)
//
// Applications that already own a lifecycle context should use
// ConsumeWithContext instead.
func (c *NastStreamConsumer) Consume(topic string, handler Consumer) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if c == nil {
		slog.Default().Error("nats consumer is nil")
		return
	}
	if handler == nil {
		c.reportError(nil, c.logger, ErrHandlerRequired)
		return
	}

	subscription, err := c.start(ctx, topic, func(msg jetstream.Msg) error {
		return handler.Consume(ctx, msg)
	})
	if err != nil {
		c.reportError(nil, c.logger, err)
		return
	}
	defer subscription.Drain()
	<-ctx.Done()
}

// ConsumeWithContext starts the consumer and stops it when ctx is cancelled.
// The returned ConsumeContext must be retained by the caller and stopped or
// drained as part of application shutdown.
func (c *NastStreamConsumer) ConsumeWithContext(ctx context.Context, topic string, handler Consumer) (jetstream.ConsumeContext, error) {
	if handler == nil {
		return nil, ErrHandlerRequired
	}
	if ctx == nil {
		return nil, errors.New("nats consumer context is nil")
	}

	return c.start(ctx, topic, func(msg jetstream.Msg) error {
		return handler.Consume(ctx, msg)
	})
}

// start creates the JetStream consumer synchronously, so configuration and
// subscription errors are returned to the caller instead of being hidden in a
// goroutine. The underlying nats.go ConsumeContext owns the subscription and
// handles reconnects; the watcher below binds it to ctx.
func (c *NastStreamConsumer) start(ctx context.Context, topic string, handler func(jetstream.Msg) error) (jetstream.ConsumeContext, error) {
	if c == nil || c.js == nil {
		return nil, errors.New("nats consumer is nil")
	}
	if strings.TrimSpace(c.streamName) == "" {
		return nil, errors.New("nats consumer stream name is required")
	}
	if strings.TrimSpace(topic) == "" {
		return nil, errors.New("nats consumer subject is required")
	}
	if ctx == nil {
		return nil, errors.New("nats consumer context is nil")
	}
	if handler == nil {
		return nil, ErrHandlerRequired
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	c.mu.Lock()
	if c.started {
		c.mu.Unlock()
		return nil, ErrConsumerAlreadyStarted
	}
	c.started = true
	cfg := c.configForTopic(topic)
	errorHandler := c.errorHandler
	logger := c.logger
	c.mu.Unlock()

	callback := func(msg jetstream.Msg) {
		if err := invokeHandler(handler, msg); err != nil {
			c.reportError(errorHandler, logger, fmt.Errorf("handle message on %q: %w", msg.Subject(), err))
			if cfg.AckPolicy != jetstream.AckNonePolicy {
				if ackErr := msg.Nak(); ackErr != nil {
					c.reportError(errorHandler, logger, fmt.Errorf("nak message on %q: %w", msg.Subject(), ackErr))
				}
			}
			return
		}

		// AckNone consumers do not have an acknowledgement subject. For the
		// explicit and all acknowledgement policies, a successful handler is
		// acknowledged by the wrapper and a handler error is Nak'ed above.
		if cfg.AckPolicy != jetstream.AckNonePolicy {
			if err := msg.Ack(); err != nil {
				c.reportError(errorHandler, logger, fmt.Errorf("ack message on %q: %w", msg.Subject(), err))
			}
		}
	}

	cons, err := c.createConsumer(ctx, cfg)
	if err != nil {
		c.resetStart()
		return nil, fmt.Errorf("create consumer for stream %q subject %q: %w", c.streamName, topic, err)
	}

	consumeContext, err := cons.consume(callback, jetstream.ConsumeErrHandler(func(_ jetstream.ConsumeContext, consumeErr error) {
		if consumeErr != nil {
			c.reportError(errorHandler, logger, fmt.Errorf("consume stream %q subject %q: %w", c.streamName, topic, consumeErr))
		}
	}))
	if err != nil {
		c.resetStart()
		return nil, fmt.Errorf("start consumer for stream %q subject %q: %w", c.streamName, topic, err)
	}

	c.mu.Lock()
	c.subscription = consumeContext
	c.mu.Unlock()

	go func() {
		select {
		case <-ctx.Done():
			// Drain allows already buffered messages to finish. The underlying
			// ConsumeContext makes this operation idempotent.
			consumeContext.Drain()
		case <-consumeContext.Closed():
		}
		c.cleanupEphemeral(cfg, cons, errorHandler, logger)
	}()

	return consumeContext, nil
}

func (c *NastStreamConsumer) cleanupEphemeral(cfg jetstream.ConsumerConfig, cons consumerHandle, errorHandler func(error), logger *slog.Logger) {
	if cfg.Durable != "" {
		return
	}

	name := cons.name()
	if name == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.js.DeleteConsumer(ctx, c.streamName, name); err != nil && !errors.Is(err, jetstream.ErrConsumerNotFound) {
		c.reportError(errorHandler, logger, fmt.Errorf("delete ephemeral consumer %q: %w", name, err))
	}
}

func (c *NastStreamConsumer) createConsumer(ctx context.Context, cfg jetstream.ConsumerConfig) (consumerHandle, error) {
	if isPushConfig(cfg) {
		consumer, err := c.js.CreateOrUpdatePushConsumer(ctx, c.streamName, cfg)
		if err != nil {
			return nil, err
		}
		return pushConsumerHandle{consumer: consumer}, nil
	}

	consumer, err := c.js.CreateOrUpdateConsumer(ctx, c.streamName, cfg)
	if err != nil {
		return nil, err
	}
	return pullConsumerHandle{consumer: consumer}, nil
}

func (c *NastStreamConsumer) resetStart() {
	c.mu.Lock()
	c.started = false
	c.subscription = nil
	c.mu.Unlock()
}

func (c *NastStreamConsumer) reportError(handler func(error), logger *slog.Logger, err error) {
	if err == nil {
		return
	}
	if handler != nil {
		handler(err)
		return
	}
	if logger != nil {
		logger.Error("nats jetstream consumer error", "stream", c.streamName, "error", err)
	}
}

// configForTopic applies the one piece of configuration that belongs to the
// Consume call. FilterSubject is the application subject; DeliverSubject is
// reserved for push transport and is generated when it was not supplied.
func (c *NastStreamConsumer) configForTopic(topic string) jetstream.ConsumerConfig {
	cfg := cloneConsumerConfig(c.option)
	cfg.FilterSubject = topic
	cfg.FilterSubjects = nil
	if cfg.DeliverSubject == topic {
		// DeliverSubject is transport-only. Treating the application subject
		// as the delivery subject creates a push consumer on the wrong subject.
		cfg.DeliverSubject = ""
	}

	if isPushConfig(cfg) {
		if cfg.DeliverSubject == "" {
			cfg.DeliverSubject = deliverySubject(c.streamName, cfg)
		}
	} else {
		// A pull consumer must not carry push-only fields. In particular, do
		// not turn the application subject into DeliverSubject.
		cfg.DeliverSubject = ""
		cfg.DeliverGroup = ""
	}

	return cfg
}

func isPushConfig(cfg jetstream.ConsumerConfig) bool {
	return cfg.DeliverSubject != "" || cfg.DeliverGroup != ""
}

func deliverySubject(stream string, cfg jetstream.ConsumerConfig) string {
	name := cfg.Durable
	if name == "" {
		name = cfg.Name
	}
	if name == "" {
		// Queue subscribers must share a delivery subject even when each
		// underlying consumer is ephemeral; otherwise every subscriber gets a
		// different subject and the queue group cannot load-balance.
		if cfg.DeliverGroup != "" {
			return fmt.Sprintf("_NAST.%s.queue.%s", stream, cfg.DeliverGroup)
		}
		return nats.NewInbox()
	}
	return fmt.Sprintf("_NAST.%s.%s", stream, name)
}

// consumerHandle is the common subset of pull and push JetStream consumers.
type consumerHandle interface {
	consume(jetstream.MessageHandler, jetstream.ConsumeErrHandler) (jetstream.ConsumeContext, error)
	name() string
}

// pullConsumerHandle adapts a pull consumer to consumerHandle.
type pullConsumerHandle struct{ consumer jetstream.Consumer }

func (c pullConsumerHandle) consume(handler jetstream.MessageHandler, errHandler jetstream.ConsumeErrHandler) (jetstream.ConsumeContext, error) {
	return c.consumer.Consume(handler, errHandler)
}

func (c pullConsumerHandle) name() string {
	if info := c.consumer.CachedInfo(); info != nil {
		return info.Name
	}
	return ""
}

// pushConsumerHandle adapts a push consumer to consumerHandle. Push and pull
// consume options are intentionally kept separate by nats.go, so the adapter
// only forwards the handler and relies on the consumer-level settings.
type pushConsumerHandle struct{ consumer jetstream.PushConsumer }

func (c pushConsumerHandle) consume(handler jetstream.MessageHandler, errHandler jetstream.ConsumeErrHandler) (jetstream.ConsumeContext, error) {
	return c.consumer.Consume(handler, errHandler)
}

func (c pushConsumerHandle) name() string {
	if info := c.consumer.CachedInfo(); info != nil {
		return info.Name
	}
	return ""
}

func invokeHandler(handler func(jetstream.Msg) error, msg jetstream.Msg) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("handler panic: %v", recovered)
		}
	}()
	return handler(msg)
}
