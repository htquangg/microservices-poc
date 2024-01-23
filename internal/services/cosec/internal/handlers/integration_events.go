package handlers

import (
	"context"

	"github.com/htquangg/di/v2"
	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/sec"
	"github.com/htquangg/microservices-poc/internal/services/cosec/constants"
	"github.com/htquangg/microservices-poc/internal/services/cosec/models"
	"github.com/htquangg/microservices-poc/internal/services/order/orderpb"
	"github.com/htquangg/microservices-poc/pkg/database"
)

type integrationHandler[T ddd.Event] struct {
	orchestrator sec.Orchestrator[*models.CreateOrderData]
}

func NewIntegrationEventHandlers(
	reg registry.Registry,
	orchestrator sec.Orchestrator[*models.CreateOrderData],
	mws ...am.MessageHandlerMiddleware,
) am.MessageHandler {
	return am.NewEventHandler(
		reg,
		integrationHandler[ddd.Event]{
			orchestrator: orchestrator,
		}, mws...)
}

func RegisterIntegrationEventHandlers(ctn di.Container, db database.DB) error {
	rawMsgHandler := am.MessageHandlerFunc(func(ctx context.Context, msg am.IncomingMessage) error {
		return db.WithTx(ctx, func(ctx context.Context) error {
			return ctn.Get(constants.IntegrationEventHandlersKey).(am.MessageHandler).HandleMessage(ctx, msg)
		})
	})

	subsciber := ctn.Get(constants.MessageSubscriberKey).(am.MessageSubscriber)

	return registerIntegrationEventHandlers(subsciber, rawMsgHandler)
}

func registerIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) (err error) {
	if _, err = subscriber.Subscribe(orderpb.OrderAggregateChannel, handlers, am.MessageFilter{
		orderpb.OrderCreatedEvent,
	}); err != nil {
		return err
	}

	return nil
}

func (h integrationHandler[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case orderpb.OrderCreatedEvent:
		return h.onOrderCreated(ctx, event)
	}

	return nil
}

func (h integrationHandler[T]) onOrderCreated(ctx context.Context, event T) error {
	payload := event.Payload().(*orderpb.OrderCreated)

	var total float64
	items := make([]*models.Item, 0, len(payload.GetItems()))
	for _, item := range payload.GetItems() {
		items = append(items, &models.Item{
			ProductID: item.GetProductId(),
			StoreID:   item.GetStoreId(),
			Price:     item.GetPrice(),
			Quantity:  int(item.GetQuantity()),
		})
		total += float64(item.GetQuantity() * int32(item.GetPrice()))
	}

	data := &models.CreateOrderData{
		OrderID:    payload.GetId(),
		CustomerID: payload.GetCustomerId(),
		PaymentID:  payload.GetPaymentId(),
		ShoppingID: payload.GetShoppingId(),
		Items:      items,
		Total:      total,
	}

	return h.orchestrator.Start(ctx, event.ID(), data)
}
