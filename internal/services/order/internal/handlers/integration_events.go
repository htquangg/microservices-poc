package handlers

import (
	"context"

	"github.com/htquangg/di/v2"
	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/services/basket/basketpb"
	"github.com/htquangg/microservices-poc/internal/services/order/constants"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/application/commands"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/database"
)

type integrationEventHandlers[T ddd.Event] struct {
	app *application.Application
}

func NewIntegrationEventHandlers(reg registry.Registry, app *application.Application, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewEventHandler(reg, integrationEventHandlers[ddd.Event]{
		app: app,
	}, mws...)
}

func RegisterIntegrationEventHandlers(container di.Container, db database.DB) error {
	rawMsgHandler := am.MessageHandlerFunc(func(ctx context.Context, msg am.IncomingMessage) error {
		ctx = container.Scoped(ctx)
		return db.WithTx(ctx, func(ctx context.Context) error {
			return di.Get(ctx, constants.IntegrationEventHandlersKey).(am.MessageHandler).HandleMessage(ctx, msg)
		})
	})

	subsciber := container.Get(constants.MessageSubscriberKey).(am.MessageSubscriber)

	return registerIntegrationEventHandlers(subsciber, rawMsgHandler)
}

func registerIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) (err error) {
	if _, err = subscriber.Subscribe(basketpb.BasketAggregateChannel, handlers, am.MessageFilter{
		basketpb.BasketCheckedOutEvent,
	}); err != nil {
		return err
	}

	return nil
}

func (h integrationEventHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case basketpb.BasketCheckedOutEvent:
		return h.onBasketCheckedOut(ctx, event)
	}

	return nil
}

func (h integrationEventHandlers[T]) onBasketCheckedOut(ctx context.Context, event T) error {
	payload := event.Payload().(*basketpb.BasketCheckedOut)

	items := make([]*domain.Item, 0, len(payload.GetItems()))
	for _, item := range payload.GetItems() {
		items = append(
			items,
			domain.NewItem(
				item.GetProductId(),
				item.GetStoreId(),
				item.GetPrice(),
				int(item.GetQuantity()),
				domain.WithProductName(item.GetProductName()),
				domain.WithStoreName(item.GetStoreName()),
			),
		)
	}

	return h.app.Commands.CreateOrderHandler.Handle(ctx, commands.CreateOrder{
		ID:         payload.GetId(),
		CustomerID: payload.GetCustomerId(),
		PaymentID:  payload.GetPaymentId(),
		Items:      items,
	})
}
