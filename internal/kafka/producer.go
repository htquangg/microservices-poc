package kafka

import (
	"context"
	"time"

	"github.com/htquangg/microservices-poc/internal/am"
	proto_msq "github.com/htquangg/microservices-poc/internal/kafka/proto"
	"github.com/htquangg/microservices-poc/pkg/logger"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Producer struct {
	log logger.Logger
	w   *kafka.Writer
}

func NewProducer(brokers []string, log logger.Logger) *Producer {
	return &Producer{
		log: log,
		w:   NewWriter(brokers, log),
	}
}

func (p *Producer) Publish(ctx context.Context, topicName string, rawMsgs ...am.Message) error {
	logFields := logger.Fields{
		"topic": topicName,
	}

	msgs := make([]kafka.Message, 0, len(rawMsgs))

	for _, rawMsg := range rawMsgs {
		logFields["message_id"] = rawMsg.ID()
		p.log.Debugw("sending message to queue", logFields)

		metadata, err := structpb.NewStruct(rawMsg.Metadata())
		if err != nil {
			return err
		}

		data, err := proto.Marshal(&proto_msq.Message{
			Id:       rawMsg.ID(),
			Name:     rawMsg.Subject(),
			Data:     rawMsg.Data(),
			Metadata: metadata,
			SentAt:   timestamppb.New(rawMsg.SentAt()),
		})
		if err != nil {
			return err
		}

		msgs = append(msgs,
			kafka.Message{
				Topic: rawMsg.Subject(),
				Value: data,
				Time:  time.Now(),
			},
		)
	}

	err := p.publish(ctx, msgs...)
	if err != nil {
		return err
	}

	p.log.Debug("messages sent to queue")

	return nil
}

func (p *Producer) Close() error {
	return p.w.Close()
}

func (p *Producer) publish(ctx context.Context, msgs ...kafka.Message) error {
	return p.w.WriteMessages(ctx, msgs...)
}
