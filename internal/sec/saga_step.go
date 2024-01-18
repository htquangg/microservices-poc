package sec

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/ddd"
)

var _ SagaStep[any] = (*sagaStep[any])(nil)

type (
	SagaStep[T any] interface {
		Action(fn StepActionFn[T]) SagaStep[T]
		Compensation(fn StepActionFn[T]) SagaStep[T]
		OnActionReply(replyName string, fn StepReplyHandlerFn[T]) SagaStep[T]
		OnCompensationReply(replyName string, fn StepReplyHandlerFn[T]) SagaStep[T]
		isInvocable(compensating bool) bool
		execute(ctx context.Context, sagaCtx *SagaContext[T]) *stepResult[T]
		handle(ctx context.Context, sagaCtx *SagaContext[T], reply ddd.Reply) error
	}

	StepActionFn[T any]       func(ctx context.Context, data T) (string, ddd.Command, error)
	StepReplyHandlerFn[T any] func(ctx context.Context, data T, reply ddd.Reply) error

	sagaStep[T any] struct {
		actions  map[bool]StepActionFn[T]
		handlers map[bool]map[string]StepReplyHandlerFn[T]
	}

	stepResult[T any] struct {
		ctx         *SagaContext[T]
		destination string
		cmd         ddd.Command
		err         error
	}
)

func (s *sagaStep[T]) Action(fn StepActionFn[T]) SagaStep[T] {
	s.actions[notCompensating] = fn
	return s
}

func (s *sagaStep[T]) Compensation(fn StepActionFn[T]) SagaStep[T] {
	s.actions[isCompensating] = fn
	return s
}

func (s *sagaStep[T]) OnActionReply(replyName string, fn StepReplyHandlerFn[T]) SagaStep[T] {
	s.handlers[notCompensating][replyName] = fn
	return s
}

func (s *sagaStep[T]) OnCompensationReply(replyName string, fn StepReplyHandlerFn[T]) SagaStep[T] {
	s.handlers[isCompensating][replyName] = fn
	return s
}

func (s *sagaStep[T]) execute(ctx context.Context, sagaCtx *SagaContext[T]) *stepResult[T] {
	if action := s.actions[sagaCtx.Compensating]; action != nil {
		destination, cmd, err := action(ctx, sagaCtx.Data)
		return &stepResult[T]{
			ctx:         sagaCtx,
			destination: destination,
			cmd:         cmd,
			err:         err,
		}
	}

	return &stepResult[T]{
		ctx: sagaCtx,
	}
}

func (s *sagaStep[T]) handle(ctx context.Context, sagaCtx *SagaContext[T], reply ddd.Reply) error {
	if handle := s.handlers[sagaCtx.Compensating][reply.ReplyName()]; handle != nil {
		return handle(ctx, sagaCtx.Data, reply)
	}

	return nil
}

func (s *sagaStep[T]) isInvocable(compensating bool) bool {
	return s.actions[compensating] != nil
}
