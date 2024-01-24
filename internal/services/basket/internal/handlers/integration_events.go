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

type integrationHandlers[T ddd.Event] struct {
	storeRepo   domain.StoreRepository
	productRepo domain.ProductRepository
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(
	reg registry.Registry,
	storeRepo domain.StoreRepository,
	productRepo domain.ProductRepository,
	mws ...am.MessageHandlerMiddleware,
) am.MessageHandler {
	return am.NewEventHandler(reg, integrationHandlers[ddd.Event]{
		storeRepo:   storeRepo,
		productRepo: productRepo,
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
	if _, err = subscriber.Subscribe(storepb.StoreAggregateChannel, handlers, am.MessageFilter{
		storepb.StoreCreatedEvent,
		storepb.StoreRebrandedEvent,
	}); err != nil {
		return err
	}

	if _, err = subscriber.Subscribe(storepb.ProductAggregateChannel, handlers, am.MessageFilter{
		storepb.ProductAddedEvent,
	}); err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case storepb.StoreCreatedEvent:
		return h.onStoreCreated(ctx, event)
	case storepb.StoreRebrandedEvent:
		return h.onStoreRebranded(ctx, event)
	case storepb.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	}

	return nil
}

func (h integrationHandlers[T]) onStoreCreated(ctx context.Context, event T) error {
	payload := (event.Payload()).(*storepb.StoreCreated)
	return h.storeRepo.Add(ctx, &domain.Store{
		ID:   payload.Id,
		Name: payload.Name,
	})
}

func (h integrationHandlers[T]) onStoreRebranded(ctx context.Context, event T) error {
	payload := (event.Payload()).(*storepb.StoreRebranded)
	return h.storeRepo.Rebrand(ctx, payload.Id, payload.Name)
}

func (h integrationHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*storepb.ProductAdded)
	return h.productRepo.Add(
		ctx,
		payload.GetId(),
		payload.GetStoreId(),
		payload.GetName(),
		payload.GetSku(),
		payload.GetPrice(),
	)
}
