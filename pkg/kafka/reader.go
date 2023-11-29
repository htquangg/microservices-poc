package kafka

import (
	"time"

	"github.com/htquangg/microservices-poc/pkg/logger"

	"github.com/segmentio/kafka-go"
)

func NewReader(brokers []string, log logger.Logger) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:                brokers,
		MinBytes:               minBytes,
		MaxBytes:               maxBytes,
		QueueCapacity:          queueCapacity,
		HeartbeatInterval:      heartbeatInterval,
		CommitInterval:         commitInterval,
		PartitionWatchInterval: partitionWatchInterval,
		Logger:                 kafka.LoggerFunc(log.Printf),
		ErrorLogger:            kafka.LoggerFunc(log.Errorf),
		MaxAttempts:            maxAttempts,
		MaxWait:                time.Second,
		Dialer: &kafka.Dialer{
			Timeout: dialTimeout,
		},
	})
}
