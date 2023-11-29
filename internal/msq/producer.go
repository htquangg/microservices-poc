package msq

import (
	"context"
	"time"

	"github.com/htquangg/microservices-poc/internal/am"
	proto_msq "github.com/htquangg/microservices-poc/internal/msq/proto"
	"github.com/htquangg/microservices-poc/pkg/kafka"
	"github.com/htquangg/microservices-poc/pkg/logger"

	seg_kafka "github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type KafkaProducer struct {
	log logger.Logger
	p   *kafka.Producer
}

func NewKafkaProducer(brokers []string, log logger.Logger) *KafkaProducer {
	return &KafkaProducer{
		log: log,
		p:   kafka.NewProducer(brokers, log),
	}
}

func (p *KafkaProducer) Publish(ctx context.Context, topicName string, rawMsgs ...am.Message) error {
	logFields := logger.Fields{
		"topic": topicName,
	}

	msgs := make([]seg_kafka.Message, 0, len(rawMsgs))

	for _, rawMsg := range rawMsgs {
		logFields["message_id"] = rawMsg.ID()
		p.log.Debugw("Sending message to kafka", logFields)

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
			seg_kafka.Message{
				Topic: rawMsg.Subject(),
				Value: data,
				Time:  time.Now(),
			},
		)
	}

	err := p.p.Publish(ctx, msgs...)
	if err != nil {
		return err
	}

	p.log.Debug("Messages sent to kafka")

	return nil
}
