package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/htquangg/microservices-poc/internal/am"
	proto_kafka "github.com/htquangg/microservices-poc/internal/kafka/proto"
	"github.com/htquangg/microservices-poc/pkg/logger"

	"github.com/segmentio/kafka-go"

	"google.golang.org/protobuf/proto"
)

type Consumer struct {
	cfg     *ConsumerConfig
	mu      sync.Mutex
	log     logger.Logger
	closing chan struct{}
	subs    []*subscription
}

func NewConsumer(cfg *ConsumerConfig) *Consumer {
	cfg.setDefaults()
	return &Consumer{
		cfg:     cfg,
		log:     cfg.Log,
		closing: make(chan struct{}),
	}
}

func (c *Consumer) Subscribe(
	topic string,
	handler am.MessageHandler,
	options ...am.SubscriberOption,
) (am.Subscription, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	subCfg := am.NewSubscriberConfig(options)

	cfgReader := c.cfg.reader
	cfgReader.Topic = topic

	s := &subscription{
		r:           kafka.NewReader(cfgReader),
		log:         c.cfg.Log,
		messageCh:   make(chan *kafka.Message, c.cfg.Concurrency),
		quit:        make(chan struct{}),
		concurrency: c.cfg.Concurrency,
		consumeFn:   c.handleMsg(subCfg, handler),
	}

	s.ctx, s.cancelFn = context.WithCancel(context.Background())
	s.consume()

	c.subs = append(c.subs, s)

	return s, nil
}

func (c *Consumer) handleMsg(subCfg am.SubscriberConfig, handler am.MessageHandler) func(*kafka.Message) error {
	var filters map[string]struct{}
	if len(subCfg.MessageFilters()) > 0 {
		filters = make(map[string]struct{}, len(subCfg.MessageFilters()))
		for _, key := range subCfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	return func(kafkaMsg *kafka.Message) error {
		m := &proto_kafka.Message{}
		if err := proto.Unmarshal(kafkaMsg.Value, m); err != nil {
			c.log.Errorf("unmarshal the *kafka.Message %s", kafkaMsg)
		}

		if filters != nil {
			if _, exists := filters[m.GetName()]; !exists {
				c.log.Warn("filtered message")
				return nil
			}
		}

		msg := &rawMessage{
			id:         m.GetId(),
			name:       m.GetName(),
			subject:    kafkaMsg.Topic,
			data:       m.GetData(),
			metadata:   m.GetMetadata().AsMap(),
			sentAt:     m.SentAt.AsTime(),
			receivedAt: time.Now(),
			acked:      false,
			// ackFn:      func() error {},
			// nackFn:     func() error {},
			// extendFn:   func() error {},
			// killFn:     func() error {},
		}

		return handler.HandleMessage(context.Background(), msg)
	}
}

func (c *Consumer) Unsubscribe() error {
	for _, sub := range c.subs {
		err := sub.Unsubscribe()
		if err != nil {
			return err
		}
	}
	return nil
}
