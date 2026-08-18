package nats_stream_consumer

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

func TestParseOptionsPreservesJetStreamDeliverySettings(t *testing.T) {
	cfg := parseOptions(&jetstream.ConsumerConfig{
		Durable:        "payments",
		DeliverPolicy:  jetstream.DeliverByStartSequencePolicy,
		OptStartSeq:    42,
		MaxDeliver:     -1,
		MaxAckPending:  -1,
		BackOff:        []time.Duration{time.Second, 2 * time.Second},
		FilterSubjects: []string{"payment.created", "payment.updated"},
		Metadata:       map[string]string{"owner": "wallet"},
	})

	if cfg.Durable != "payments" || cfg.DeliverPolicy != jetstream.DeliverByStartSequencePolicy || cfg.OptStartSeq != 42 {
		t.Fatalf("consumer delivery settings were not preserved: %+v", cfg)
	}
	if cfg.MaxDeliver != -1 || cfg.MaxAckPending != -1 {
		t.Fatalf("negative JetStream limits were not preserved: %+v", cfg)
	}
	if len(cfg.BackOff) != 2 || len(cfg.FilterSubjects) != 2 || cfg.Metadata["owner"] != "wallet" {
		t.Fatalf("extended consumer settings were not preserved: %+v", cfg)
	}
	if cfg.AckPolicy != jetstream.AckExplicitPolicy || cfg.ReplayPolicy != jetstream.ReplayInstantPolicy {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestConfigForTopicUsesFilterSubjectAndNotDeliverSubject(t *testing.T) {
	consumer := NewNastStreamConsumer(nil, "PAYMENTS")

	cfg := consumer.configForTopic("payment.created")
	if cfg.FilterSubject != "payment.created" {
		t.Fatalf("FilterSubject = %q", cfg.FilterSubject)
	}
	if cfg.DeliverSubject != "" || cfg.DeliverGroup != "" {
		t.Fatalf("pull consumer unexpectedly has push settings: %+v", cfg)
	}
	if cfg.DeliverPolicy != jetstream.DeliverAllPolicy {
		t.Fatalf("default DeliverPolicy = %v, want DeliverAllPolicy", cfg.DeliverPolicy)
	}

	push := NewNastStreamConsumer(nil, "PAYMENTS", &jetstream.ConsumerConfig{
		DeliverGroup: "workers",
	})
	pushCfg := push.configForTopic("payment.created")
	if pushCfg.FilterSubject != "payment.created" || pushCfg.DeliverGroup != "workers" {
		t.Fatalf("push configuration was not applied correctly: %+v", pushCfg)
	}
	if pushCfg.DeliverSubject == "" || pushCfg.DeliverSubject == "payment.created" {
		t.Fatalf("invalid internal DeliverSubject: %q", pushCfg.DeliverSubject)
	}
}

func TestConsumeFuncAcknowledgesAndNaks(t *testing.T) {
	js := newFakeJetStream()
	consumer := NewNastStreamConsumer(js, "PAYMENTS")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cc, err := consumer.ConsumeFuncWithContext(ctx, "payment.created", func(msg jetstream.Msg) error {
		return nil
	})
	if err != nil {
		t.Fatalf("ConsumeFuncWithContext: %v", err)
	}

	ackMsg := &fakeMsg{subject: "payment.created"}
	js.pull.deliver(ackMsg)
	if ackMsg.ackCount() != 1 || ackMsg.nakCount() != 0 {
		t.Fatalf("successful handler ack/nak counts = %d/%d, want 1/0", ackMsg.ackCount(), ackMsg.nakCount())
	}

	if _, err := consumer.ConsumeFunc("payment.created", func(jetstream.Msg) error { return nil }); !errors.Is(err, ErrConsumerAlreadyStarted) {
		t.Fatalf("second start error = %v, want ErrConsumerAlreadyStarted", err)
	}

	cancel()
	select {
	case <-cc.Closed():
	case <-time.After(time.Second):
		t.Fatal("consumer did not close after context cancellation")
	}
}

func TestConsumeFuncNaksHandlerErrors(t *testing.T) {
	js := newFakeJetStream()
	consumer := NewNastStreamConsumer(js, "PAYMENTS")
	cc, err := consumer.ConsumeFunc("payment.created", func(jetstream.Msg) error {
		return errors.New("temporary failure")
	})
	if err != nil {
		t.Fatalf("ConsumeFunc: %v", err)
	}
	defer cc.Stop()

	msg := &fakeMsg{subject: "payment.created"}
	js.pull.deliver(msg)
	if msg.ackCount() != 0 || msg.nakCount() != 1 {
		t.Fatalf("failed handler ack/nak counts = %d/%d, want 0/1", msg.ackCount(), msg.nakCount())
	}
}

func TestDurableConsumerKeepsExistingDeliveryPolicy(t *testing.T) {
	js := newFakeJetStream()
	js.existingPull = &fakePullConsumer{
		ctx: newFakeConsumeContext(),
		info: &jetstream.ConsumerInfo{
			Name: "wallet-transactions",
			Config: jetstream.ConsumerConfig{
				DeliverPolicy: jetstream.DeliverNewPolicy,
			},
		},
	}
	consumer := NewNastStreamConsumer(js, "PAYMENTS", &jetstream.ConsumerConfig{
		Durable:       "wallet-transactions",
		DeliverPolicy: jetstream.DeliverAllPolicy,
	})

	cc, err := consumer.ConsumeFunc("payment.created", func(jetstream.Msg) error { return nil })
	if err != nil {
		t.Fatalf("ConsumeFunc: %v", err)
	}
	defer cc.Stop()

	if js.pullCfg.DeliverPolicy != jetstream.DeliverNewPolicy {
		t.Fatalf("DeliverPolicy = %v, want existing %v", js.pullCfg.DeliverPolicy, jetstream.DeliverNewPolicy)
	}
}

func TestConsumeFuncUsesPushConsumerWhenQueueGroupIsConfigured(t *testing.T) {
	js := newFakeJetStream()
	consumer := NewNastStreamConsumer(js, "PAYMENTS", &jetstream.ConsumerConfig{
		DeliverGroup: "workers",
	})

	cc, err := consumer.ConsumeFunc("payment.created", func(jetstream.Msg) error {
		return nil
	})
	if err != nil {
		t.Fatalf("ConsumeFunc: %v", err)
	}
	defer cc.Stop()

	if js.push.handler == nil {
		t.Fatal("push consumer was not started")
	}
	if js.pushCfg.FilterSubject != "payment.created" || js.pushCfg.DeliverGroup != "workers" {
		t.Fatalf("unexpected push config: %+v", js.pushCfg)
	}
	if js.pull.handler != nil {
		t.Fatal("pull consumer was started for a queue group")
	}
}

type fakeJetStream struct {
	jetstream.JetStream
	pull *fakePullConsumer
	push *fakePushConsumer

	existingPull *fakePullConsumer

	mu      sync.Mutex
	pullCfg jetstream.ConsumerConfig
	pushCfg jetstream.ConsumerConfig
	deleted bool
}

func newFakeJetStream() *fakeJetStream {
	return &fakeJetStream{
		pull: &fakePullConsumer{ctx: newFakeConsumeContext()},
		push: &fakePushConsumer{ctx: newFakeConsumeContext()},
	}
}

func (f *fakeJetStream) CreateOrUpdateConsumer(_ context.Context, _ string, cfg jetstream.ConsumerConfig) (jetstream.Consumer, error) {
	f.mu.Lock()
	f.pullCfg = cfg
	f.mu.Unlock()
	return f.pull, nil
}

func (f *fakeJetStream) Consumer(_ context.Context, _ string, _ string) (jetstream.Consumer, error) {
	if f.existingPull == nil {
		return nil, jetstream.ErrConsumerNotFound
	}
	return f.existingPull, nil
}

func (f *fakeJetStream) CreateOrUpdatePushConsumer(_ context.Context, _ string, cfg jetstream.ConsumerConfig) (jetstream.PushConsumer, error) {
	f.mu.Lock()
	f.pushCfg = cfg
	f.mu.Unlock()
	return f.push, nil
}

func (f *fakeJetStream) DeleteConsumer(_ context.Context, _ string, _ string) error {
	f.mu.Lock()
	f.deleted = true
	f.mu.Unlock()
	return nil
}

type fakePullConsumer struct {
	jetstream.Consumer
	ctx     *fakeConsumeContext
	handler jetstream.MessageHandler
	info    *jetstream.ConsumerInfo
}

func (f *fakePullConsumer) Consume(handler jetstream.MessageHandler, _ ...jetstream.PullConsumeOpt) (jetstream.ConsumeContext, error) {
	f.handler = handler
	return f.ctx, nil
}

func (f *fakePullConsumer) CachedInfo() *jetstream.ConsumerInfo {
	if f.info != nil {
		return f.info
	}
	return &jetstream.ConsumerInfo{Name: "fake-pull"}
}

func (f *fakePullConsumer) deliver(msg jetstream.Msg) {
	f.handler(msg)
}

type fakePushConsumer struct {
	jetstream.PushConsumer
	ctx     *fakeConsumeContext
	handler jetstream.MessageHandler
}

func (f *fakePushConsumer) Consume(handler jetstream.MessageHandler, _ ...jetstream.PushConsumeOpt) (jetstream.ConsumeContext, error) {
	f.handler = handler
	return f.ctx, nil
}

func (f *fakePushConsumer) CachedInfo() *jetstream.ConsumerInfo {
	return &jetstream.ConsumerInfo{Name: "fake-push"}
}

type fakeConsumeContext struct {
	closed chan struct{}
	once   sync.Once
}

func newFakeConsumeContext() *fakeConsumeContext {
	return &fakeConsumeContext{closed: make(chan struct{})}
}

func (f *fakeConsumeContext) Stop()                   { f.once.Do(func() { close(f.closed) }) }
func (f *fakeConsumeContext) Drain()                  { f.Stop() }
func (f *fakeConsumeContext) Closed() <-chan struct{} { return f.closed }

type fakeMsg struct {
	jetstream.Msg
	subject string
	mu      sync.Mutex
	acks    int
	naks    int
}

func (m *fakeMsg) Subject() string { return m.subject }

func (m *fakeMsg) Ack() error {
	m.mu.Lock()
	m.acks++
	m.mu.Unlock()
	return nil
}

func (m *fakeMsg) Nak() error {
	m.mu.Lock()
	m.naks++
	m.mu.Unlock()
	return nil
}

func (m *fakeMsg) ackCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.acks
}

func (m *fakeMsg) nakCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.naks
}

func TestDeliverySubjectForNamedPushConsumerIsStable(t *testing.T) {
	first := deliverySubject("PAYMENTS", jetstream.ConsumerConfig{Durable: "worker"})
	second := deliverySubject("PAYMENTS", jetstream.ConsumerConfig{Durable: "worker"})
	if first != second || !strings.HasPrefix(first, "_NAST.PAYMENTS.worker") {
		t.Fatalf("delivery subject = %q, second = %q", first, second)
	}
}

func TestDeliverySubjectForEphemeralQueueConsumerIsShared(t *testing.T) {
	first := deliverySubject("PAYMENTS", jetstream.ConsumerConfig{DeliverGroup: "workers"})
	second := deliverySubject("PAYMENTS", jetstream.ConsumerConfig{DeliverGroup: "workers"})
	if first != second || !strings.HasPrefix(first, "_NAST.PAYMENTS.queue.workers") {
		t.Fatalf("delivery subject = %q, second = %q", first, second)
	}
}
