package kafka

import (
	"context"

	"github.com/htquangg/microservices-poc/pkg/logger"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	brokers []string
	log     logger.Logger
	w       *kafka.Writer
}

func NewProducer(brokers []string, log logger.Logger) *Producer {
	return &Producer{
		brokers: brokers,
		log:     log,
		w:       NewWriter(brokers, log),
	}
}

func (p *Producer) Writer() *kafka.Writer {
	return p.w
}

func (p *Producer) Publish(ctx context.Context, msgs ...kafka.Message) error {
	return p.w.WriteMessages(ctx, msgs...)
}

func (p *Producer) Close() error {
	return p.w.Close()
}
