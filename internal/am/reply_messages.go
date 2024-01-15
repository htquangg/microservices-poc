package am

import (
	"context"
	"time"

	proto_am "github.com/htquangg/microservices-poc/internal/am/proto"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/registry"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	SuccessReply = "am.Success"
	FailureReply = "am.Failure"

	OutcomeSuccess = "SUCCESS"
	OutcomeFailure = "FAILURE"

	ReplyHandlerPrefix  = "REPLY_"
	ReplyNameHandler    = ReplyHandlerPrefix + "NAME"
	ReplyOutcomeHandler = ReplyHandlerPrefix + "OUTCOME"
)

type (
	ReplyMessage interface {
		MessageBase
		ddd.Reply
	}

	IncomingReplyMessage interface {
		IncomingMessageBase
		ddd.Reply
	}

	ReplyPublisher interface {
		Publish(ctx context.Context, topicName string, reply ddd.Reply) error
	}

	replyPublisher struct {
		reg       registry.Registry
		publisher MessagePublisher
	}

	replyMessage struct {
		id         string
		name       string
		payload    ddd.ReplyPayload
		occurredAt time.Time
		msg        IncomingMessageBase
	}

	replyMessageHandler struct {
		reg     registry.Registry
		handler ddd.ReplyHandler[ddd.Reply]
	}
)

var (
	_ ReplyPublisher = (*replyPublisher)(nil)
	_ ReplyMessage   = (*replyMessage)(nil)
)

func NewReplyPublisher(
	reg registry.Registry,
	publisher MessagePublisher,
	mws ...MessagePublisherMiddleware,
) ReplyPublisher {
	return &replyPublisher{
		reg:       reg,
		publisher: publisher,
	}
}

func (r *replyPublisher) Publish(ctx context.Context, topicName string, reply ddd.Reply) (err error) {
	var payload []byte

	if reply.ReplyName() != SuccessReply && reply.ReplyName() != FailureReply {
		payload, err = r.reg.Serialize(reply.ReplyName(), reply.Payload())
		if err != nil {
			return err
		}
	}

	data, err := proto.Marshal(&proto_am.ReplyMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(reply.OccurredAt()),
	})
	if err != nil {
		return err
	}

	return r.publisher.Publish(ctx, topicName, message{
		id:       reply.ID(),
		subject:  topicName,
		name:     reply.ReplyName(),
		data:     data,
		metadata: reply.Metadata(),
		sentAt:   time.Now(),
	})
}

func (r replyMessage) ID() string {
	return r.id
}

func (r replyMessage) ReplyName() string {
	return r.name
}

func (r replyMessage) Payload() ddd.ReplyPayload {
	return r.payload
}

func (r replyMessage) Metadata() ddd.Metadata {
	return r.msg.Metadata()
}

func (r replyMessage) OccurredAt() time.Time {
	return r.occurredAt
}

func (r replyMessage) Subject() string {
	return r.msg.Subject()
}

func (r replyMessage) MessageName() string {
	return r.msg.MessageName()
}

func (r replyMessage) SentAt() time.Time {
	return r.msg.SentAt()
}

func (r replyMessage) ReceivedAt() time.Time {
	return r.msg.ReceivedAt()
}

func (r replyMessage) Ack() error {
	return r.msg.Ack()
}

func (r replyMessage) NAck() error {
	return r.msg.NAck()
}

func (r replyMessage) Extend() error {
	return r.msg.Extend()
}

func (r replyMessage) Kill() error {
	return r.msg.Kill()
}

func NewReplyHandler(
	reg registry.Registry,
	handler ddd.ReplyHandler[ddd.Reply],
	mws ...MessageHandlerMiddleware,
) MessageHandler {
	return MessageHandlerWithMiddleware(replyMessageHandler{
		reg:     reg,
		handler: handler,
	}, mws...)
}

func (h replyMessageHandler) HandleMessage(ctx context.Context, msg IncomingMessage) error {
	var replyData proto_am.ReplyMessageData

	err := proto.Unmarshal(msg.Data(), &replyData)
	if err != nil {
		return err
	}

	replyName := msg.MessageName()

	var payload any

	if replyName != SuccessReply && replyName != FailureReply {
		payload, err = h.reg.Deserialize(replyName, replyData.GetPayload())
		if err != nil {
			return err
		}
	}

	replyMsg := &replyMessage{
		id:         msg.ID(),
		name:       replyName,
		payload:    payload,
		occurredAt: replyData.GetOccurredAt().AsTime(),
		msg:        msg,
	}

	return h.handler.HandleReply(ctx, replyMsg)
}
