package kafka

import (
	"time"

	"github.com/htquangg/microservices-poc/pkg/logger"

	"github.com/segmentio/kafka-go"
)

type ConsumerConfig struct {
	Brokers        []string
	Log            logger.Logger
	CommitInterval time.Duration
	Concurrency    int
	reader         kafka.ReaderConfig
}

func (cfg *ConsumerConfig) setDefaults() {
	if cfg.Concurrency == 0 {
		cfg.Concurrency = 1
	}

	if cfg.CommitInterval == 0 {
		cfg.CommitInterval = time.Second
		// Kafka-go library default value is 0, we need to also change this.
		cfg.reader.CommitInterval = time.Second
	} else {
		cfg.reader.CommitInterval = cfg.CommitInterval
	}

	cfg.reader.Brokers = cfg.Brokers
	cfg.reader.MinBytes = minBytes
	cfg.reader.MaxBytes = maxBytes
	cfg.reader.QueueCapacity = queueCapacity
	cfg.reader.HeartbeatInterval = heartbeatInterval
	cfg.reader.PartitionWatchInterval = partitionWatchInterval
	cfg.reader.Logger = kafka.LoggerFunc(cfg.Log.Printf)
	cfg.reader.ErrorLogger = kafka.LoggerFunc(cfg.Log.Errorf)
	cfg.reader.MaxAttempts = maxAttempts
	cfg.reader.MaxWait = maxWait
	cfg.reader.Dialer = &kafka.Dialer{
		Timeout:   dialTimeout,
		KeepAlive: dialKeepAlive,
	}
}
