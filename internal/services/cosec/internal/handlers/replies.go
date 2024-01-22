package handlers

import (
	"context"

	"github.com/htquangg/di/v2"
	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/sec"
	"github.com/htquangg/microservices-poc/internal/services/cosec/constants"
	"github.com/htquangg/microservices-poc/internal/services/cosec/internal/saga"
	"github.com/htquangg/microservices-poc/internal/services/cosec/models"
	"github.com/htquangg/microservices-poc/pkg/database"
)

func NewReplyHandlers(
	reg registry.Registry,
	orchestrator sec.Orchestrator[*models.CreateOrderData],
	mws ...am.MessageHandlerMiddleware,
) am.MessageHandler {
	return am.NewReplyHandler(reg, orchestrator, mws...)
}

func RegisterReplyHandlers(container di.Container, db database.DB) error {
	rawMsgHandler := am.MessageHandlerFunc(func(ctx context.Context, msg am.IncomingMessage) error {
		ctx = container.Scoped(ctx)
		return db.WithTx(ctx, func(ctx context.Context) error {
			return di.Get(ctx, constants.ReplyHandlersKey).(am.MessageHandler).HandleMessage(ctx, msg)
		})
	})

	subsciber := container.Get(constants.MessageSubscriberKey).(am.MessageSubscriber)

	return registerReplyHandlers(subsciber, rawMsgHandler)
}

func registerReplyHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) (err error) {
	if _, err = subscriber.Subscribe(saga.CreateOrderReplyChannel, handlers); err != nil {
		return err
	}

	return nil
}
