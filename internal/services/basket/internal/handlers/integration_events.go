package handlers

import (
	"context"

	"github.com/htquangg/di/v2"
	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/services/basket/constants"
	"github.com/htquangg/microservices-poc/internal/services/basket/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/database"

	"github.com/htquangg/microservices-poc/internal/services/store/storepb"
)

type intergrationHandlers[T ddd.Event] struct {
	storeRepo   domain.StoreRepository
	productRepo domain.ProductRepository
}

var _ ddd.EventHandler[ddd.Event] = (*intergrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(
	reg registry.Registry,
	storeRepo domain.StoreRepository,
	productRepo domain.ProductRepository,
	mws ...am.MessageHandlerMiddleware,
) am.MessageHandler {
	return am.NewEventHandler(reg, intergrationHandlers[ddd.Event]{
		storeRepo:   storeRepo,
		productRepo: productRepo,
	}, mws...)
}

func RegisterIntergrationEventHandlers(container di.Container, db database.DB) error {
	rawMsgHandler := am.MessageHandlerFunc(func(ctx context.Context, msg am.IncomingMessage) error {
		ctx = container.Scoped(ctx)
		return db.WithTx(ctx, func(ctx context.Context) error {
			return di.Get(ctx, constants.IntegrationEventHandlersKey).(am.MessageHandler).HandleMessage(ctx, msg)
		})
	})

	subsciber := container.Get(constants.MessageSubscriberKey).(am.MessageSubscriber)

	return registerIntergrationEventHandlers(subsciber, rawMsgHandler)
}

func registerIntergrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(storepb.StoreAggregateChannel, handlers, am.MessageFilter{
		storepb.StoreCreatedEvent,
		storepb.StoreRebrandedEvent,
	})
	return err
}

func (h intergrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case storepb.StoreCreatedEvent:
		return h.onStoreCreated(ctx, event)
	case storepb.StoreRebrandedEvent:
		return h.onStoreRebranded(ctx, event)
	}

	return nil
}

func (h intergrationHandlers[T]) onStoreCreated(ctx context.Context, event T) error {
	payload := (event.Payload()).(*storepb.StoreCreated)
	return h.storeRepo.Add(ctx, &domain.Store{
		ID:   payload.Id,
		Name: payload.Name,
	})
}

func (h intergrationHandlers[T]) onStoreRebranded(ctx context.Context, event T) error {
	payload := (event.Payload()).(*storepb.StoreRebranded)
	return h.storeRepo.Rebrand(ctx, payload.Id, payload.Name)
}
