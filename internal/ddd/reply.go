package ddd

import (
	"context"
	"time"

	"github.com/htquangg/microservices-poc/pkg/uid"
)

type (
	Reply interface {
		IDer
		ReplyName() string
		Payload() ReplyPayload
		Metadata() Metadata
		OccurredAt() time.Time
	}

	ReplyHandler[T Reply] interface {
		HandleReply(ctx context.Context, reply T) error
	}

	ReplyHandlerFunc[T Reply] func(ctx context.Context, reply T) error

	ReplyPayload interface{}

	ReplyOption interface {
		configureReply(*reply)
	}

	reply struct {
		Entity
		payload    interface{}
		metadata   Metadata
		occurredAt time.Time
	}
)

var _ Reply = (*reply)(nil)

func NewReply(name string, payload ReplyPayload, options ...ReplyOption) Reply {
	return newReply(name, payload, options...)
}

func newReply(name string, payload ReplyPayload, options ...ReplyOption) reply {
	rep := reply{
		Entity:     NewEntity(uid.GetManager().ID(), name),
		payload:    payload,
		metadata:   make(Metadata),
		occurredAt: time.Now(),
	}

	for _, option := range options {
		option.configureReply(&rep)
	}

	return rep
}

func (e reply) ReplyName() string {
	return e.EntityName()
}

func (e reply) Payload() ReplyPayload {
	return e.payload
}

func (e reply) Metadata() Metadata {
	return e.metadata
}

func (e reply) OccurredAt() time.Time {
	return e.occurredAt
}

func (f ReplyHandlerFunc[T]) HandleReply(ctx context.Context, reply T) error {
	return f(ctx, reply)
}
