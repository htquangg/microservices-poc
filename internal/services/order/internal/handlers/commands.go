package handlers

import (
	"context"

	"github.com/htquangg/di/v2"
	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/services/cosec/constants"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/application/commands"
	"github.com/htquangg/microservices-poc/internal/services/order/orderpb"
	"github.com/htquangg/microservices-poc/pkg/database"
)

type commandHandlers struct {
	app *application.Application
}

func NewCommandHandlers(
	reg registry.Registry,
	replyPublisher am.ReplyPublisher,
	app *application.Application,
	mws ...am.MessageHandlerMiddleware,
) am.MessageHandler {
	return am.NewCommandHandler(reg, replyPublisher, &commandHandlers{
		app: app,
	}, mws...)
}

func RegisterCommandHandlers(ctn di.Container, db database.DB) error {
	rawMsgHandler := am.MessageHandlerFunc(func(ctx context.Context, msg am.IncomingMessage) error {
		return db.WithTx(ctx, func(ctx context.Context) error {
			return ctn.Get(constants.CommandHandlersKey).(am.MessageHandler).HandleMessage(ctx, msg)
		})
	})

	subsciber := ctn.Get(constants.MessageSubscriberKey).(am.MessageSubscriber)

	return registerCommandHandlers(subsciber, rawMsgHandler)
}

func registerCommandHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) (err error) {
	if _, err = subscriber.Subscribe(orderpb.CommandChannel, handlers, am.MessageFilter{
		orderpb.ApproveOrderCommand,
		orderpb.RejectOrderCommand,
	}); err != nil {
		return err
	}

	return nil
}

func (h *commandHandlers) HandleCommand(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	switch cmd.CommandName() {
	case orderpb.ApproveOrderCommand:
		return h.doApproveOrder(ctx, cmd)
	}

	return nil, nil
}

func (h *commandHandlers) doApproveOrder(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*orderpb.ApproveOrder)

	return nil, h.app.Commands.ApproveOrderHandler.Handle(ctx, commands.ApproveOrder{
		ID:         payload.GetId(),
		ShoppingID: "",
	})
}
