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

type (
	EventMessage interface {
		MessageBase
		ddd.Event
	}

	IncommingEventMessage interface {
		IncomingMessageBase
		ddd.Event
	}

	EventPublisher interface {
		Publish(ctx context.Context, topicName string, event ddd.Event) error
	}

	eventPublisher struct {
		reg       registry.Registry
		publisher MessagePublisher
	}

	eventMessage struct {
		id         string
		name       string
		payload    ddd.EventPayload
		occurredAt time.Time
		msg        IncomingMessageBase
	}
)

var (
	_ EventMessage   = (*eventMessage)(nil)
	_ EventPublisher = (*eventPublisher)(nil)
)

func NewEventPublisher(
	reg registry.Registry,
	msgPublisher MessagePublisher,
	mws ...MessagePublisherMiddleware,
) EventPublisher {
	return &eventPublisher{
		reg:       reg,
		publisher: MessagePublisherWithMiddleware(msgPublisher, mws...),
	}
}

func (p *eventPublisher) Publish(ctx context.Context, topicName string, event ddd.Event) error {
	payload, err := p.reg.Serialize(event.EventName(), event.Payload())
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&proto_am.EventMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(event.OccurredAt()),
	})
	if err != nil {
		return err
	}

	return p.publisher.Publish(ctx, topicName, message{
		id:       event.ID(),
		name:     event.EventName(),
		subject:  topicName,
		data:     data,
		metadata: event.Metadata(),
		sentAt:   time.Now(),
	})
}

func (m eventMessage) ID() string {
	return m.id
}

func (m eventMessage) Subject() string {
	return m.msg.Subject()
}

func (m eventMessage) EventName() string {
	return m.name
}

func (m eventMessage) MessageName() string {
	return m.msg.MessageName()
}

func (m eventMessage) Payload() ddd.EventPayload {
	return m.payload
}

func (m eventMessage) Metadata() ddd.Metadata {
	return m.msg.Metadata()
}

func (m eventMessage) OccurredAt() time.Time {
	return m.occurredAt
}

func (m eventMessage) SentAt() time.Time {
	return m.msg.SentAt()
}

func (m eventMessage) ReceivedAt() time.Time {
	return m.msg.ReceivedAt()
}

func (e eventMessage) Ack() error {
	return e.msg.Ack()
}

func (e eventMessage) NAck() error {
	return e.msg.NAck()
}

func (e eventMessage) Extend() error {
	return e.msg.Extend()
}

func (e eventMessage) Kill() error {
	return e.msg.Kill()
}

type eventMessageHandler struct {
	reg     registry.Registry
	handler ddd.EventHandler[ddd.Event]
}

func NewEventHandler(
	reg registry.Registry,
	handler ddd.EventHandler[ddd.Event],
	mws ...MessageHandlerMiddleware,
) MessageHandler {
	return MessageHandlerWithMiddleware(&eventMessageHandler{
		reg:     reg,
		handler: handler,
	}, mws...)
}

func (h *eventMessageHandler) HandleMessage(ctx context.Context, msg IncomingMessage) error {
	var eventData proto_am.EventMessageData

	err := proto.Unmarshal(msg.Data(), &eventData)
	if err != nil {
		return err
	}

	eventName := msg.MessageName()

	payload, err := h.reg.Deserialize(eventName, eventData.GetPayload())
	if err != nil {
		return err
	}

	eventMsg := eventMessage{
		id:         msg.ID(),
		name:       eventName,
		payload:    payload,
		occurredAt: eventData.GetOccurredAt().AsTime(),
		msg:        msg,
	}

	return h.handler.HandleEvent(ctx, eventMsg)
}
