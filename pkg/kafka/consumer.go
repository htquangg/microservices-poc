package kafka

import (
	"context"

	"github.com/htquangg/microservices-poc/pkg/logger"
)

type Consumer struct {
	brokers []string
	log     logger.Logger
	closing chan struct{}
}

func NewConsumer(brokers []string, log logger.Logger) *Consumer {
	return &Consumer{
		brokers: brokers,
		log:     log,
		closing: make(chan struct{}),
	}
}

func (c *Consumer) Subscribe(ctx context.Context, topic string) error {
	panic("unimplemented")
}
